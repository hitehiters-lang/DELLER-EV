package handlers

import (
	"deller/internal/repository"
	"deller/internal/service"

	"github.com/gin-gonic/gin"
)

// @Tags Товары
// @Summary Получить все товары определённой категории
// @ID get_products
// @Accept json
// @Produce json
// @Param category query string true "uuid категории" format(uuid)
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Router /products/list [get]
func GetProducts(c *gin.Context) {

	category_str := c.Query("category")

	err, count, sch := repository.SelectProductsByCategory(c, category_str)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	} else if count == 0 {
		service.ErrorResponse(c, 404, "Ни одного продукта по заданной категории не найдено")
		return
	}

	res := service.ConvertProducts(sch)

	c.JSON(200, res)
}
