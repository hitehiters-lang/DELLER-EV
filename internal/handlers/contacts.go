package handlers

import (
	"database/sql"
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/repository"
	"deller/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

// @Tags Связь
// @Summary Отправить письмо
// @Description Отправляет письмо на почту компании
// @ID contact_us
// @Accept json
// @Produce json
// @Param contacter body dto.SendContactMail true "Поля ввода"
// @Success 200 "{"message": "Успешно! Ждите звонка, мы с вами свяжемся."}"
// @Failure 400 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /contact [post]
func Contact(c *gin.Context) {
	var body dto.SendContactMail
	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}
	if body.Email == "" || body.PhoneNumber == "" {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}
	ban := service.CheckForBanwords(body)
	if ban {
		repository.AddMail(c, body, true)
		c.JSON(200, map[string]string{"message": "Успешно! Ждите звонка, мы с вами свяжемся."})
		fmt.Println("[SERVER] - contact message autoban " + body.Comment)
		return
	}
	err := service.ContactToAdmin(body)
	if err != nil {
		service.ErrorResponse(c, 500, "Повторите попытку позже")
		fmt.Println(err)
		return
	}
	err = repository.AddMail(c, body, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.JSON(200, map[string]string{"message": "Успешно! Ждите звонка, мы с вами свяжемся."})
}

// @Tags Связь
// @Summary Получить все письма.
// @Description Получить все письма из бд. Доступно только админам.
// @ID get_mails
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param filter query string false "Фильтр забаненных. n - письма прошедшие проверку (default), b - только забаненные письма, a - все"
// @Success 200 {object} dto.ContactsResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/contact [get]
func GetMails(c *gin.Context) {

	r, _ := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	filter := c.Query("filter")
	if filter == "" {
		filter = "n"
	}

	err, count, mails := repository.GetMails(c, filter)
	if err != nil && err != sql.ErrNoRows {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	} else if count == 0 {
		service.ErrorResponse(c, 404, "Ни одного письма не найдено")
		return
	}
	res := service.ConvertMails(mails)

	c.JSON(200, res)
}
