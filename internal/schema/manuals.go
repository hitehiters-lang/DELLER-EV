package schema

import "time"

type Manuals struct {
	Id        int       `bun:"id,pk,autoincrement"`
	Slug      string    `bun:"slug"`
	Name      string    `bun:"name"`
	ImagePath string    `bun:"image_path"`
	FilePath  string    `bun:"file_path"`
	CreatedAt time.Time `bun:"created_at"`
}
