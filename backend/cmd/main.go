package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/controllers"
	"github.com/ryanozx/skillnet/database"
)

func main() {
	router := gin.Default()
	database.ConnectDatabase()

	router.GET("/posts", controllers.GetPosts)
	router.POST("/posts", controllers.PostPost)
	router.GET("/posts/:id", controllers.GetPostByID)
	router.PATCH("/posts/:id", controllers.UpdatePost)
	router.DELETE("/posts/:id", controllers.DeletePost)

	routerErr := router.Run("localhost:8080")
	if routerErr != nil {
		return
	}
}
