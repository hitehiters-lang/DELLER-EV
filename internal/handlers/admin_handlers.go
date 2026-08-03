package handlers

import (
	"database/sql"
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/repository"
	"deller/internal/service"
	"errors"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var caching int32 = 0

// @Tags Админ
// @Summary Сохранить файл
// @Description Сохраняет файл в статику. Должно использоваться совместно с другими ручками, например добавления руководства.
// @ID add_file
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file false "Файл"
// @Param file_name formData string true "Имя файла (без расширения)"
// @Success 201 {object} dto.UploadMessage
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/files [post]
func AddFile(c *gin.Context) {

	r, _ := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	file_name := c.PostForm("file_name")

	file, err := c.FormFile("file")
	if err != nil {
		service.ErrorResponse(c, 400, "Нет файла")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		service.ErrorResponse(c, 400, "Файл должен иметь расширение")
		return
	}

	core_name, err := service.GenerateUniqueFileName(file_name+ext, "./shared/files")
	if err != nil {
		service.ErrorResponse(c, 400, "Плохое название файла: "+err.Error())
		return
	}

	dst := "./shared/files/" + core_name
	if err := c.SaveUploadedFile(file, dst); err != nil {
		service.ErrorResponse(c, 500, "Не удалось сохранить файл: "+err.Error())
		return
	}

	res := dto.UploadMessage{
		Code:    201,
		Message: "Файл успешно сохранён",
		Link:    "/files/" + core_name,
	}
	c.JSON(201, res)
}

// @Tags Админ
// @Summary Сохранить главную картинку
// @Description Сохраняет картинку в статику. Должно использоваться совместно с другими ручками, например добавления товара или категории.
// @ID add_picture
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "PNG изображение"
// @Param image_name formData string true "Имя изображения"
// @Success 201 {object} dto.UploadMessage
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/images [post]
func AddPicture(c *gin.Context) {

	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	image_name := c.PostForm("image_name")

	file, err := c.FormFile("file")
	if err != nil {
		service.ErrorResponse(c, 400, "Нет файла")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" {
		service.ErrorResponse(c, 400, "Неправильное расширение файла")
		return
	}

	core_name, err := service.GenerateUniqueFileName(image_name+".png", "./shared/images")
	if err != nil {
		service.ErrorResponse(c, 400, "Плохое название файла: "+err.Error())
		return
	}

	dst := "./shared/images/" + core_name
	if err := c.SaveUploadedFile(file, dst); err != nil {
		service.ErrorResponse(c, 500, "Не удалось сохранить файл: "+err.Error())
		return
	}

	res := dto.UploadMessage{
		Code:    201,
		Message: "Файл успешно сохранён",
		Link:    "/images/" + core_name,
	}
	c.JSON(201, res)
}

// @Tags Админ
// @Summary Прикрепить картинку
// @Description Сохраняет картинку в статику. Прикрепляет эту картинку к товару.
// @ID bind_picture
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "PNG изображение"
// @Param image_name formData string true "Имя изображения"
// @Param product_id formData string true "id товара"
// @Success 201 {object} dto.ProductImage
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/image-bind [post]
func BindPicture(c *gin.Context) {
	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	image_name := c.PostForm("image_name")
	product_id := c.PostForm("product_id")

	file, err := c.FormFile("file")
	if err != nil {
		service.ErrorResponse(c, 400, "Нет файла")
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" {
		service.ErrorResponse(c, 400, "Неправильное расширение файла")
		return
	}

	core_name, err := service.GenerateUniqueFileName(image_name+".png", "./shared/images")
	if err != nil {
		service.ErrorResponse(c, 400, "Плохое название файла: "+err.Error())
		return
	}

	err = repository.CheckProduct(c, product_id)
	if err != nil {
		service.ErrorResponse(c, 400, "Продукта не существует")
		return
	}
	pic_uuid, err := repository.BindPicture(c, core_name, product_id)
	if err != nil {
		service.ErrorResponse(c, 400, "Ошибка")
		return
	}
	dst := "./shared/images/" + core_name
	if err := c.SaveUploadedFile(file, dst); err != nil {
		service.ErrorResponse(c, 500, "Не удалось сохранить файл: "+err.Error())
		repository.DeletePicture(c, pic_uuid)
		return
	}
	res := dto.ProductImage{
		Id:        pic_uuid,
		ProductId: product_id,
		ImagePath: "/images/" + core_name,
	}
	c.JSON(201, res)
}

// @Tags Админ
// @Summary Открепить картинку
// @Description Удаляет картинку из БД. Не трогает её в статике. Нужно для замены картинки товара и полного управления ими
// @ID unbind_picture
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id query string true "id изображения"
// @Success 201 {object} dto.Message
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/image-unbind [delete]
func UnbindPicture(c *gin.Context) {
	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	pic_id := c.Query("id")
	pic_uuid, err := uuid.Parse(pic_id)
	if err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	err = repository.DeletePicture(c, pic_uuid)
	if err != nil {
		service.ErrorResponse(c, 404, "Картинки нет в БД")
		return
	}

	res := dto.Message{
		Code:    200,
		Message: "Картинка удалена из БД",
	}
	c.JSON(201, res)
}

// @Tags Админ
// @Summary Создать новую категорию
// @ID add_category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param login body dto.AddCategoryBody true "Данные категории"
// @Success 201 {object} dto.Category
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 409 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/categories [post]
func AddCategory(c *gin.Context) {

	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	var body dto.AddCategoryBody
	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	err = repository.CheckCategory(c, body.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			new_category, err := repository.AddCategory(c, body)
			if err != nil {
				service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
				return
			}
			res := service.ConvertCategory(new_category)
			c.JSON(201, res)
			go service.CacheCategory(new_category)
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}
	service.ErrorResponse(c, 409, "Категория уже существует")
}

// @Tags Админ
// @Summary Удалить категорию
// @ID del_category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "id категории"
// @Success 200 {object} dto.Message
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/categories/{id} [delete]
func DelCategory(c *gin.Context) {

	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	category := c.Param("id")
	if category == "" {
		service.ErrorResponse(c, 400, "Категория не указана")
		return
	}

	err = repository.CheckCategory(c, category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 404, "Категория не существует")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}

	err = repository.DeleteCategory(c, category)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}
	c.JSON(200, dto.Message{
		Code:    200,
		Message: "Категория " + category + " успешно удалена",
	})
}

