package dto

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id                uuid.UUID `json:"id"`
	CategoryId        uuid.UUID `json:"category_id"`
	Stack             uuid.UUID `json:"stack"`
	Name              string    `json:"name"`
	Short             string    `json:"short"`
	Price             int       `json:"price"`
	Discount          int       `json:"discount"`
	PriceWithDiscount int       `json:"price_with_discount"`
}

type ProductMore struct {
	Id                uuid.UUID `json:"id"`
	CategoryId        uuid.UUID `json:"category_id"`
	Stack             uuid.UUID `json:"stack"`
	Name              string    `json:"name"`
	Short             string    `json:"short"`
	Description       string    `json:"description"`
	Price             int       `json:"price"`
	Discount          int       `json:"discount"`
	Tags              []string  `json:"tags"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	DeletedAt         time.Time `json:"deleted_at"`
	PriceWithDiscount int       `json:"price_with_discount"`
}

type ProductResponse struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
}
