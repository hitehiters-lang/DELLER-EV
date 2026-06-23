package schema

import (
	"time"

	"github.com/google/uuid"
)

type Products struct {
	Id          uuid.UUID `bun:"id,pk,autoincrement"`
	CategoryId  uuid.UUID `bun:"category_id"`
	Spec        uuid.UUID `bun:"spec"`
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
