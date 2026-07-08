package service

import (
	"deller/internal/dto"
	"deller/internal/schema"
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
