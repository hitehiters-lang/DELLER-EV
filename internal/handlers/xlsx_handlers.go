package handlers

import (
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	pkgxlsx "deller/internal/pkg/xlsx"
	"deller/internal/repository"
	"deller/internal/schema"
	"deller/internal/service"
	"os"

	"github.com/gin-gonic/gin"
)

// @Tags Админ
// @Summary Импорт цен из XLSX
// @Description Загружает XLSX-файл с ценами, сопоставляет товары по полю xlsx_name в БД и обновляет цены пачкой в одной транзакции. Файл удаляется после обработки.
// @Description
// @Description ### Как работают колонки
// @Description Модуль **автоматически определяет** колонки с именем и ценой в каждой строке:
// @Description - **Колонка имени**: первая непустая ячейка в строке.
// @Description - **Колонка цены**: следующая ячейка сразу за колонкой имени.
// @Description
// @Description ### Нормализация пробелов
// @Description Перед сопоставлением все пробельные символы в названии (и из XLSX, и из БД) схлопываются до одного пробела через `strings.Fields` + `strings.Join`.
// @Description Это значит, что `"Электромотор  Deller   QSD90A-80      1,0 кВт"` и `"Электромотор Deller QSD90A-80 1,0 кВт"` считаются одним и тем же товаром.
// @Description
// @Description ### Парсинг цены
// @Description Из строки цены удаляются суффиксы `р.`, `р`, `руб.`, `руб`, запятые и пробелы, после чего значение преобразуется в `int`.
// @Description Примеры: `"16,590 р."` → `16590`, `"480 р."` → `480`.
// @Description
// @Description ### Сопоставление с БД
// @Description Товары сопоставляются **точным совпадением** (после нормализации пробелов) текста из колонки имени с полем `xlsx_name` в таблице `products`.
// @Description Если поле `xlsx_name` у товара не заполнено — он не участвует в сопоставлении.
// @Description
// @Description ### Фильтрация мусора
// @Description Строки без цены, с непарсящейся ценой, а также полностью пустые строки — пропускаются.
// @Description
// @Description ---
// @Description
// @Description ### [CODEX] Инструкция для фронтендера
// @Description
// @Description **1. xlsx_name в API-ответах.**
// @Description Поле `xlsx_name` теперь возвращается во **всех** роутах, где фигурирует товар:
// @Description - `GET /categories/:id` — в каждом объекте `Product` внутри `ProductsResponse`
// @Description - `GET /products/:id` — в `ProductDetail`
// @Description - `GET /products/:id/relative` — в каждом объекте `Product`
// @Description Используй это поле чтобы показывать xlsx_name в админке (в списке товаров, при редактировании) и чтобы администратор видел, по какому ключу товар сопоставляется с прайс-листом.
// @Description
// @Description **2. Страница импорта XLSX (новая).**
// @Description Добавь в админ-панель страницу «Импорт цен из XLSX» с формой загрузки файла. Рядом с формой размести пояснение для пользователя:
// @Description > «Загрузите XLSX-файл с прайс-листом. В каждой строке первая непустая ячейка считается названием товара, следующая за ней — ценой. Название товара из XLSX сравнивается с полем xlsx_name в базе. Лишние пробелы в названиях игнорируются. Пример корректной строки: колонка A (пусто), колонка B — "Электромотор Deller QSD90A-80 1,0 кВт V1/V2", колонка C — "16,590 р."»
// @Description После загрузки показывай результат в таблице: строка в файле | название | статус (updated / not_found / parse_error / unchanged) | старая цена | новая цена | ошибка.
// @Description
// @Description **3. Форма добавления/редактирования товара (доработка существующей).**
// @Description При вводе названия товара (поле `name`) автоматически заполняй:
// @Description - `slug` — транслитерация названия латиницей, пробелы → дефисы, только `[a-z0-9-]`. Используй готовую функцию с бекенда или сделай на фронте (пример: `"Электромотор Deller QSD90A-80 1,0 кВт"` → `"elektromotor-deller-qsd90a-80-1-0-kvt"`).
// @Description - `xlsx_name` — формируется как **«[Категория в единственном числе] [Название товара]»**. Например: категория «Электромотоциклы» → единственное число «Электромотоцикл», товар «DellerБайк Электро 400» → `xlsx_name` = `"Электромотоцикл DellerБайк Электро 400"`. Категорию приводи к единственному числу через словарь замен (простой mapping на фронте: «Электромоторы»→«Электромотор», «Контроллеры»→«Контроллер», «Батареи»→«Батарея», «Зарядные устройства»→«Зарядное устройство», «Электромотоциклы»→«Электромотоцикл», «Электропитбайки»→«Электропитбайк», «Карты электро»→«Карт электро», «Аксессуары»→«Аксессуар» и т.д.). Если категория не найдена в словаре — просто отрезай последнюю букву (грубая эвристика).
// @Description Поле `xlsx_name` должно быть видимым и редактируемым в форме (на случай если название в прайс-листе отличается от сгенерированного).
// @Description
// @Description **4. Ручное управление xlsx_name.**
// @Description В списке товаров админки добавь колонку `xlsx_name` (или показывай при редактировании), чтобы администратор мог проверить/подправить значение, по которому идёт сопоставление с прайсом.
// @Description
// @Description **Важно:** xlsx_name в БД храни как есть, без нормализации. Нормализация пробелов происходит только на бекенде в момент сопоставления — множественные пробелы схлопываются при сравнении.
// @ID xlsx_import
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "XLSX-файл с ценами"
// @Success 200 {object} dto.XlsxImportResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/xlsx-import [post]
func ImportPricesFromXlsx(c *gin.Context) {
	r, _ := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		service.ErrorResponse(c, 400, "Файл не найден в запросе")
		return
	}

	filename := file.Filename
	if len(filename) < 5 || filename[len(filename)-5:] != ".xlsx" {
		service.ErrorResponse(c, 400, "Файл должен иметь расширение .xlsx")
		return
	}

	tmpPath := "./shared/tmp_" + filename
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		service.ErrorResponse(c, 500, "Не удалось сохранить файл: "+err.Error())
		return
	}
	defer os.Remove(tmpPath)

	processor := pkgxlsx.NewProcessor()
	result, err := processor.ProcessFile(tmpPath)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка обработки файла: "+err.Error())
		return
	}

	details := make([]dto.XlsxRowDetail, len(result.Details))
	for i, d := range result.Details {
		details[i] = dto.XlsxRowDetail{
			Row:      d.Row,
			Name:     d.Name,
			Status:   d.Status,
			OldPrice: d.OldPrice,
			NewPrice: d.NewPrice,
			Error:    d.Error,
		}
	}

	c.JSON(200, dto.XlsxImportResponse{
		Total:    result.Total,
		Matched:  result.Matched,
		Updated:  result.Updated,
		NotFound: result.NotFound,
		Errors:   result.Errors,
		Details:  details,
	})
}

