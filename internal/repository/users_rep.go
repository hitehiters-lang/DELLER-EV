package repository

import (
	"deller/internal/database"
	"deller/internal/schema"

	"github.com/gin-gonic/gin"
)

func CheckLogin(c *gin.Context, login string) (schema.Users, error) {
	var user schema.Users
	ctx := c.Request.Context()
	db := database.GetDatabase()

	err := db.NewSelect().Model(&user).Where("email = ? OR name = ?", login, login).Scan(ctx)
	if err != nil {
		return schema.Users{}, err
	}
	return user, nil
}
