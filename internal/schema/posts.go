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
	ImagePath  string    `bun:"image_path"`

	Author *Users `bun:"rel:belongs-to,join:author_id=id"`
}

type Comments struct {
	Id        uuid.UUID `bun:"id,pk"`
	AuthorId  uuid.UUID `bun:"author_id"`
	PostSlug  string    `bun:"post_slug"`
	Content   string    `bun:"content"`
	ParentId  uuid.UUID `bun:"parent_id,nullzero"`
	CreatedAt time.Time `bun:"created_at"`
	DeletedAt time.Time `bun:"deleted_at,nullzero"`

	Author *Users    `bun:"rel:belongs-to,join:author_id=id"`
	Post   *Posts    `bun:"rel:belongs-to,join:post_slug=slug"`
	Parent *Comments `bun:"rel:belongs-to,join:parent_id=id"`
}
