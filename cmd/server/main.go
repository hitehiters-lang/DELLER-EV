package main

import (
	"deller/docs"
	"deller/internal/config"
	"deller/internal/database"
	"deller/internal/handlers"
	pkgauth "deller/internal/pkg/auth"
	pkgblog "deller/internal/pkg/blog_gen"
	pkgtrans "deller/internal/pkg/translator"
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
// @version 1.3
// @description Документация к API.
// @description
// @description ## Обновление API 1.3; Изменения:
// @description
// @description
// @description
// @description ## Обновление API 1.2; Изменения:
// @description
// @description ### 1. Роуты подкатегорий.
// @description Подкатегория - новая таблица в бд и поля в товарах.
// @description Она привязывается к Категории по FK, а товар привязывается и к категории и к подкатегории.
// @description Категория у товара - обязательный параметр, а вот подкатегория необязательна и может быть равна ничему ("")
// @description Поиск по Категории всё ещё актуален и работает, но для поиска конкретно по подкатегориям нужно использовать новые роуты.
// @description
// @description ### 2. Подкатегории в товарах.
// @description Теперь при возвращении товара возвращается и его подкатегория. Её не стоит отображать в карточке,
// @description так как она может быть null, "" или вообще не возвращаться на роуте (тоже учитывать)
// @description В админских роутах создания и изменения товара тоже было добавлено поле subcategory. Он не обязателен, если не выбран то должен отправляться на роут как "".
// @description При изменении товара, и категория и подкатегория если не изменяются пользователем - должны отправляться на роут те же, что и пришли.
// @description Так же и с полем category_id у подкатегории.
// @BasePath /
func main() {
	config.LoadENV()
	database.InitDatabase()

	err := pkgtrans.TranslatorCheck()
	if err != nil {
		return
	}

	pkgblog.Start()

	r := gin.Default()
	r.Use(cors.Default())
	r.MaxMultipartMemory = 8 << 20

	api := r.Group("/api/v1")
	{
		api.Static("/images", "./shared/images")
		api.Static("/files", "./shared/files")

		api.GET("/manuals", handlers.GetManuals)

		api.GET("/categories", handlers.GetCategories)
		api.GET("/categories/:id", handlers.GetProducts)

		api.GET("/subcategories/:id/categories", handlers.GetSubcategories)
		api.GET("/subcategories/:id/products", handlers.GetProductsSub)

		api.GET("/products/:id", handlers.GetProductByID)
		api.GET("/products/:id/relative", handlers.GetProductsByTag)

		api.POST("/auth/login", handlers.LogIn)
		api.POST("/auth/reg", handlers.Register)

		api.GET("/blog/posts", handlers.GetPosts)
		api.GET("/blog/posts/:slug", handlers.GetPostDetail)
		api.GET("/blog/comments/:slug", handlers.GetPostComments)

		api.GET("/products/:id/reviews", handlers.GetProductReviews)

		api.POST("/orders", handlers.AddOrder)
		api.POST("/contact", handlers.Contact)

		auth := api.Group("")
		auth.Use(pkgauth.JWTAuthMiddleware())
		{
			auth.GET("/auth/me", handlers.Me)

			auth.POST("/blog/posts", handlers.NewPost)
			auth.POST("/blog/comments/:slug", handlers.AddPostComment)

			auth.POST("/admin/files", handlers.AddFile)

			auth.POST("/admin/images", handlers.AddPicture)
			auth.POST("/admin/image-bind", handlers.BindPicture)
			auth.DELETE("/admin/image-unbind", handlers.UnbindPicture)

			auth.POST("/admin/subcategories", handlers.AddSubcategory)
			auth.DELETE("/admin/subcategories/:id", handlers.DelSubcategory)
			auth.PUT("/admin/subcategories/:id", handlers.PutSubcategory)

			auth.POST("/admin/categories", handlers.AddCategory)
			auth.DELETE("/admin/categories/:id", handlers.DelCategory)
			auth.PUT("/admin/categories/:id", handlers.PutCategory)

			auth.POST("/admin/products", handlers.AddProduct)
			auth.DELETE("/admin/products/:id", handlers.DelProduct)
			auth.PUT("/admin/products/:id", handlers.PutProduct)

			auth.DELETE("/admin/comments/:id", handlers.DelPostComment)
			auth.DELETE("/admin/posts/:slug", handlers.DelPost)
			auth.PUT("/admin/posts/:slug", handlers.PutPost)
			auth.DELETE("/admin/reviews/:id", handlers.DelProductReview)

			auth.DELETE("/admin/orders/:id", handlers.DelOrder)
			auth.GET("/admin/orders", handlers.GetOrders)

			auth.POST("/products/:id/reviews", handlers.AddProductReview)

			auth.POST("/admin/manuals", handlers.AddManual)
			auth.DELETE("/admin/manuals/:slug", handlers.DelManual)

			auth.POST("/admin/xlsx-import", handlers.ImportPricesFromXlsx)
			auth.GET("/admin/xlsx-export", handlers.ExportPricesToXlsx)

			auth.POST("/admin/cache", handlers.CacheAll)

			auth.GET("/admin/contact", handlers.GetMails)
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
