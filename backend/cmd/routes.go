package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/controllers"
	"github.com/ryanozx/skillnet/middleware"
)

func (server *serverConfig) setupRoutes() {
	// TODO: Research implications of allowing cross-domain requests
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	server.router.Use(cors.New(corsConfig))
	server.router.Use(sessions.Sessions("mysession", server.store))

	apiEnv := &controllers.APIEnv{
		DB: server.db,
	}
	public := server.router.Group("/")
	registerPublicRoutes(public, apiEnv)

	auth := server.router.Group("/auth")
	auth.Use(middleware.AuthRequired)
	registerPrivateRoutes(auth, apiEnv)
}

func registerPublicRoutes(routerGroup *gin.RouterGroup, api *controllers.APIEnv) {
	routerGroup.POST("/login", controllers.PostLogin)
	routerGroup.GET("/login", controllers.GetLogin)
	routerGroup.GET("/posts", api.GetPosts)
	routerGroup.GET("/posts/:id", api.GetPostByID)
}

func registerPrivateRoutes(routerGroup *gin.RouterGroup, api *controllers.APIEnv) {
	routerGroup.POST("/posts", api.CreatePost)
	routerGroup.PATCH("/posts/:id", api.UpdatePost)
	routerGroup.DELETE("/posts/:id", api.DeletePost)
	routerGroup.GET("/test", func(context *gin.Context) {
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "authorised",
		})
	})

	routerGroup.GET("/logout", controllers.GetLogout)
}
