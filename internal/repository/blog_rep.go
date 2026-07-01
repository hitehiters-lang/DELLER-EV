package repository

import (
	"deller/internal/database"
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/schema"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
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

func GetPosts(c *gin.Context, lim, off int) (error, int, []schema.Posts) {
	var posts []schema.Posts
	ctx := c.Request.Context()
	db := database.GetDatabase()

	count, err := db.NewSelect().Model((*schema.Posts)(nil)).Count(ctx)
	if err != nil || count == 0 {
		return err, count, posts
	}

	err = db.NewSelect().
		Model(&posts).
		Where("deleted_at IS NULL").
		Relation("Author").
		OrderBy("created_at", bun.OrderAsc).
		Limit(lim).Offset(off).
		Scan(ctx)

	return err, count, posts
}

func GetPostDetail(c *gin.Context, slug string) (schema.Posts, error) {
	var post schema.Posts
	ctx := c.Request.Context()
	db := database.GetDatabase()

	err := db.NewSelect().
		Model(&post).
		Where("slug = ?", slug).
		Relation("Author").
		Scan(ctx)

	return post, err
}
