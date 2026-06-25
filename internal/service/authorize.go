package service

import (
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/schema"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(h, p string) error {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	if err != nil {
		return err
	}
	return nil
}

func LoginUser(id uuid.UUID) dto.JWTResponse {
	jwt, err := pkgauth.GenerateJWT(id)
	if err != nil {
		return dto.JWTResponse{}
	}
	return dto.JWTResponse{
		Token: "Bearer " + jwt,
	}
}

func ConvertUser(u schema.Users) dto.User {
	user := dto.User{
		Id:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
	}
	return user
}
