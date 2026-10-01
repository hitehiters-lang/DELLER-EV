package service

import (
	"deller/internal/config"
	"deller/internal/dto"
	"deller/internal/schema"
	"fmt"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

func ConvertOrders(ord []schema.Orders, prods []schema.OrderProducts, total int) dto.OrdersResponse {
	var orders []dto.Order
	for _, e := range ord {
		var order_prods []dto.OrderProduct
		for _, p := range prods {
			if p.OrderId == e.Id {
				order_prod := dto.OrderProduct{
					ProductSlug: p.ProductId,
					Quantity:    p.Quantity,
					Price:       p.Price,
				}
				order_prods = append(order_prods, order_prod)
			}
		}
		order := dto.Order{
			Id:          e.Id,
			Name:        e.Name,
			PhoneNumber: e.PhoneNumber,
			Email:       e.Email,
			Address:     e.Address,
			Comment:     e.Comment,
			Products:    order_prods,
			CreatedAt:   e.CreatedAt,
		}
		orders = append(orders, order)
	}
	res := dto.OrdersResponse{
		Orders: orders,
		Total:  total,
	}
	return res
}

func MailToOrderer(ord schema.Orders, prods []schema.OrderProducts) {

	mail, pass, _ := config.GetMail()

	m := gomail.NewMessage()
	m.SetHeader("From", mail)
	m.SetHeader("To", ord.Email)
	m.SetHeader("Subject", "DellerEv - Заказ получен")

	var body strings.Builder
	body.WriteString("Здравствуйте!\n")
	body.WriteString("Администрация DellerEv получила ваш заказ и он уже оформляется!\n")
	body.WriteString("Как только мы сможем, то свяжемся с вами для уточнения деталей или отмены.\n\n")
	body.WriteString("Детали заказа:\n\n")

	var full_price int
	for i, e := range prods {
		fmt.Fprintf(&body, "%d. %v, %dшт. - %dр.\n", i+1, e.Product.Name, e.Quantity, e.TotalPrice)
		full_price += e.TotalPrice
	}
	fmt.Fprintf(&body, "\nИтоговая цена заказа - %dр.\n\n", full_price)

	body.WriteString("(Это письмо является автоматическим подтверждением, что ваш заказ попал в нашу базу данных и обрабатывается)")
	m.SetBody("text/plain", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, mail, pass)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("[SERVER] - Ошибка доставки сообщения: " + err.Error())
		return
	}
	fmt.Println("[SERVER] - Сообщение заказчику доставлено")
}

func MailToAdmin(ord schema.Orders, prods []schema.OrderProducts) {

	mail, pass, mailto := config.GetMail()

	m := gomail.NewMessage()
	m.SetHeader("From", mail)
	m.SetHeader("To", mailto)
	m.SetHeader("Reply-To", ord.Email)
	m.SetHeader("Subject", "DellerEv - Новый заказ")

	moscowZone := time.FixedZone("MSK", 3*60*60)

	var body strings.Builder
	body.WriteString("Добрый день\n")
	fmt.Fprintf(&body, "Сегодня (%v) в %v был составлен новый заказ на сайте DellerEV.\n\n",
		ord.CreatedAt.In(moscowZone).Format("2006-01-02"),
		ord.CreatedAt.In(moscowZone).Format("15:04:05"))
	body.WriteString("Детали заказчика:\n\n")

	fmt.Fprintf(&body, "id: %v\nИмя: %v\nНомер телефона: %v\nЭлектронная почта: %v\nАдрес: %v\nКомментарий: %v\nСоздан: %v\n\n",
		ord.Id, ord.Name, ord.PhoneNumber, ord.Email, ord.Address, ord.Comment, ord.CreatedAt.In(moscowZone).Format("2006-01-02 15:04:05"))

	body.WriteString("---\n\nДетали заказа:\n\n")

	var full_price int
	for i, e := range prods {
		fmt.Fprintf(&body, "%d. %v, %dшт. - %dр.\n", i+1, e.Product.Name, e.Quantity, e.TotalPrice)
		full_price += e.TotalPrice
	}
	fmt.Fprintf(&body, "\nИтоговая цена заказа - %dр.\n\n", full_price)

	body.WriteString("Полный обзор на заказ можно будет сделать в панели администратора на сайте.\n\n")
	m.SetBody("text/plain", body.String())

	d := gomail.NewDialer("smtp.gmail.com", 587, mail, pass)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("[SERVER] - Ошибка доставки сообщения: " + err.Error())
		return
	}
	fmt.Println("[SERVER] - Сообщение админу доставлено")
}
