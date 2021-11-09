package main

import (
	"os"

	controller "personal-blog/controller"
	middleware "personal-blog/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/login", controller.Login)
	router.POST("/posts", middleware.TokenAuthMiddleware(), controller.CreatePost)
	router.PUT("/posts", middleware.TokenAuthMiddleware(), controller.UpdatePost)
	router.GET("/posts/id/:id", controller.GetPostById)
	router.GET("/posts", controller.GetPosts)
	router.GET("/posts/categoryid/:categoryid", controller.GetPostsByCategoryId)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)

}
