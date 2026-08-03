package service

import (
	"deller/internal/config"
	"deller/internal/dto"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"gopkg.in/gomail.v2"
)

func CheckForBanwords(b dto.ContactMail) bool {
	data, err := os.ReadFile("shared/banwords_list.txt")
	if err != nil {
		fmt.Println("[SERVER] - can't read banwords_list.txt:", err)
		return true
	}

	commentWords := make(map[string]struct{})
	for _, token := range strings.FieldsFunc(strings.ToLower(b.Comment), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) {
		commentWords[token] = struct{}{}
	}

	for _, line := range strings.Split(string(data), "\n") {
		word := strings.ToLower(strings.TrimSpace(line))
		if word == "" {
			continue
		}
		if _, ok := commentWords[word]; ok {
			return true
		}
	}

	return false
}

func ContactToAdmin(b dto.ContactMail) error {

	mail, pass, mailto := config.GetMail()

	m := gomail.NewMessage()
	m.SetHeader("From", mail)
	m.SetHeader("To", mailto)
	m.SetHeader("Subject", "DellerEv - Новое обращение")

	moscowZone := time.FixedZone("MSK", 3*60*60)

	var body strings.Builder
	body.WriteString("Добрый день\n")
	fmt.Fprintf(&body, "Сегодня (%v) в %v было составлено новое обращение на сайте DellerEV.\n\n",
		time.Now().In(moscowZone).Format("2006-01-02"),
		time.Now().In(moscowZone).Format("15:04:05"))
	body.WriteString("Детали обращения:\n\n")

	fmt.Fprintf(&body, "Имя: %v\nНомер телефона: %v\nСоздано: %v\nСообщение: \n\t%v\n",
		b.Name, b.PhoneNumber, time.Now().In(moscowZone).Format("2006-01-02 15:04:05"), b.Comment)

	m.SetBody("text/plain", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, mail, pass)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	fmt.Println("[SERVER] - Сообщение админу доставлено")
	return nil
}
