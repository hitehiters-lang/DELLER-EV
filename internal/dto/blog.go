package dto

import "time"

type Post struct {
	Slug          string    `json:"slug"`
	Title         string    `json:"title"`
	AuthorName    string    `json:"author_name"`
	AuthorIsAdmin bool      `json:"author_is_admin"`
	IsApproved    bool      `json:"is_approved"`
	CreatedAt     time.Time `json:"created_at"`
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
}

type NewPost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type SlugReturn struct {
	Slug string `json:"slug"`
}

type PostResponse struct {
	Posts []Post `json:"posts"`
	Total int    `json:"total"`
}