// @Tags Админ
// @Summary Экспорт прайс-листа в XLSX
// @Description Экспортирует товары в XLSX-файл в формате прайс-листа DellerQS. Файл содержит шапку с реквизитами компании и список товаров с ценами.
// @Description
// @Description ### Параметр category
// @Description - `all` — все товары, сгруппированные по категориям. Между группами — пустая строка-разделитель.
// @Description - `id категории` — товары только указанной категории, без разделителей.
// @Description
// @Description ### Название товара в экспорте
// @Description Используется поле `xlsx_name`. Если оно не заполнено — `name`.
// @Description
// @Description ### Формат цены
// @Description Цена пишется как число с пользовательским форматом `#,##0 "р."` — Excel отобразит `16,590 р.`.
// @Description
// @Description ### Для фронтендера
// @Description [CODEX] Добавь на страницу импорта XLSX кнопку «Экспортировать прайс-лист».
// @Description Выпадающий список с категориями + опция «Все категории». При выборе — GET-запрос на этот эндпоинт с параметром `?category=all` или `?category=<id>`.
// @Description Ответ — файл `.xlsx`, сохраняй через `Content-Disposition: attachment; filename=price.xlsx`.
// @ID xlsx_export
// @Security BearerAuth
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param category query string true "id категории или 'all' для всех"
// @Success 200 {file} file
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/xlsx-export [get]
func ExportPricesToXlsx(c *gin.Context) {
	r, _ := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	category := c.Query("category")
	if category == "" {
		service.ErrorResponse(c, 400, "Параметр category обязателен (id категории или 'all')")
		return
	}

	var products []schema.Products
	var catMap map[string]string
	var groupByCategory bool

	if category == "all" {
		err, prods, cats := repository.GetAllProductsForExport(c)
		if err != nil {
			service.ErrorResponse(c, 500, "Ошибка получения товаров: "+err.Error())
			return
		}
		products = prods
		catMap = cats
		groupByCategory = true
	} else {
		err, prods := repository.GetProductsByCategoryForExport(c, category)
		if err != nil {
			service.ErrorResponse(c, 500, "Ошибка получения товаров: "+err.Error())
			return
		}
		products = prods
		catMap = make(map[string]string)
		groupByCategory = false
	}

	filePath, err := pkgxlsx.ExportPriceList(products, catMap, groupByCategory)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка генерации XLSX: "+err.Error())
		return
	}
	defer os.Remove(filePath)

	c.Header("Content-Disposition", "attachment; filename=price.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File(filePath)
}
