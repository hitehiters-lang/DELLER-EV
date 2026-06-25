package handlers

import (
	"database/sql"
	"deller/internal/repository"
	"deller/internal/service"

	"github.com/gin-gonic/gin"
)

// @Tags Товары
// @Summary Получить все категории
// @Description Возвращает id и имена категорий
// @ID get_categories
// @Accept json
// @Produce json
// @Success 200 {object} dto.CategoriesResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /categories [get]
func GetCategories(c *gin.Context) {

	err, count, cats := repository.GetCategories(c)
	if err != nil && err != sql.ErrNoRows {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	} else if count == 0 {
		service.ErrorResponse(c, 404, "Ни одной категории не найдено")
		return
	}

	res := service.ConvertCategories(cats)

	c.JSON(200, res)
}

// @Tags Товары
// @Summary Получить все товары определённой категории
// @ID get_products
// @Accept json
// @Produce json
// @Param id path string true "id категории"
// @Success 200 {object} dto.ProductsResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /categories/{id} [get]
func GetProducts(c *gin.Context) {

	category := c.Param("id")
	if category == "" {
		service.ErrorResponse(c, 400, "Категория не указана")
		return
	}

	err, count, sch := repository.SelectProductsByCategory(c, category)
	if err != nil && err != sql.ErrNoRows {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	} else if count == 0 {
		service.ErrorResponse(c, 404, "Ни одного продукта по заданной категории не найдено")
		return
	}

	res := service.ConvertProducts(sch)

	c.JSON(200, res)
}

// @Tags Товары
// @Summary Получить товар по id
// @Description Возвращает подробное описание товара и его характеристики (если есть)
// @ID get_product_by_id
// @Accept json
// @Produce json
// @Param id path string true "id товара"
// @Success 200 {object} dto.ProductDetail
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /products/{id} [get]
func GetProductByID(c *gin.Context) {

	prod_id := c.Param("id")
	if prod_id == "" {
		service.ErrorResponse(c, 400, "Товар не указан")
		return
	}

	err, count, prod := repository.GetProductDetail(c, prod_id)
	if err != nil && err != sql.ErrNoRows {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	} else if count == 0 {
		service.ErrorResponse(c, 404, "Ни одного продукта по заданному id не найдено")
		return
	}
	err, count, imgs := repository.GetProductPictures(c, prod_id)

	img_arr := service.GetImages(imgs)
	res := service.ConvertProductDetail(prod, img_arr)

	c.JSON(200, res)
}

// @Tags Товары
// @Summary Получить похожие товары по тегу
// @Description Возвращает товары которые с таким же тегом, как и у входящего товара. Если таких нет, или их не 4 - добивает рандомными до 4ёх.
// @ID get_products_by_tag
// @Accept json
// @Produce json
// @Param id path string true "id товара"
// @Success 200 {object} dto.ProductsResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /products/{id}/relative [get]
func GetProductsByTag(c *gin.Context) {

	prod_id := c.Param("id")
	if prod_id == "" {
		service.ErrorResponse(c, 400, "Товар не указан")
		return
	}

	err, count, related := repository.GetProductsByTag(c, prod_id)
	if err != nil && err != sql.ErrNoRows {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	} else if count == 0 {
		service.ErrorResponse(c, 404, "Ни одного продукта по заданному id не найдено")
		return
	}

	res := service.ConvertProducts(related)

	c.JSON(200, res)
}
