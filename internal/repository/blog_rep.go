package repository

import (
	"database/sql"
	"deller/internal/database"
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/schema"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		ImagePath:  "/images/" + b.ImageName + ".png",
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

	count, err := db.NewSelect().Model((*schema.Posts)(nil)).
		Where("deleted_at IS NULL").Count(ctx)
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

	ct, err := db.NewSelect().
		Model(&post).
		Where("slug = ?", slug).
		Where("deleted_at IS NULL").
		Relation("Author").
		ScanAndCount(ctx)
	if ct == 0 {
		err = sql.ErrNoRows
	}

	return post, err
}

func GetComments(c *gin.Context, slug string) (error, int, []schema.Comments) {
	var comments []schema.Comments
	ctx := c.Request.Context()
	db := database.GetDatabase()

	count, err := db.NewSelect().Model((*schema.Comments)(nil)).Where("comments.post_slug = ?", slug).
		Where("comments.deleted_at IS NULL").Count(ctx)
	if err != nil || count == 0 {
		return err, count, comments
	}

	err = db.NewSelect().
		Model(&comments).
		Where("comments.post_slug = ?", slug).
		Where("comments.deleted_at IS NULL").
		Relation("Author").
		Relation("Parent").
		Relation("Post").
		OrderBy("comments.created_at", bun.OrderAsc).
		Scan(ctx)

	return err, count, comments
}

func CheckComment(c *gin.Context, id uuid.UUID) error {
	var comments []schema.Comments
	ctx := c.Request.Context()
	db := database.GetDatabase()

	ct, err := db.NewSelect().Model(&comments).Where("comments.id = ? AND comments.deleted_at IS NULL", id).ScanAndCount(ctx)
	if ct == 0 {
		err = sql.ErrNoRows
	}
	return err
}

func AddComment(c *gin.Context, b dto.NewComment, slug string, p uuid.UUID) (schema.Comments, error) {

	author, _ := pkgauth.GetCurrentUser(c)
	com_uuid, _ := uuid.NewV7()

	add_comm := schema.Comments{
		Id:        com_uuid,
		AuthorId:  author.Id,
		PostSlug:  slug,
		Content:   b.Content,
		ParentId:  p,
		CreatedAt: time.Now(),
	}
	ctx := c.Request.Context()

	db := database.GetDatabase()
	err := db.NewInsert().Model(&add_comm).Scan(ctx)
	if err != nil {
		return schema.Comments{}, err
	}
	err = db.NewSelect().Model(&add_comm).Where("comments.id = ?", com_uuid).Relation("Author").Relation("Parent").Relation("Post").Scan(ctx)
	return add_comm, err
}

func DeletePost(c *gin.Context, slug string) error {
	del_post := schema.Posts{
		DeletedAt: time.Now(),
	}
	ctx := c.Request.Context()

	db := database.GetDatabase()
	_, err := db.NewUpdate().Model(&del_post).Where("slug = ?", slug).Column("deleted_at").Exec(ctx)
	return err
}

func DeleteComment(c *gin.Context, id uuid.UUID) error {
	del_comm := schema.Comments{
		DeletedAt: time.Now(),
	}
	ctx := c.Request.Context()

	db := database.GetDatabase()
	_, err := db.NewUpdate().Model(&del_comm).Where("id = ?", id).Column("deleted_at").Exec(ctx)
	return err
}
