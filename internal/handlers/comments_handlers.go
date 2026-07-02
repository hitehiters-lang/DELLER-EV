package handlers

import (
	"database/sql"
	"deller/internal/dto"
	"deller/internal/repository"
	"deller/internal/schema"
	"deller/internal/service"
	"errors"

	"github.com/gin-gonic/gin"
)

// @Tags Комментарии
// @Summary Получить комментарии поста
// @Description Возвращает все комментарии поста. parent_id это к какому комментарию был ответ.
// @ID get_comments
// @Accept json
// @Produce json
// @Param slug path string true "Слаг поста"
// @Success 200 {object} dto.CommentResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /blog/comments/{slug} [get]
func GetPostComments(c *gin.Context) {

	slug := c.Param("slug")
	if slug == "" {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	err, count, comments := repository.GetComments(c, slug)
	if err != nil && err != sql.ErrNoRows {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	} else if count == 0 {
		service.ErrorResponse(c, 404, "Ни одного комментария не найдено")
		return
	}

	res := service.ConvertComments(comments, count)

	c.JSON(200, res)
}

// @Tags Комментарии
// @Summary Написать комментарий к посту
// @Description Добавляеет комментарий к посту в БД. parent_id это к какому комментарию был ответ. Возвращает созданный комментарий в массиве
// @ID add_comment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param slug path string true "Слаг поста"
// @Param com body dto.NewComment true "Данные комментария"
// @Success 201 {object} dto.CommentResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /blog/comments/{slug} [post]
func AddPostComment(c *gin.Context) {

	slug := c.Param("slug")
	if slug == "" {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	var body dto.NewComment
	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	err := service.ValidateComBody(body)
	if err != nil {
		service.ErrorResponse(c, 400, err.Error())
		return
	}

	err = repository.CheckComParent(c, body.ParentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 404, "Комментарий для ответа не найден")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка проверки ответа на комментарий: "+err.Error())
		return
	}
	com, err := repository.AddComment(c, body, slug)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка отправки комментария: "+err.Error())
		return
	}

	rc := []schema.Comments{}
	rc = append(rc, com)

	res := service.ConvertComments(rc, 1)

	c.JSON(200, res)
}
