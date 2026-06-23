package main

import (
	"deller/docs"
	"deller/internal/config"
	"deller/internal/database"
	"deller/internal/handlers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @bearerFormat JWT
// @description Введите JWT токен

// @title Deller-Ev API
// @version 0.1
// @description Документация к API.
// @BasePath /
func main() {
	config.LoadENV()
	database.InitDatabase()

	r := gin.Default()
	r.Use(cors.Default())

	api := r.Group("/api/v1")
	{
		api.GET("/products/list", handlers.GetProducts)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"

	api.GET("/docs/*any", func(c *gin.Context) {
		if c.Param("any") == "" || c.Param("any") == "/" {
			c.Redirect(http.StatusMovedPermanently, "/api/v1/docs/index.html")
			return
		}
		ginSwagger.WrapHandler(swaggerfiles.Handler)(c)
	})

	port := config.GetServer()
	server := "0.0.0.0:" + port
	r.Run(server)
}
