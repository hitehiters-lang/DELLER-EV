package blog_gen

import (
	"context"
	"deller/internal/database"
	"deller/internal/schema"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func findOrCreateUser(name string) (uuid.UUID, error) {
	ctx := context.Background()
	db := database.GetDatabase()

	var existing schema.Users
	err := db.NewSelect().Model(&existing).Where("name = ?", name).Scan(ctx)
	if err == nil {
		return existing.Id, nil
	}

	var count int
	count, err = db.NewSelect().Model((*schema.Users)(nil)).Where("name = ?", name).Count(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("check user: %w", err)
	}
	if count > 0 {
		var user schema.Users
		err = db.NewSelect().Model(&user).Where("name = ?", name).Scan(ctx)
		if err != nil {
			return uuid.Nil, fmt.Errorf("get existing user: %w", err)
		}
		return user.Id, nil
	}

	newID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("generate uuid: %w", err)
	}

	emailName := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	email := emailName + "@deller.pro"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("hash password: %w", err)
	}

	newUser := schema.Users{
		Id:        newID,
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		IsAdmin:   false,
		CreatedAt: time.Now(),
	}

	_, err = db.NewInsert().Model(&newUser).Exec(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert user: %w", err)
	}

	return newUser.Id, nil
}
