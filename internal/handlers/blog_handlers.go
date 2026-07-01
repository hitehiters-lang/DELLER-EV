package handlers

import (
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/repository"
	"deller/internal/service"

	"github.com/gin-gonic/gin"
)

// @Tags Блог
// @Summary Создать новый пост
// @Description Создаёт новый пост в БД. Доступно только админам.
// @ID new_post
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param login body dto.NewPost true "Данные поста"
// @Success 200 {object} dto.SlugReturn
// @Failure 400 {object} dto.ErrorMessage
// @Failure 409 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /blog/posts [post]
func NewPost(c *gin.Context) {

	r, _ := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	var body dto.NewPost

	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	err := service.ValidatePostBody(body)
	if err != nil {
		service.ErrorResponse(c, 400, err.Error())
		return
	}

	new_slug := service.CreateSlugFromTitle(body.Title)
	add_post, err := repository.AddPost(c, body, new_slug)

	c.JSON(201, dto.SlugReturn{Slug: add_post.Slug})
}
