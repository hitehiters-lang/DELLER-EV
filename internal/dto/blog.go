package dto

import "time"

type Post struct {
	Slug       string    `json:"slug"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	AuthorName string    `json:"author_name"`
	IsApproved bool      `json:"is_approved"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type NewPost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type SlugReturn struct {
	Slug string `json:"slug"`
}
