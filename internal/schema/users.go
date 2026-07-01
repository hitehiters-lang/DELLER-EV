package schema

import (
	"time"

	"github.com/google/uuid"
)

type Users struct {
	Id        uuid.UUID `bun:"id,pk"`
	Name      string    `bun:"name"`
	Email     string    `bun:"email"`
	Password  string    `bun:"password"`
	IsAdmin   bool      `bun:"is_admin"`
	CreatedAt time.Time `bun:"created_at"`
}
