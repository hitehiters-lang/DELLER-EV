package service

import (
	"deller/internal/dto"
	"fmt"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, err int, message string) {
	res := dto.ErrorMessage{
		Code:    err,
		Message: message,
	}
	fmt.Println("[SERVER] -", err, message)
	c.JSON(err, res)
}
