package handlers

import (
	"database/sql"
	"deller/internal/repository"
	"deller/internal/service"

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
