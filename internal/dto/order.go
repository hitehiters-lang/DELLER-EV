package dto

import (
	"time"

	"github.com/google/uuid"
)

type OrderProduct struct {
	ProductSlug string `json:"product_slug"`
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price,omitempty"`
}

type Order struct {
	Id          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	PhoneNumber string         `json:"phone_number"`
	Email       string         `json:"email"`
	Address     string         `json:"address"`
	Comment     string         `json:"comment"`
	Products    []OrderProduct `json:"products"`
	CreatedAt   time.Time      `json:"created_at"`
}

type CreateOrder struct {
	Name        string         `json:"name"`
	PhoneNumber string         `json:"phone_number"`
	Email       string         `json:"email"`
	Address     string         `json:"address"`
	Comment     string         `json:"comment"`
	Products    []OrderProduct `json:"products"`
}

type OrdersResponse struct {
	Orders []Order `json:"orders"`
	Total  int     `json:"total"`
}
