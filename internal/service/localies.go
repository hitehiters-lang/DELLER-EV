package service

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func LoadLocalisation(c *gin.Context) string {
	l := strings.ToLower(c.Query("l"))
	switch l {
	case "ru", "en", "kk", "zh":
		return l
	}
	return ""
}
