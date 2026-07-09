package repository

import (
	"database/sql"
	"deller/internal/database"
	"deller/internal/dto"
	"deller/internal/schema"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddOrder(c *gin.Context, body dto.CreateOrder) (schema.Orders, []schema.OrderProducts, error) {

	order_uuid, _ := uuid.NewV7()
	new_order := schema.Orders{
		Id:          order_uuid,
		Name:        body.Name,
		PhoneNumber: body.PhoneNumber,
		Email:       body.Email,
		Address:     body.Address,
		Comment:     body.Comment,
		CreatedAt:   time.Now(),
	}

	var prods []schema.OrderProducts

	for _, e := range body.Products {

		err, _, p := GetProductDetail(c, e.ProductSlug)
		if err != nil {
			return schema.Orders{}, []schema.OrderProducts{}, err
		}
		if e.Quantity < 1 {
			return schema.Orders{}, []schema.OrderProducts{}, sql.ErrNoRows
		}
		pwd := p.Price - int(float64(p.Price)*(float64(p.Discount)/100))

		prod_uuid, _ := uuid.NewV7()
		prod := schema.OrderProducts{
			Id:         prod_uuid,
			ProductId:  e.ProductSlug,
			OrderId:    order_uuid,
			Quantity:   e.Quantity,
			Price:      pwd,
			TotalPrice: e.Quantity * pwd,
		}
		prods = append(prods, prod)
	}

	ctx := c.Request.Context()
	db := database.GetDatabase()
	tx, err := db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return schema.Orders{}, []schema.OrderProducts{}, err
	}

	_, err = tx.NewInsert().Model(&new_order).Exec(ctx)
	if err != nil {
		return schema.Orders{}, []schema.OrderProducts{}, err
	}
	_, err = tx.NewInsert().Model(&prods).Exec(ctx)
	if err != nil {
		return schema.Orders{}, []schema.OrderProducts{}, err
	}

	var sel_order schema.Orders
	var sel_prods []schema.OrderProducts
	err = tx.NewSelect().Model(&sel_order).Where("id = ?", order_uuid).Scan(ctx)
	if err != nil {
		return schema.Orders{}, []schema.OrderProducts{}, err
	}
	err = tx.NewSelect().Model(&sel_prods).Where("order_id = ?", sel_order.Id).Relation("Product").Scan(ctx)
	if err != nil {
		return schema.Orders{}, []schema.OrderProducts{}, err
	}

	err = tx.Commit()
	return sel_order, sel_prods, err
}

func CheckOrder(c *gin.Context, order_id uuid.UUID) error {
	var order schema.Orders
	ctx := c.Request.Context()

	db := database.GetDatabase()
	ct, err := db.NewSelect().Model(&order).Where("id = ?", order_id).Where("deleted_at IS NULL").ScanAndCount(ctx)
	if ct == 0 {
		err = sql.ErrNoRows
	}
	return err
}

func DeleteOrder(c *gin.Context, order_id uuid.UUID) error {
	del_order := schema.Orders{
		DeletedAt: time.Now(),
	}
	ctx := c.Request.Context()

	db := database.GetDatabase()
	_, err := db.NewUpdate().Model(&del_order).Where("id = ?", order_id).Column("deleted_at").Exec(ctx)
	return err
}

func GetOrders(c *gin.Context) (error, int, []schema.Orders, []schema.OrderProducts) {
	var orders []schema.Orders
	var prods []schema.OrderProducts
	ctx := c.Request.Context()

	db := database.GetDatabase()
	ct, err := db.NewSelect().Model(&orders).Where("deleted_at IS NULL").ScanAndCount(ctx)
	if err != nil {
		return err, 0, nil, nil
	}
	err = db.NewSelect().Model(&prods).Scan(ctx)
	return err, ct, orders, prods
}
