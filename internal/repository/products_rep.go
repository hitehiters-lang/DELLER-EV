package repository

import (
	"database/sql"
	"deller/internal/database"
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/schema"
	"errors"
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
		Join("LEFT JOIN products AS p ON p.category_id = categories.id AND p.deleted_at IS NULL").
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
		Where("deleted_at IS NULL").
		OrderBy("is_trending", bun.OrderDesc).
		ScanAndCount(ctx)
	return err, count, products
}

func GetProductDetail(c *gin.Context, prod_id string) (error, int, schema.Products) {
	var product schema.Products
	ctx := c.Request.Context()

	db := database.GetDatabase()
	count, err := db.NewSelect().Model(&product).Relation("Specification").Where("products.id = ?", prod_id).
		Where("deleted_at IS NULL").ScanAndCount(ctx)
	return err, count, product
}

func GetProductsByTag(c *gin.Context, prod_id string) (error, int, []schema.Products) {
	var product schema.Products
	var related_products []schema.Products
	ctx := c.Request.Context()

	db := database.GetDatabase()
	count, err := db.NewSelect().Model(&product).Where("products.id = ?", prod_id).Where("deleted_at IS NULL").ScanAndCount(ctx)
	if err != nil || count == 0 {
		return err, count, related_products
	}

	tags := product.Tags
	count, err = db.NewSelect().Model(&related_products).Where("products.tags && ?", pgdialect.Array(tags)).
		Limit(4).Where("products.id != ?", prod_id).Where("deleted_at IS NULL").ScanAndCount(ctx)

	if count < 4 {
		excludedIDs := []string{prod_id}
		for _, p := range related_products {
			excludedIDs = append(excludedIDs, p.Id)
		}

		neededFromCategory := 4 - len(related_products)
		if neededFromCategory > 0 && product.CategoryId != "" {
			var category_products []schema.Products
			err = db.NewSelect().
				Model(&category_products).
				Where("products.category_id = ?", product.CategoryId).
				Where("products.id NOT IN ? AND deleted_at IS NULL", bun.Tuple(excludedIDs)).
				OrderExpr("RANDOM()").
				Limit(neededFromCategory).
				Scan(ctx)

			if err == nil && len(category_products) > 0 {
				related_products = append(related_products, category_products...)
				for _, p := range category_products {
					excludedIDs = append(excludedIDs, p.Id)
				}
			}
		}

		neededRandom := 4 - len(related_products)
		if neededRandom > 0 {
			var added_products []schema.Products
			err = db.NewSelect().
				Model(&added_products).
				Where("products.id NOT IN ? AND deleted_at IS NULL", bun.Tuple(excludedIDs)).
				OrderExpr("RANDOM()").
				Limit(neededRandom).
				Scan(ctx)

			if err != nil {
				return err, len(related_products), related_products
			}
			related_products = append(related_products, added_products...)
		}

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
	ct, err := db.NewSelect().Model(&products).Where("id = ?", prod_id).Where("deleted_at IS NULL").ScanAndCount(ctx)
	if ct == 0 {
		err = sql.ErrNoRows
	}
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
		IsTrending:  body.IsTrending,
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

func DeleteProduct(c *gin.Context, prod string) error {
	del_prod := schema.Products{
		DeletedAt: time.Now(),
	}
	ctx := c.Request.Context()

	db := database.GetDatabase()
	_, err := db.NewUpdate().Model(&del_prod).Where("id = ?", prod).Column("deleted_at").Exec(ctx)
	return err
}

func PutProduct(c *gin.Context, body dto.PutProductBody, id string) (schema.Products, error) {
	s := body.Specification
	ctx := c.Request.Context()

	db := database.GetDatabase()
	var prod schema.Products
	err := db.NewSelect().Model(&prod).Relation("Specification").Where("products.id = ?", id).Scan(ctx)
	if err != nil {
		return schema.Products{}, err
	}

	if prod.Spec == uuid.Nil {
		return schema.Products{}, errors.New("product has no specification")
	}

	new_specs := schema.MotoSpecifications{
		Id:                 prod.Spec,
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
		Id:          prod.Id,
		CategoryId:  body.CategoryId,
		Spec:        prod.Spec,
		Name:        body.Name,
		Short:       body.Short,
		Description: body.Description,
		Price:       body.Price,
		Discount:    body.Discount,
		Tags:        body.Tags,
		CreatedAt:   prod.CreatedAt,
		UpdatedAt:   time.Now(),
		MainImage:   body.MainImage,
		IsTrending:  body.IsTrending,
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return schema.Products{}, err
	}
	defer tx.Rollback()
	_, err = tx.NewUpdate().Model(&new_specs).Where("id = ?", prod.Spec).Exec(ctx)
	if err != nil {
		return schema.Products{}, err
	}
	_, err = tx.NewUpdate().Model(&new_product).WherePK().ExcludeColumn("spec").Exec(ctx)
	if err != nil {
		return schema.Products{}, err
	}
	err = tx.NewSelect().Model(&new_product).WherePK().Relation("Specification").Scan(ctx)
	if err != nil {
		return schema.Products{}, err
	}
	err = tx.Commit()
	if err != nil {
		return schema.Products{}, err
	}
	return new_product, nil
}

func PutCategory(c *gin.Context, body dto.PutCategoryBody, id string) (schema.Categories, error) {
	ctx := c.Request.Context()

	db := database.GetDatabase()
	upd_cat := schema.Categories{
		Id:        body.Id,
		Name:      body.Name,
		ImagePath: "/images/" + body.Name + ".png",
	}
	err := db.NewUpdate().Model(&upd_cat).Where("id = ?", id).Returning("*").Scan(ctx)
	if err != nil {
		return schema.Categories{}, err
	}
	return upd_cat, nil
}

func BindPicture(c *gin.Context, image_name, product_id string) (uuid.UUID, error) {
	pic_uuid, _ := uuid.NewV7()
	picture := schema.ProductImages{
		Id:        pic_uuid,
		ProductId: product_id,
		ImagePath: "/images/" + image_name,
	}
	ctx := c.Request.Context()

	db := database.GetDatabase()
	_, err := db.NewInsert().Model(&picture).Exec(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	return pic_uuid, nil
}

func DeletePicture(c *gin.Context, pic_uuid uuid.UUID) error {
	var del_img schema.ProductImages
	ctx := c.Request.Context()

	db := database.GetDatabase()
	res, err := db.NewDelete().Model(&del_img).Where("id = ?", pic_uuid).Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func CheckReview(c *gin.Context, id uuid.UUID) error {
	var reviews []schema.Reviews
	ctx := c.Request.Context()
	db := database.GetDatabase()

	ct, err := db.NewSelect().Model(&reviews).Where("reviews.id = ? AND reviews.deleted_at IS NULL", id).ScanAndCount(ctx)
	if ct == 0 {
		err = sql.ErrNoRows
	}
	return err
}

func AddReview(c *gin.Context, b dto.NewReview, id string, p uuid.UUID) (schema.Reviews, error) {

	author, _ := pkgauth.GetCurrentUser(c)
	rev_uuid, _ := uuid.NewV7()

	add_rev := schema.Reviews{
		Id:        rev_uuid,
		AuthorId:  author.Id,
		ProductId: id,
		Rate:      b.Rate,
		Content:   b.Content,
		ParentId:  p,
		CreatedAt: time.Now(),
	}
	ctx := c.Request.Context()

	db := database.GetDatabase()
	err := db.NewInsert().Model(&add_rev).Scan(ctx)
	if err != nil {
		return schema.Reviews{}, err
	}
	err = db.NewSelect().Model(&add_rev).Where("reviews.id = ?", rev_uuid).Relation("Author").Relation("Parent").Relation("Product").Scan(ctx)
	return add_rev, err
}

func GetReviews(c *gin.Context, id string) (error, int, []schema.Reviews) {
	var reviews []schema.Reviews
	ctx := c.Request.Context()
	db := database.GetDatabase()

	count, err := db.NewSelect().Model((*schema.Reviews)(nil)).Where("reviews.product_id = ?", id).
		Where("reviews.deleted_at IS NULL").Count(ctx)
	if err != nil || count == 0 {
		return err, count, reviews
	}

	err = db.NewSelect().
		Model(&reviews).
		Where("reviews.product_id = ?", id).
		Where("reviews.deleted_at IS NULL").
		Relation("Author").
		Relation("Parent").
		Relation("Product").
		OrderBy("reviews.created_at", bun.OrderAsc).
		Scan(ctx)

	return err, count, reviews
}

func DeleteReview(c *gin.Context, id uuid.UUID) error {
	del_rev := schema.Reviews{
		DeletedAt: time.Now(),
	}
	ctx := c.Request.Context()

	db := database.GetDatabase()
	_, err := db.NewUpdate().Model(&del_rev).Where("id = ?", id).Column("deleted_at").Exec(ctx)
	return err
}

func CheckForRevRepeat(c *gin.Context, id string) (bool, error) {
	var review schema.Reviews
	ctx := c.Request.Context()
	db := database.GetDatabase()
	cur_user, _ := pkgauth.GetCurrentUser(c)

	bl, err := db.NewSelect().
		Model(&review).
		Where("reviews.product_id = ?", id).
		Where("reviews.author_id = ?", cur_user.Id).
		Where("reviews.deleted_at IS NULL").
		Exists(ctx)
	return bl, err
}

func GetAllProductsForExport(c *gin.Context) (error, []schema.Products, map[string]string) {
	var products []schema.Products
	var categories []schema.Categories
	ctx := c.Request.Context()
	db := database.GetDatabase()

	err := db.NewSelect().Model(&products).
		Where("deleted_at IS NULL").
		Order("category_id", "name").
		Scan(ctx)
	if err != nil {
		return err, nil, nil
	}

	err = db.NewSelect().Model(&categories).Scan(ctx)
	if err != nil {
		return err, nil, nil
	}

	catMap := make(map[string]string)
	for _, cat := range categories {
		catMap[cat.Id] = cat.Name
	}

	return nil, products, catMap
}

func GetProductsByCategoryForExport(c *gin.Context, catID string) (error, []schema.Products) {
	var products []schema.Products
	ctx := c.Request.Context()
	db := database.GetDatabase()

	err := db.NewSelect().Model(&products).
		Where("category_id = ?", catID).
		Where("deleted_at IS NULL").
		Order("name").
		Scan(ctx)

	return err, products
}
