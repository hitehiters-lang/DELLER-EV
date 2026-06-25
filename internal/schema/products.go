package schema

import (
	"time"

	"github.com/google/uuid"
)

type Categories struct {
	Id        string `bun:"id"`
	Name      string `bun:"name"`
	ImagePath string `bun:"image_path"`

	ProductsCount int `bun:",scanonly"`
}

type Products struct {
	Id          string    `bun:"id,pk,autoincrement"`
	CategoryId  string    `bun:"category_id"`
	Spec        int       `bun:"spec"`
	Name        string    `bun:"name"`
	Short       string    `bun:"short"`
	Description string    `bun:"description"`
	Price       int       `bun:"price"`
	Discount    int       `bun:"discount"`
	Tags        []string  `bun:"tags,array"`
	CreatedAt   time.Time `bun:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at"`
	DeletedAt   time.Time `bun:"deleted_at"`

	Specification *MotoSpecifications `bun:"rel:belongs-to,join:spec=id"`
}

type ProductImages struct {
	Id        uuid.UUID `bun:"id"`
	ProductId string    `bun:"product_id"`
	IsMain    bool      `bun:"is_main"`
	ImagePath string    `bun:"image_path"`
}
