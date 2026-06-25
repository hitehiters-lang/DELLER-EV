package dto

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	ImagePath       string `json:"image_path"`
	TotalInCategory int    `json:"total_in_category"`
}

type Product struct {
	Id                string    `json:"id"`
	CategoryId        string    `json:"category_id"`
	Spec              uuid.UUID `json:"spec"`
	Name              string    `json:"name"`
	Short             string    `json:"short"`
	Price             int       `json:"price"`
	Discount          int       `json:"discount"`
	PriceWithDiscount int       `json:"price_with_discount"`
	MainImage         string    `json:"main_image"`
}

type ProductDetail struct {
	Id                string         `json:"id"`
	CategoryId        string         `json:"category_id"`
	Spec              uuid.UUID      `json:"spec"`
	Name              string         `json:"name"`
	Short             string         `json:"short"`
	Description       string         `json:"description"`
	Price             int            `json:"price"`
	Discount          int            `json:"discount"`
	Tags              []string       `json:"tags"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         time.Time      `json:"deleted_at"`
	PriceWithDiscount int            `json:"price_with_discount"`
	MainImage         string         `json:"main_image"`
	Specification     Specification  `json:"specification"`
	Images            []ProductImage `json:"images"`
}

type Specification struct {
	Id                 uuid.UUID `json:"id"`
	Gear               string    `json:"gear"`
	PowerRatedMax      string    `json:"power_rated_max"`
	TorqueMoment       string    `json:"torque_moment"`
	MaxSpeed           string    `json:"max_speed"`
	Battery            string    `json:"battery"`
	PowerReserve       string    `json:"power_reserve"`
	Weight             string    `json:"weight"`
	LoadCap            string    `json:"load_cap"`
	DimensionsLWH      string    `json:"dimensions_l_w_h"`
	SeatHeight         string    `json:"seat_height"`
	Tires              string    `json:"tires"`
	Break              string    `json:"break"`
	Clearance          string    `json:"clearance"`
	Charging           string    `json:"charging"`
	MoistureProtection string    `json:"moisture_protection"`
}

type CategoriesResponse struct {
	Categories []Category `json:"categories"`
	Total      int        `json:"total"`
}

type ProductsResponse struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
}

type ProductImage struct {
	Id        uuid.UUID `json:"id"`
	ProductId string    `json:"product_id"`
	IsMain    bool      `json:"is_main"`
	ImagePath string    `json:"image_path"`
}

type AddCategoryBody struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	ImageName string `json:"image_name"`
}
