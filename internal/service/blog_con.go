package service

import (
	"deller/internal/dto"
	"deller/internal/schema"
	"errors"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func ValidatePostBody(body dto.NewPost) error {
	var ret_err string
	var err error

	if utf8.RuneCountInString(body.Content) < 4 {
		ret_err = "Содержание поста слишком короткое"
	} else if utf8.RuneCountInString(body.Content) > 5000 {
		ret_err = "Содержание поста слишком длинное"
	}

	if utf8.RuneCountInString(body.Title) < 10 {
		ret_err = "Название слишком короткое"
	} else if utf8.RuneCountInString(body.Title) > 500 {
		ret_err = "Название слишком длинное"
	}

	if ret_err != "" {
		err = errors.New(ret_err)
	}

	return err
}

func CreateSlugFromTitle(t string) string {
	s := slug.Make(t)
	r := []rune(s)
	if len(r) > 80 {
		s = strings.TrimSuffix(string(r[:80]), "-")
	}
	return s
}

func GetPostsQuery(c *gin.Context) (int, int, error) {

	off := 0
	lim := 15
	if c.Query("off") != "" {
		offset, err := strconv.Atoi(c.Query("off"))
		if err != nil {
			return 0, 0, err
		}
		off = offset
	}
	if c.Query("lim") != "" {
		limit, err := strconv.Atoi(c.Query("lim"))
		if err != nil {
			return 0, 0, err
		}
		lim = limit
	}
	if lim > 30 {
		lim = 30
	} else if lim < 1 {
		lim = 15
	}
	if off < 0 {
		off = 0
	}
	return lim, off, nil
}

func ConvertPosts(c int, s []schema.Posts) dto.PostResponse {
	var po []dto.Post
	for _, e := range s {
		sts := dto.Post{
			Slug:          e.Slug,
			Title:         e.Title,
			AuthorName:    e.Author.Name,
			AuthorIsAdmin: e.Author.IsAdmin,
			IsApproved:    e.IsApproved,
			CreatedAt:     e.CreatedAt,
		}
		po = append(po, sts)
	}
	return dto.PostResponse{
		Posts: po,
		Total: c,
	}
}

func PostDetails(p schema.Posts) dto.PostDetail {
	dto_post := dto.PostDetail{
		Slug:          p.Slug,
		Title:         p.Title,
		Content:       p.Content,
		AuthorName:    p.Author.Name,
		AuthorIsAdmin: p.Author.IsAdmin,
		IsApproved:    p.IsApproved,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.CreatedAt,
	}
	return dto_post
}

func ConvertComments(c []schema.Comments, t int) dto.CommentResponse {
	var com []dto.Comment
	for _, e := range c {
		ment := dto.Comment{
			Id:                 e.Id,
			Content:            e.Content,
			AuthorName:         e.Author.Name,
			AuthorIsAdmin:      e.Author.IsAdmin,
			AuthorIsPostAuthor: false,
			CreatedAt:          time.Now(),
		}
		if e.Parent != nil {
			ment.ParentId = e.ParentId
		}
		if e.Post.AuthorID == e.AuthorId {
			ment.AuthorIsPostAuthor = true
		}
		com = append(com, ment)
	}
	return dto.CommentResponse{
		Comments: com,
		Total:    t,
	}
}

func ValidateComBody(body dto.NewComment) error {
	var ret_err string
	var err error

	if utf8.RuneCountInString(body.Content) > 999 {
		ret_err = "Содержание комментария слишком длинное"
	}

	if ret_err != "" {
		err = errors.New(ret_err)
	}

	return err
}
