package dto

type ContactMail struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Comment     string `json:"comment"`
}
