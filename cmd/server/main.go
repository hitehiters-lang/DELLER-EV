package main

import (
	"deller/docs"
	"deller/internal/config"
	"deller/internal/database"
	"deller/internal/handlers"
	pkgauth "deller/internal/pkg/auth"
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
// @version 0.6
// @description Документация к API.
// @BasePath /
func main() {
	config.LoadENV()
	database.InitDatabase()

	r := gin.Default()
	r.Use(cors.Default())

	api := r.Group("/api/v1")
	{
		api.Static("/images", "./shared")

		api.GET("/categories", handlers.GetCategories)
		api.GET("/categories/:id", handlers.GetProducts)

		api.GET("/products/:id", handlers.GetProductByID)
		api.GET("/products/:id/relative", handlers.GetProductsByTag)

		api.POST("/auth/login", handlers.LogIn)
		api.POST("/auth/reg", handlers.Register)

		api.GET("/blog/posts", handlers.GetPosts)
		api.GET("/blog/posts/:slug", handlers.GetPostDetail)
		api.GET("/blog/comments/:slug", handlers.GetPostComments)

		api.GET("/products/:id/reviews", handlers.GetProductReviews)

		auth := api.Group("")
		auth.Use(pkgauth.JWTAuthMiddleware())
		{
			auth.GET("/auth/me", handlers.Me)

			auth.POST("/blog/posts", handlers.NewPost)
			auth.POST("/blog/comments/:slug", handlers.AddPostComment)

			auth.POST("/admin/images", handlers.AddPicture)
			auth.POST("/admin/image-bind", handlers.BindPicture)
			auth.DELETE("/admin/image-unbind", handlers.UnbindPicture)

			auth.POST("/admin/categories", handlers.AddCategory)
			auth.DELETE("/admin/categories/:id", handlers.DelCategory)
			auth.PUT("/admin/categories/:id", handlers.PutCategory)

			auth.POST("/admin/products", handlers.AddProduct)
			auth.DELETE("/admin/products/:id", handlers.DelProduct)
			auth.PUT("/admin/products/:id", handlers.PutProduct)

			auth.DELETE("/admin/comments/:id", handlers.DelPostComment)
			auth.DELETE("/admin/posts/:slug", handlers.DelPost)
			auth.DELETE("/admin/reviews/:id", handlers.DelProductReview)

			auth.POST("/products/:id/reviews", handlers.AddProductReview)
		}
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
