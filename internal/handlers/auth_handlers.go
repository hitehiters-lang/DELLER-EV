package handlers

import (
	"database/sql"
	"deller/internal/dto"
	"deller/internal/repository"
	"deller/internal/service"
	"errors"

	"github.com/gin-gonic/gin"
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
