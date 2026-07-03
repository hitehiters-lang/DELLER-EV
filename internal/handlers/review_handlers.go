package handlers

import (
	"database/sql"
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/repository"
	"deller/internal/schema"
	"deller/internal/service"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Tags Отзывы
// @Summary Получить отзывы о товаре
// @Description Возвращает все отзывы о товаре. parent_id это к какому отзыву был ответ. У ответов на отзывы не должна отображаться оценка
// @ID get_reviews
// @Accept json
// @Produce json
// @Param id path string true "id товара"
// @Success 200 {object} dto.ReviewResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /products/{id}/reviews [get]
func GetProductReviews(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	err, count, reviews := repository.GetReviews(c, id)
	if err != nil && err != sql.ErrNoRows {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	} else if count == 0 {
		service.ErrorResponse(c, 404, "Ни одного отзыва не найдено")
		return
	}

	res := service.ConvertReviews(reviews)

	c.JSON(200, res)
}

// @Tags Отзывы
// @Summary Написать отзыв к товару
// @Description Добавляет отзыв к товару в БД. parent_id это к какому отзыву был ответ. Возвращает созданный отзыв в массиве.
// @ID add_review
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param parent_id query string false "uuid отзыва для ответа (пусто - комментарий к посту). Ответы на отзывы доступны только админам. При ответе на отзыв, оценка игнорируется, её на фронт не пихаем (ставь 1 балл у ответов всегда)"
// @Param id path string true "id товара"
// @Param rev body dto.NewReview true "Данные отзыва"
// @Success 201 {object} dto.ReviewResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /products/{id}/reviews [post]
func AddProductReview(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	parent_uuid := c.Query("parent_id")
	is_answer := false
	parent, err := uuid.Parse(parent_uuid)
	if err == nil {
		is_answer = true

		rights, err := pkgauth.GetCurrentRights(c)
		if err != nil || (rights != "a" && rights != "r") {
			service.ErrorResponse(c, 403, "Недостаточно прав для ответа: Отвечать на отзывы может только админ")
			return
		}
	}

	var body dto.NewReview
	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	err = service.ValidateRevBody(body)
	if err != nil {
		service.ErrorResponse(c, 400, err.Error())
		return
	}

	err = repository.CheckProduct(c, id)
	if err != nil {
		service.ErrorResponse(c, 404, "Товар не найден")
		return
	}

	var rev schema.Reviews
	if is_answer {
		err = repository.CheckReview(c, parent)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				service.ErrorResponse(c, 404, "Отзыв для ответа не найден")
				return
			}
			service.ErrorResponse(c, 500, "Ошибка проверки ответа на отзыв: "+err.Error())
			return
		}
	} else {
		bl, err := repository.CheckForRevRepeat(c, id)
		if err != nil {
			service.ErrorResponse(c, 500, err.Error())
			return
		}
		if bl {
			service.ErrorResponse(c, 409, "Ваш отзыв о товаре уже существует")
			return
		}
	}
	rev, err = repository.AddReview(c, body, id, parent)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка отправки отзыва: "+err.Error())
		return
	}

	rr := []schema.Reviews{}
	rr = append(rr, rev)

	res := service.ConvertReviews(rr)

	c.JSON(200, res)
}

// @Tags Отзывы
// @Summary Удалить отзыв
// @Description Удаляет отзыв по id. Доступно только админам.
// @ID del_review
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "id отзыва"
// @Success 201 {object} dto.Message
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/reviews/{id} [delete]
func DelProductReview(c *gin.Context) {

	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	rev_uuid := c.Param("id")
	if rev_uuid == "" {
		service.ErrorResponse(c, 400, "Отзыв не указан")
		return
	}
	rev := uuid.MustParse(rev_uuid)

	err = repository.CheckReview(c, rev)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 404, "Отзыва не существует")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка проверки отзыва: "+err.Error())
		return
	}

	err = repository.DeleteReview(c, rev)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка удаления: "+err.Error())
		return
	}
	c.JSON(200, dto.Message{
		Code:    200,
		Message: "Отзыв успешно удалён",
	})
}
