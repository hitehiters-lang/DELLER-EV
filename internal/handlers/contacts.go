package handlers

import (
	"deller/internal/dto"
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
// @Param contacter body dto.ContactMail true "Поля ввода"
// @Success 200 "{"message": "Успешно! Ждите звонка, мы с вами свяжемся."}"
// @Failure 400 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /contact [post]
func Contact(c *gin.Context) {
	var body dto.ContactMail
	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}
	ban := service.CheckForBanwords(body)
	if ban {
		c.JSON(200, map[string]string{"message": "Успешно! Ждите звонка, мы с вами свяжемся."})
		fmt.Println("[SERVER] - contact message autoban " + body.Comment)
		return
	}
	err := service.ContactToAdmin(body)
	if err != nil {
		service.ErrorResponse(c, 500, "Повторите попытку позже")
		return
	}
	c.JSON(200, map[string]string{"message": "Успешно! Ждите звонка, мы с вами свяжемся."})
}
