package repository

import (
	"database/sql"
	"deller/internal/database"
	"deller/internal/dto"
	"deller/internal/schema"
	"errors"
	"time"

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

func CheckName(c *gin.Context, name string) string {
	var user schema.Users
	ctx := c.Request.Context()
	db := database.GetDatabase()

	err := db.NewSelect().Model(&user).Where("name = ?", name).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ""
		}
		return err.Error()
	}
	return "user exists"
}

func CheckMail(c *gin.Context, mail string) string {
	var user schema.Users
	ctx := c.Request.Context()
	db := database.GetDatabase()

	err := db.NewSelect().Model(&user).Where("email = ?", mail).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ""
		}
		return err.Error()
	}
	return "user exists"
}

func RegisterUser(c *gin.Context, b dto.RegisterBody, h []byte) (uuid.UUID, error) {
	ctx := c.Request.Context()
	db := database.GetDatabase()

	user_uuid, _ := uuid.NewV7()
	new_user := schema.Users{
		Id:        user_uuid,
		Name:      b.Name,
		Email:     b.Email,
		Password:  string(h),
		IsAdmin:   false,
		CreatedAt: time.Now(),
	}

	_, err := db.NewInsert().Model(&new_user).Exec(ctx)
	return user_uuid, err
}
