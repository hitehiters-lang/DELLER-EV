package handlers

import (
	"deller/internal/config"
	"deller/internal/service"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

type ContactMail struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Comment     string `json:"comment"`
}

// @Tags Связь
// @Summary Отправить письмо
// @Description Отправляет письмо на почту компании
// @ID contact_us
// @Accept json
// @Produce json
// @Param contacter body ContactMail true "Поля ввода"
// @Success 200 "{"message": "Успешно! Ждите звонка, мы с вами свяжемся."}"
// @Failure 400 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /contact [post]
func Contact(c *gin.Context) {
	var body ContactMail
	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}
	err := ContactToAdmin(body)
	if err != nil {
		service.ErrorResponse(c, 500, "Повторите попытку позже")
		return
	}
	c.JSON(200, map[string]string{"message": "Успешно! Ждите звонка, мы с вами свяжемся."})
}

func ContactToAdmin(b ContactMail) error {

	mail, pass, mailto := config.GetMail()

	m := gomail.NewMessage()
	m.SetHeader("From", mail)
	m.SetHeader("To", mailto)
	m.SetHeader("Subject", "DellerEv - Новое обращение")

	var body strings.Builder
	body.WriteString("Добрый день\n")
	fmt.Fprintf(&body, "Сегодня (%v) в %v был составлено новое обращение на сайте DellerEV.\n\n",
		time.Now().Format("2006-01-02"),
		time.Now().Format("15:04:05"))
	body.WriteString("Детали обращения:\n\n")

	moscowZone := time.FixedZone("MSK", 3*60*60)
	fmt.Fprintf(&body, "Имя: %v\nНомер телефона: %v\nСоздано: %v\nСообщение: \n\t%v\n",
		b.Name, b.PhoneNumber, b.Comment, time.Now().In(moscowZone).Format("2006-01-02 15:04:05"))

	m.SetBody("text/plain", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, mail, pass)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	fmt.Println("[SERVER] - Сообщение админу доставлено")
	return nil
}
