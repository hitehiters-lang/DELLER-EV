package handlers

import (
	"database/sql"
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/repository"
	"deller/internal/service"
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Tags Аутентификация
// @Summary Вход
// @Description Вход по логину (email/name) и паролю. Возвращает JWT.
// @ID log_in
// @Accept json
// @Produce json
// @Param login body dto.LoginBody true "Данные для входа"
// @Success 200 {object} dto.JWTResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /auth/login [post]
func LogIn(c *gin.Context) {

	var body dto.LoginBody

	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	user, err := repository.CheckLogin(c, body.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 400, "Неверный логин или пароль")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}

	err = service.CheckPassword(user.Password, body.Password)
	if err != nil {
		service.ErrorResponse(c, 400, "Неверный логин или пароль")
		return
	}

	res := service.LoginUser(user.Id)

	c.JSON(200, res)
}

// @Tags Аутентификация
// @Summary Мой профиль
// @Description Возвращает информацию о своём профиле.
// @ID me
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.User
// @Failure 400 {object} dto.ErrorMessage
// @Failure 401 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /auth/me [get]
func Me(c *gin.Context) {

	self_user, err := pkgauth.GetCurrentUser(c)

	user, err := repository.GetUserByID(c, self_user.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 400, "Пользователя не существует")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}

	res := service.ConvertUser(user)

	c.JSON(200, res)
}

// @Tags Аутентификация
// @Summary Регистрация
// @Description Регистрация пользователя в БД. Недопустимы одинаковые почты или имена. Возвращает JWT. (todo: add email confirmation)
// @ID reg
// @Accept json
// @Produce json
// @Param login body dto.RegisterBody true "Данные для регистрации"
// @Success 200 {object} dto.JWTResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 409 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /auth/reg [post]
func Register(c *gin.Context) {

	var body dto.RegisterBody

	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	err := service.ValidateRegBody(body)
	if err != nil {
		service.ErrorResponse(c, 400, err.Error())
		return
	}

	errf := repository.CheckName(c, body.Name)
	if errf != "" {
		if errf == "user exists" {
			service.ErrorResponse(c, 409, "Пользователь с таким именем уже существует")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка проверки имени: "+errf)
		return
	}

	errf = repository.CheckMail(c, body.Email)
	if errf != "" {
		if errf == "user exists" {
			service.ErrorResponse(c, 409, "Пользователь с такой почтой уже существует")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка проверки почты: "+errf)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка хеширования: "+err.Error())
		return
	}

	id, err := repository.RegisterUser(c, body, hash)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка регистрации: "+err.Error())
		return
	}
	token, err := pkgauth.GenerateJWT(id)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка генерации токена: "+err.Error())
		return
	}

	c.JSON(201, dto.JWTResponse{Token: "Bearer " + token})
}
