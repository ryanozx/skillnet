package main

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/controllers"
	"github.com/ryanozx/skillnet/middleware"
)

func registerRoutes(router *gin.Engine, store redis.Store) {
	router.Use(sessions.Sessions("mysession", store))

	public := router.Group("/")
	registerPublicRoutes(public)

	auth := router.Group("/auth")
	auth.Use(middleware.AuthRequired)
	registerPrivateRoutes(auth)
}

func registerPublicRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/login", controllers.LoginPostHandler())
	routerGroup.GET("/login", controllers.LoginGetHandler())
	routerGroup.GET("/posts", controllers.GetPosts)
	routerGroup.GET("/posts/:id", controllers.GetPostByID)
}

func registerPrivateRoutes(routerGroup *gin.RouterGroup) {
	routerGroup.POST("/posts", controllers.PostPost)
	routerGroup.PATCH("/posts/:id", controllers.UpdatePost)
	routerGroup.DELETE("/posts/:id", controllers.DeletePost)
	routerGroup.GET("/test", func(context *gin.Context) {
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "authorised",
		})
	})

	routerGroup.GET("/logout", controllers.LogoutGetHandler())
}
