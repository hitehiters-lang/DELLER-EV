package schema

import (
	"time"

	"github.com/google/uuid"
)

type Posts struct {
	Slug       string    `bun:"slug,pk"`
	Title      string    `bun:"title"`
	Content    string    `bun:"content"`
	AuthorID   uuid.UUID `bun:"author_id"`
	IsApproved bool      `bun:"is_approved"`
	CreatedAt  time.Time `bun:"created_at"`
	UpdatedAt  time.Time `bun:"updated_at,nullzero"`
	DeletedAt  time.Time `bun:"deleted_at,nullzero"`

	Author *Users `bun:"rel:belongs-to,join:author_id=id"`
}
