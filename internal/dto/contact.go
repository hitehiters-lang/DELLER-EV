package dto

import "time"

type SendContactMail struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Comment     string `json:"comment"`
}

type ContactMail struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Comment     string    `json:"comment"`
	IsBanned    bool      `json:"is_banned"`
	CreatedAt   time.Time `json:"created_at"`
}

type ContactsResponse struct {
	Mails []ContactMail `json:"mails"`
	Total int           `json:"total"`
}
