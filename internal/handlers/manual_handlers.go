package handlers

import (
	"database/sql"
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/repository"
	"deller/internal/schema"
	"deller/internal/service"
	"errors"

	"github.com/gin-gonic/gin"
)

// @Tags Руководства
// @Summary Получить все руководства
// @Description Возвращает все данные о руководствах
// @ID get_manuals
// @Accept json
// @Produce json
// @Param l query string false "Языковой код (ru, en, kk, zh)"
// @Success 200 {object} dto.ManualResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /manuals [get]
func GetManuals(c *gin.Context) {

	l := service.LoadLocalisation(c)
	if l == "" {
		l = "ru"
	}

	err, count, manuals := repository.GetManuals(c)
	if err != nil && err != sql.ErrNoRows {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	} else if count == 0 {
		service.ErrorResponse(c, 404, "Ни одного руководства не найдено")
		return
	}

	res := service.ConvertManuals(manuals, count, l)

	c.JSON(200, res)
}

// @Tags Руководства
// @Summary Создать новое руководство
// @Description Создаёт новое руководство в БД. Возвращает массив только с созданным руководством. Достуно только админу. Картинка прикрепляется через ручку главных картинок, как к товарам. Файл прикрепляется ручкой в админке. Файл пишется здесь с расширением
// @ID add_manual
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param login body dto.AddManualBody true "Данные руководства"
// @Success 201 {object} dto.ManualResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 409 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/manuals [post]
func AddManual(c *gin.Context) {

	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	var body dto.AddManualBody
	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	err = repository.CheckManual(c, body.Slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			new_manual, err := repository.AddManual(c, body)
			if err != nil {
				service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
				return
			}
			manual_arr := []schema.Manuals{new_manual}
			res := service.ConvertManuals(manual_arr, 1, "ru")
			c.JSON(201, res)
			go service.CacheManual(new_manual)
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}
	service.ErrorResponse(c, 409, "Руководство с таким slug уже существует")
}

// @Tags Руководства
// @Summary Удалить руководство
// @Description Удаляет руководство из БД. Достуно только админу.
// @ID del_manual
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param slug path string true "slug руководства"
// @Success 201 {array} dto.Manual
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 409 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/manuals/{slug} [delete]
func DelManual(c *gin.Context) {

	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	manual_slug := c.Param("slug")
	if manual_slug == "" {
		service.ErrorResponse(c, 400, "slug руководства не указан")
		return
	}

	err = repository.CheckManual(c, manual_slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 404, "Руководство не существует")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}

	err = repository.DeleteManual(c, manual_slug)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}
	c.JSON(200, dto.Message{
		Code:    200,
		Message: "Руководство " + manual_slug + " успешно удалено",
	})
}
