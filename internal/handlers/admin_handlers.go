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

			res := service.ConvertProductDetail(prod, img_arr)
			c.JSON(201, res)
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

	err = repository.PutProduct(c, body, product)
	if err != nil {
		service.ErrorResponse(c, 400, "Ошибка: "+err.Error())
		return
	}
	c.JSON(200, dto.Message{
		Code:    200,
		Message: "Товар " + product + " успешно изменён",
	})
}
