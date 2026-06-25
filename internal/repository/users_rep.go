package repository

import (
	"deller/internal/database"
	"deller/internal/schema"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func GetUserByID(c *gin.Context, id uuid.UUID) (schema.Users, error) {
	var user schema.Users
	ctx := c.Request.Context()
	db := database.GetDatabase()

	err := db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return schema.Users{}, err
	}
	return user, nil
}
