package dto

import "time"

type Manual struct {
	Id        int       `json:"id"`
	Slug      string    `json:"slug"`
	Name      string    `json:"name"`
	ImagePath string    `json:"image_path"`
	FilePath  string    `json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
}

type ManualResponse struct {
	Manuals []Manual `json:"manuals"`
	Total   int      `json:"total"`
}

type AddManualBody struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	ImageName string `json:"image_name"`
	FileName  string `json:"file_name"`
}
