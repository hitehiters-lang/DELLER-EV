package schema

import (
	"time"

	"github.com/google/uuid"
)

type OrderProducts struct {
	Id         uuid.UUID `bun:"id,pk"`
	ProductId  string    `bun:"product_id"`
	OrderId    uuid.UUID `bun:"order_id"`
	Quantity   int       `bun:"quantity"`
	Price      int       `bun:"price"`
	TotalPrice int       `bun:"total_price"`

	Order   *Orders   `bun:"rel:belongs-to,join:order_id=id"`
	Product *Products `bun:"rel:belongs-to,join:product_id=id"`
}

type Orders struct {
	Id          uuid.UUID `bun:"id,pk"`
	Name        string    `bun:"name"`
	PhoneNumber string    `bun:"phone_number"`
	Email       string    `bun:"email"`
	Address     string    `bun:"address"`
	Comment     string    `bun:"comment"`
	CreatedAt   time.Time `bun:"created_at"`
	DeletedAt   time.Time `bun:"deleted_at,nullzero"`
}
