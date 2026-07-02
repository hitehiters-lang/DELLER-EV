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
// @Description Добавляет комментарий к посту в БД. parent_id это к какому комментарию был ответ. Возвращает созданный комментарий в массиве
// @ID add_comment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param parent_id query string false "uuid комментария для ответа (пусто - комментарий к посту)"
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

	parent_uuid := c.Query("parent_id")
	if slug == "" {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}
	parent, err := uuid.Parse(parent_uuid)
	is_answer := false
	if err == nil {
		is_answer = true
	}

	var body dto.NewComment
	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	err = service.ValidateComBody(body)
	if err != nil {
		service.ErrorResponse(c, 400, err.Error())
		return
	}

	var com schema.Comments
	if is_answer {
		err = repository.CheckComment(c, parent)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				service.ErrorResponse(c, 404, "Комментарий для ответа не найден")
				return
			}
			service.ErrorResponse(c, 500, "Ошибка проверки ответа на комментарий: "+err.Error())
			return
		}
	}
	com, err = repository.AddComment(c, body, slug, parent)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка отправки комментария: "+err.Error())
		return
	}

	rc := []schema.Comments{}
	rc = append(rc, com)

	res := service.ConvertComments(rc, 1)

	c.JSON(200, res)
}

// @Tags Комментарии
// @Summary Удалить комментарий
// @Description Удаляет комментарий по id. Доступно только админам
// @ID del_comment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "id комментария"
// @Success 201 {object} dto.Message
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/comments/{id} [delete]
func DelPostComment(c *gin.Context) {

	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	comm_uuid := c.Param("id")
	if comm_uuid == "" {
		service.ErrorResponse(c, 400, "Комментарий не указан")
		return
	}
	comm := uuid.MustParse(comm_uuid)

	err = repository.CheckComment(c, comm)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 404, "Комментария не существует")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка проверки комментария: "+err.Error())
		return
	}

	err = repository.DeleteComment(c, comm)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка удаления: "+err.Error())
		return
	}
	c.JSON(200, dto.Message{
		Code:    200,
		Message: "Комментарий успешно удалён",
	})
}
