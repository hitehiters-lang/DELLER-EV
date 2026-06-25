package pkgauth

import (
	"deller/internal/database"
	"deller/internal/schema"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCurrentUser(c *gin.Context) (schema.Users, error) {
	var db = database.GetDatabase()
	var user schema.Users
	var ctx = c.Request.Context()

	id, exists := c.Get(UserIDKey)
	if !exists {
		return schema.Users{}, errors.New("user id not found in context")
	}

	userID, ok := id.(uuid.UUID)
	if !ok {
		return schema.Users{}, errors.New("invalid user id type in context")
	}

	err := db.NewSelect().
		Model(&user).
		Where("id = ?", userID).
		Scan(ctx)
	if err != nil {
		return schema.Users{}, err
	}

	return user, nil
}

// r - root /
// a - admin /
// u - user
func GetCurrentRights(c *gin.Context) (string, error) {
	user, err := GetCurrentUser(c)
	if err != nil {
		return "", err
	}

	if user.Name == "root" && user.Id == uuid.MustParse("123e4567-e89b-12d3-a456-426614174111") {
		return "r", nil
	}
	if user.IsAdmin == true {
		return "a", nil
	}
	return "u", nil
}
