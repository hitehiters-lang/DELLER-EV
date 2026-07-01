package repository

import (
	"deller/internal/database"
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/schema"
	"time"

	"github.com/gin-gonic/gin"
)

func AddPost(c *gin.Context, b dto.NewPost, slug string) (schema.Posts, error) {

	author, _ := pkgauth.GetCurrentUser(c)

	add_post := schema.Posts{
		Slug:       slug,
		Title:      b.Title,
		Content:    b.Content,
		AuthorID:   author.Id,
		IsApproved: true,
		CreatedAt:  time.Now(),
	}
	ctx := c.Request.Context()

	db := database.GetDatabase()
	err := db.NewInsert().Model(&add_post).Returning("*").Scan(ctx)
	return add_post, err
}
