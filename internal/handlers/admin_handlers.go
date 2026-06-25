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

	"github.com/gin-gonic/gin"
)

// @Tags Админ
// @Summary Сохранить картинку
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
	if r != "a" {
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

	dst := "./shared/" + image_name + ".png"
	if err := c.SaveUploadedFile(file, dst); err != nil {
		service.ErrorResponse(c, 500, "Не удалось сохранить файл: "+err.Error())
		return
	}

	res := dto.UploadMessage{
		Code:    201,
		Message: "Файл успешно сохранён",
		Link:    "/images/" + image_name + ".png",
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
	if r != "a" {
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
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}
	service.ErrorResponse(c, 409, "Категория уже существует")
}
