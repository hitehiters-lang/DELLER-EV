package dto

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Slug          string    `json:"slug"`
	Title         string    `json:"title"`
	AuthorName    string    `json:"author_name"`
	AuthorIsAdmin bool      `json:"author_is_admin"`
	IsApproved    bool      `json:"is_approved"`
	CreatedAt     time.Time `json:"created_at"`
	ImagePath     string    `json:"image_path"`
}

type PostDetail struct {
	Slug          string    `json:"slug"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AuthorName    string    `json:"author_name"`
	AuthorIsAdmin bool      `json:"author_is_admin"`
	IsApproved    bool      `json:"is_approved"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ImagePath     string    `json:"image_path"`
}

type SlugReturn struct {
	Slug string `json:"slug"`
}

type PostResponse struct {
	Posts []Post `json:"posts"`
	Total int    `json:"total"`
}

type Comment struct {
	Id                 uuid.UUID `json:"id"`
	Content            string    `json:"content"`
	AuthorName         string    `json:"author_name"`
	AuthorIsAdmin      bool      `json:"author_is_admin"`
	AuthorIsPostAuthor bool      `json:"author_is_post_author"`
	ParentId           uuid.UUID `json:"parent_id"`
	PostSlug           string    `json:"post_slug"`
	CreatedAt          time.Time `json:"created_at"`
	DeletedAt          time.Time `json:"deleted_at"`
}

type CommentResponse struct {
	Comments []Comment `json:"comments"`
	Total    int       `json:"total"`
}
