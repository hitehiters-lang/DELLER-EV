package repository

import (
	"deller/internal/database"
	"deller/internal/dto"
	"deller/internal/schema"
	"time"

	"github.com/gin-gonic/gin"
)

func GetManuals(c *gin.Context) (error, int, []schema.Manuals) {
	var manuals []schema.Manuals
	ctx := c.Request.Context()

	db := database.GetDatabase()
	count, err := db.NewSelect().Model(&manuals).ScanAndCount(ctx)
	return err, count, manuals
}

func CheckManual(c *gin.Context, slug string) error {
	var manuals schema.Manuals
	ctx := c.Request.Context()

	db := database.GetDatabase()
	err := db.NewSelect().Model(&manuals).Where("slug = ?", slug).Scan(ctx)
	return err
}

func AddManual(c *gin.Context, body dto.AddManualBody) (schema.Manuals, error) {
	new_manual := schema.Manuals{
		Slug:      body.Slug,
		Name:      body.Name,
		ImagePath: "/images/" + body.ImageName + ".png",
		FilePath:  "/files/" + body.FileName + ".pdf",
		CreatedAt: time.Now(),
	}
	ctx := c.Request.Context()

	db := database.GetDatabase()
	err := db.NewInsert().Model(&new_manual).Returning("*").Scan(ctx)
	return new_manual, err
}

func DeleteManual(c *gin.Context, slug string) error {
	var del_manual schema.Manuals
	ctx := c.Request.Context()

	db := database.GetDatabase()
	_, err := db.NewDelete().Model(&del_manual).Where("slug = ?", slug).Exec(ctx)
	return err
}
