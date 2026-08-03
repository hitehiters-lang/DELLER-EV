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
	Id          string    `bun:"id,pk"`
	CategoryId  string    `bun:"category_id"`
	Spec        uuid.UUID `bun:"spec"`
	Name        string    `bun:"name"`
	Short       string    `bun:"short"`
	Description string    `bun:"description"`
	Price       int       `bun:"price"`
	Discount    int       `bun:"discount"`
	Tags        []string  `bun:"tags,array"`
	CreatedAt   time.Time `bun:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero"`
	DeletedAt   time.Time `bun:"deleted_at,nullzero"`
	MainImage   string    `bun:"main_image"`
	IsTrending  bool      `bun:"is_trending"`
	XlsxName    string    `bun:"xlsx_name"`
	OnRequest   bool      `bun:"on_request"`
	OnStock     bool      `bun:"on_stock"`

	Specification *MotoSpecifications `bun:"rel:belongs-to,join:spec=id"`
}

type ProductImages struct {
	Id        uuid.UUID `bun:"id,pk"`
	ProductId string    `bun:"product_id"`
	ImagePath string    `bun:"image_path"`
}

type Reviews struct {
	Id        uuid.UUID `bun:"id,pk"`
	AuthorId  uuid.UUID `bun:"author_id"`
	ProductId string    `bun:"product_id"`
	Rate      int       `bun:"rate"`
	Content   string    `bun:"content"`
	ParentId  uuid.UUID `bun:"parent_id,nullzero"`
	CreatedAt time.Time `bun:"created_at"`
	DeletedAt time.Time `bun:"deleted_at,nullzero"`

	Author  *Users    `bun:"rel:belongs-to,join:author_id=id"`
	Product *Products `bun:"rel:belongs-to,join:product_id=id"`
	Parent  *Comments `bun:"rel:belongs-to,join:parent_id=id"`
}