// @Tags Админ
// @Summary Создать новый товар
// @ID add_product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param login body dto.AddProductBody true "Данные товара"
// @Success 201 {object} dto.ProductDetail
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 409 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/products [post]
func AddProduct(c *gin.Context) {

	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	var body dto.AddProductBody
	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	err = repository.CheckProduct(c, body.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			new_product, err := repository.AddProduct(c, body)
			if err != nil {
				service.ErrorResponse(c, 400, "Ошибка сервера: "+err.Error())
				return
			}
			_, _, prod := repository.GetProductDetail(c, new_product.Id)
			imgs := repository.GetProductPictures(c, new_product.Id)
			img_arr := service.GetImages(imgs)

			res := service.ConvertProductDetail(prod, img_arr, "ru")
			c.JSON(201, res)
			go service.CacheProduct(prod)
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}
	service.ErrorResponse(c, 409, "Товар уже существует")
}

// @Tags Админ
// @Summary Удалить товар
// @ID del_products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "id товара"
// @Success 200 {object} dto.Message
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/products/{id} [delete]
func DelProduct(c *gin.Context) {

	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	product := c.Param("id")
	if product == "" {
		service.ErrorResponse(c, 400, "Товар не указан")
		return
	}

	err = repository.CheckProduct(c, product)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 404, "Товара не существует")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}

	err = repository.DeleteProduct(c, product)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}
	c.JSON(200, dto.Message{
		Code:    200,
		Message: "Товар " + product + " успешно удалён",
	})
}

// @Tags Админ
// @Summary Изменить товар
// @ID put_product
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "id товара"
// @Param login body dto.AddProductBody true "Данные товара"
// @Success 201 {object} dto.ProductDetail
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/products/{id} [put]
func PutProduct(c *gin.Context) {

	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	var body dto.PutProductBody
	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	product := c.Param("id")
	if product == "" {
		service.ErrorResponse(c, 400, "Товар не указан")
		return
	}

	err = repository.CheckProduct(c, product)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 404, "Товара не существует")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}

	new_product, err := repository.PutProduct(c, body, product)
	if err != nil {
		service.ErrorResponse(c, 400, "Ошибка: "+err.Error())
		return
	}
	go service.CacheProduct(new_product)
	c.JSON(200, dto.Message{
		Code:    200,
		Message: "Товар " + product + " успешно изменён",
	})
}

// @Tags Админ
// @Summary Изменить категорию
// @ID put_category
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "id категории"
// @Param login body dto.PutCategoryBody true "Данные товара"
// @Success 201 {object} dto.ProductDetail
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/categories/{id} [put]
func PutCategory(c *gin.Context) {

	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	var body dto.PutCategoryBody
	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	category := c.Param("id")
	if category == "" {
		service.ErrorResponse(c, 400, "Категория не указана")
		return
	}

	err = repository.CheckCategory(c, category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 404, "Категории не существует")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}

	cats, err := repository.PutCategory(c, body, category)
	if err != nil {
		service.ErrorResponse(c, 400, "Ошибка: "+err.Error())
		return
	}
	go service.CacheCategory(cats)
	c.JSON(200, dto.Message{
		Code:    200,
		Message: "Категория " + category + " успешно изменена",
	})
}

// @Tags Админ
// @Summary Кешировать переводы на всём сайте
// @Description Переводит все категории, посты, товары и их характеристики на языки сайта..
// @ID cache_all
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 "{"message":"Кеширование начато"}""
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/cache [post]
func CacheAll(c *gin.Context) {
	if atomic.CompareAndSwapInt32(&caching, 0, 1) {
		c.JSON(200, map[string]string{"message": "Кеширование начато"})

		go func() {
			defer endCaching()
			service.CacheAllContent()
		}()
		return
	}
	c.JSON(200, map[string]string{"message": "Кеширование в процессе"})
}

func endCaching() {
	defer atomic.StoreInt32(&caching, 0)
}
