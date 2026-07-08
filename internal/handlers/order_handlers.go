package handlers

import (
	"database/sql"
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/repository"
	"deller/internal/service"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Tags Заказы
// @Summary Создать заказ
// @Description Добавляет заказ пользователя в бд и высылает ему письмо на почту. Входить необязательно. price необязателен и в бд не записывается. quantity это количество и оно обязательно
// @ID add_order
// @Accept json
// @Produce json
// @Param order body dto.CreateOrder true "Данные заказа"
// @Success 201 "{"message": "Заказ принят, ждите звонка."}"
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 409 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /orders [post]
func AddOrder(c *gin.Context) {

	var body dto.CreateOrder
	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}
	if len(body.Products) > 15 {
		service.ErrorResponse(c, 400, "Превышен лимит максимального количества товаров")
		return
	}
	err := repository.AddOrder(c, body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 400, "Плохой запрос")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}

	c.JSON(201, map[string]string{"message": "Заказ принят, ждите звонка или письмо на электронную почту."})
}

// @Tags Заказы
// @Summary Удалить заказ
// @Description Удаляет заказ, подтверждая что он был отменён или доставлен. Доступно только админам.
// @ID del_order
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "uuid заказа"
// @Success 200 {object} dto.Message
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/orders/{id} [delete]
func DelOrder(c *gin.Context) {

	r, _ := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	order := c.Param("id")
	if order == "" {
		service.ErrorResponse(c, 400, "Заказ не указан")
		return
	}
	order_uuid, err := uuid.Parse(order)
	if err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	err = repository.CheckOrder(c, order_uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 404, "Заказа не существует")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}

	err = repository.DeleteOrder(c, order_uuid)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}
	c.JSON(200, dto.Message{
		Code:    200,
		Message: "Заказ успешно удалён",
	})
}

// @Tags Заказы
// @Summary Получить все заказы
// @Description Получить все заказы из бд. price - цена за 1 товар. Доступно только админам.
// @ID get_orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.OrdersResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/orders [get]
func GetOrders(c *gin.Context) {

	r, _ := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	err, count, orders, order_prods := repository.GetOrders(c)
	if err != nil && err != sql.ErrNoRows {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	} else if count == 0 {
		service.ErrorResponse(c, 404, "Ни одного заказа не найдено")
		return
	}
	res := service.ConvertOrders(orders, order_prods, count)

	c.JSON(200, res)
}
