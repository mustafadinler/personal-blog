package main

import (
	"os"

	controller "personal-blog/controller"
	middleware "personal-blog/middleware"

	docs "personal-blog/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		post := v1.Group("/posts")
		{
			post.POST("", middleware.TokenAuthMiddleware(), controller.CreatePost)
			post.PUT("", middleware.TokenAuthMiddleware(), controller.UpdatePost)
			post.GET("/id/:id", controller.GetPostById)
			post.GET("", controller.GetPosts)
			post.GET("/categoryid/:categoryid", controller.GetPostsByCategoryId)
		}

		login := v1.Group("/login")
		{
			login.POST("/login", controller.Login)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":" + port)

}
