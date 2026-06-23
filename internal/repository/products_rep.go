package repository

import (
	"deller/internal/database"
	"deller/internal/schema"

	"github.com/gin-gonic/gin"
)

func SelectProductsByCategory(c *gin.Context, cat_id string) (error, int, []schema.Products) {
	var products []schema.Products

	db := database.GetDatabase()
	count, err := db.NewSelect().Model(&products).Where("category_id = ?", cat_id).ScanAndCount(c.Request.Context())
	return err, count, products
}

func GetProductDetail(c *gin.Context, uuid string) (error, int, schema.Products) {
	var product schema.Products

	db := database.GetDatabase()
	count, err := db.NewSelect().Model(&product).Relation("Specification").Where("products.id = ?", uuid).ScanAndCount(c.Request.Context())
	return err, count, product
}
