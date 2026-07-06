package handlers

import (
	"database/sql"
	"deller/internal/dto"
	pkgauth "deller/internal/pkg/auth"
	"deller/internal/repository"
	"deller/internal/service"
	"errors"

	"github.com/gin-gonic/gin"
)

// @Tags Блог
// @Summary Создать новый пост
// @Description Создаёт новый пост в БД. Доступно только админам. Картинки добавляются через роут /admin/images. Название картинки пишется без расширения.
// @ID new_post
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param login body dto.NewPost true "Данные поста"
// @Success 201 {object} dto.SlugReturn
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

// @Tags Блог
// @Summary Получить посты
// @Description Возвращает самые новые посты с пагинацией, а также сколько всего постов есть в БД на данный момент.
// @ID get_posts
// @Accept json
// @Produce json
// @Param off query int false "Смещение (def - 0)"
// @Param lim query int false "Лимит (def - 15)"
// @Success 200 {object} dto.PostResponse
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /blog/posts [get]
func GetPosts(c *gin.Context) {

	lim, off, err := service.GetPostsQuery(c)
	if err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
	}

	err, count, posts := repository.GetPosts(c, lim, off)
	if err != nil && err != sql.ErrNoRows {
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	} else if count == 0 {
		service.ErrorResponse(c, 404, "Ни одного поста не найдено")
		return
	}

	res := service.ConvertPosts(count, posts)

	c.JSON(200, res)
}

// @Tags Блог
// @Summary Получить пост по слагу
// @Description Возвращает полный пост.
// @ID get_post_detail
// @Accept json
// @Produce json
// @Param slug path string true "Слаг поста"
// @Success 200 {object} dto.PostDetail
// @Failure 400 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /blog/posts/{slug} [get]
func GetPostDetail(c *gin.Context) {

	slug := c.Param("slug")
	if slug == "" {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	post, err := repository.GetPostDetail(c, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			service.ErrorResponse(c, 404, "Ни одного поста не найдено")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}

	res := service.PostDetails(post)

	c.JSON(200, res)
}

// @Tags Блог
// @Summary Удалить пост
// @Description Удаляет пост по слагу. Доступно только админам
// @ID del_post
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param slug path string true "Слаг поста"
// @Success 201 {object} dto.Message
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/posts/{slug} [delete]
func DelPost(c *gin.Context) {

	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	slug := c.Param("slug")
	if slug == "" {
		service.ErrorResponse(c, 400, "Пост не указан")
		return
	}

	_, err = repository.GetPostDetail(c, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 404, "Поста не существует")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка проверки поста: "+err.Error())
		return
	}

	err = repository.DeletePost(c, slug)
	if err != nil {
		service.ErrorResponse(c, 500, "Ошибка удаления: "+err.Error())
		return
	}
	c.JSON(200, dto.Message{
		Code:    200,
		Message: "Пост успешно удалён",
	})
}

// @Tags Блог
// @Summary Изменить пост
// @Description Изменяет пост по слагу. Доступно только админам. КОДЕКСУ: Дополнительно добавь чтобы можно было вставлять изображения для блогов, как ты это сделал в товарах.
// @ID put_post
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param slug path string true "Слаг поста"
// @Param putter body dto.PutPostBody true "Тело запроса"
// @Success 201 {object} dto.Message
// @Failure 400 {object} dto.ErrorMessage
// @Failure 403 {object} dto.ErrorMessage
// @Failure 404 {object} dto.ErrorMessage
// @Failure 500 {object} dto.ErrorMessage
// @Router /admin/posts/{slug} [put]
func PutPost(c *gin.Context) {

	r, err := pkgauth.GetCurrentRights(c)
	if r != "a" && r != "r" {
		service.ErrorResponse(c, 403, "Недостаточно прав")
		return
	}

	var body dto.PutPostBody
	if err := c.ShouldBindJSON(&body); err != nil {
		service.ErrorResponse(c, 400, "Плохой запрос")
		return
	}

	slug := c.Param("slug")
	if slug == "" {
		service.ErrorResponse(c, 400, "Пост не указан")
		return
	}

	_, err = repository.GetPostDetail(c, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			service.ErrorResponse(c, 404, "Поста не существует")
			return
		}
		service.ErrorResponse(c, 500, "Ошибка сервера: "+err.Error())
		return
	}

	err = repository.PutPost(c, body, slug)
	if err != nil {
		service.ErrorResponse(c, 400, "Ошибка: "+err.Error())
		return
	}
	c.JSON(200, dto.Message{
		Code:    200,
		Message: "Пост " + slug + " успешно изменён",
	})
}
