package repository

import (
	"deller/internal/database"
	"deller/internal/dto"
	"deller/internal/schema"
	"time"

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

	count, err := db.NewSelect().
		Model(&products).
		Where("category_id = ?", cat_id).
		Where("deleted_at < ?", time.Now()).
		ScanAndCount(ctx)
	return err, count, products
}

func GetProductDetail(c *gin.Context, prod_id string) (error, int, schema.Products) {
	var product schema.Products
	ctx := c.Request.Context()

	db := database.GetDatabase()
	count, err := db.NewSelect().Model(&product).Relation("Specification").Where("products.id = ?", prod_id).ScanAndCount(ctx)
	return err, count, product
}

func GetProductsByTag(c *gin.Context, prod_id string) (error, int, []schema.Products) {
	var product schema.Products
	var related_products []schema.Products
	ctx := c.Request.Context()

	db := database.GetDatabase()
	count, err := db.NewSelect().Model(&product).Where("products.id = ?", prod_id).ScanAndCount(ctx)
	if err != nil || count == 0 {
		return err, count, related_products
	}

	tags := product.Tags
	count, err = db.NewSelect().Model(&related_products).Where("products.tags && ?", pgdialect.Array(tags)).
		OrderExpr("RANDOM()").Limit(4).Where("products.id != ?", prod_id).ScanAndCount(ctx)
	if count < 4 {
		var added_products []schema.Products
		needed := 4 - count

		excludedIDs := []string{prod_id}
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

func GetProductPictures(c *gin.Context, prod_id string) []schema.ProductImages {
	var pictures []schema.ProductImages
	ctx := c.Request.Context()

	db := database.GetDatabase()
	db.NewSelect().Model(&pictures).Where("product_id = ?", prod_id).Scan(ctx)
	return pictures
}

func CheckCategory(c *gin.Context, cat_id string) error {
	var categories schema.Categories
	ctx := c.Request.Context()

	db := database.GetDatabase()
	err := db.NewSelect().Model(&categories).Where("id = ?", cat_id).Scan(ctx)
	return err
}

func AddCategory(c *gin.Context, body dto.AddCategoryBody) (schema.Categories, error) {
	new_category := schema.Categories{
		Id:        body.Id,
		Name:      body.Name,
		ImagePath: "/images/" + body.ImageName + ".png",
	}
	ctx := c.Request.Context()

	db := database.GetDatabase()
	err := db.NewInsert().Model(&new_category).Returning("*").Scan(ctx)
	return new_category, err
}

func DeleteCategory(c *gin.Context, cat string) error {
	var del_category schema.Categories
	ctx := c.Request.Context()

	db := database.GetDatabase()
	_, err := db.NewDelete().Model(&del_category).Where("id = ?", cat).Exec(ctx)
	return err
}

func CheckProduct(c *gin.Context, prod_id string) error {
	var products schema.Products
	ctx := c.Request.Context()

	db := database.GetDatabase()
	err := db.NewSelect().Model(&products).Where("id = ?", prod_id).Scan(ctx)
	return err
}

func AddProduct(c *gin.Context, body dto.AddProductBody) (schema.Products, error) {
	spec_uuid, _ := uuid.NewV7()
	s := body.Specification
	new_specs := schema.MotoSpecifications{
		Id:                 spec_uuid,
		Gear:               s.Gear,
		PowerRatedMax:      s.PowerRatedMax,
		TorqueMoment:       s.TorqueMoment,
		MaxSpeed:           s.MaxSpeed,
		Battery:            s.Battery,
		PowerReserve:       s.PowerReserve,
		Weight:             s.Weight,
		LoadCap:            s.LoadCap,
		DimensionsLWH:      s.DimensionsLWH,
		SeatHeight:         s.SeatHeight,
		Tires:              s.Tires,
		Break:              s.Break,
		Clearance:          s.Clearance,
		Charging:           s.Charging,
		MoistureProtection: s.MoistureProtection,
	}
	new_product := schema.Products{
		Id:          body.Id,
		CategoryId:  body.CategoryId,
		Spec:        spec_uuid,
		Name:        body.Name,
		Short:       body.Short,
		Description: body.Description,
		Price:       body.Price,
		Discount:    body.Discount,
		Tags:        body.Tags,
		CreatedAt:   time.Now(),
		MainImage:   body.MainImage,
	}
	ctx := c.Request.Context()

	db := database.GetDatabase()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return schema.Products{}, err
	}
	defer tx.Rollback()
	err = tx.NewInsert().Model(&new_specs).Returning("*").Scan(ctx)
	if err != nil {
		return schema.Products{}, err
	}
	err = tx.NewInsert().Model(&new_product).Returning("*").Scan(ctx)
	if err != nil {
		return schema.Products{}, err
	}

	err = tx.NewSelect().Model(&new_product).Relation("Specification").Where("products.id = ?", new_product.Id).Scan(ctx)
	if err != nil {
		return schema.Products{}, err
	}
	err = tx.Commit()
	if err != nil {
		return schema.Products{}, err
	}
	return new_product, nil
}
