package repository

import (
	"deller/internal/database"
	"deller/internal/schema"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func GetCategories(c *gin.Context) (error, int, []schema.Categories) {
	var categories []schema.Categories
	ctx := c.Request.Context()

	db := database.GetDatabase()
	count, err := db.NewSelect().Model(&categories).
		ColumnExpr("categories.*").
		ColumnExpr("COUNT(p.id) AS products_count").
		Join("LEFT JOIN products AS p ON p.category_id = categories.id").
		Group("categories.id").
		ScanAndCount(ctx)
	return err, count, categories
}

func SelectProductsByCategory(c *gin.Context, cat_id string) (error, int, []schema.Products) {
	var products []schema.Products
	ctx := c.Request.Context()

	db := database.GetDatabase()
	count, err := db.NewSelect().Model(&products).Where("category_id = ?", cat_id).ScanAndCount(ctx)
	return err, count, products
}

func GetProductDetail(c *gin.Context, prod_uuid string) (error, int, schema.Products) {
	var product schema.Products
	ctx := c.Request.Context()

	db := database.GetDatabase()
	count, err := db.NewSelect().Model(&product).Relation("Specification").Where("products.id = ?", prod_uuid).ScanAndCount(ctx)
	return err, count, product
}

func GetProductsByTag(c *gin.Context, prod_uuid string) (error, int, []schema.Products) {
	var product schema.Products
	var related_products []schema.Products
	ctx := c.Request.Context()

	db := database.GetDatabase()
	count, err := db.NewSelect().Model(&product).Where("products.id = ?", prod_uuid).ScanAndCount(ctx)
	if err != nil || count == 0 {
		return err, count, related_products
	}

	tags := product.Tags
	count, err = db.NewSelect().Model(&related_products).Where("products.tags && ?", pgdialect.Array(tags)).
		OrderExpr("RANDOM()").Limit(4).Where("products.id != ?", prod_uuid).ScanAndCount(ctx)
	if count < 4 {
		var added_products []schema.Products
		needed := 4 - count

		excludedIDs := []uuid.UUID{uuid.MustParse(prod_uuid)}
		for _, p := range related_products {
			excludedIDs = append(excludedIDs, p.Id)
		}

		err = db.NewSelect().
			Model(&added_products).
			Where("products.id NOT IN ?", bun.Tuple(excludedIDs)).
			OrderExpr("RANDOM()").
			Limit(needed).
			Scan(ctx)
		if err != nil {
			return err, count, related_products
		}

		related_products = append(related_products, added_products...)
		count = len(related_products)
	}
	return err, count, related_products
}
