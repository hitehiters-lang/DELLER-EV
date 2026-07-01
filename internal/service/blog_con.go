package service

import (
	"deller/internal/dto"
	"errors"
	"strings"
	"unicode/utf8"

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
