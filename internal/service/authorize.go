package service

import (
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/schema"
	"errors"
	"net/mail"

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

func ValidateRegBody(body dto.RegisterBody) error {
	var ret_err string
	var err error

	if len(body.Password) < 6 {
		ret_err = "Пароль должен быть минимум 6 символов"
	} else if len(body.Password) > 200 {
		ret_err = "Пароль слишком длинный"
	}

	_, merr := mail.ParseAddress(body.Email)
	if merr != nil {
		ret_err = "Электронная почта введена неверно"
	}
	if len(body.Email) > 200 {
		ret_err = "Почта слишком длинная"
	}

	if len(body.Name) < 4 {
		ret_err = "Имя пользователя слишком короткое"
	} else if len(body.Name) > 100 {
		ret_err = "Имя пользователя слишком длинное"
	}

	if ret_err != "" {
		err = errors.New(ret_err)
	}

	return err
}
