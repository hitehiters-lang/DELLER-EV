package dto

import (
	"time"

	"github.com/google/uuid"
)

type JWTResponse struct {
	Token string `json:"token"`
}

type LoginBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterBody struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type User struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	IsAdmin   bool      `json:"is_admin"`
}
