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
	server.router.Use(sessions.Sessions("skillnet", server.store))

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
	routerGroup.POST("/login", api.PostLogin)
	routerGroup.GET("/login", api.GetLogin)
	routerGroup.GET("/posts", api.GetPosts)
	routerGroup.GET("/posts/:id", api.GetPostByID)
	routerGroup.POST("/signup", api.CreateUser)
	routerGroup.GET("/users/:username", api.GetProfile)
}

func registerPrivateRoutes(routerGroup *gin.RouterGroup, api *controllers.APIEnv) {
	routerGroup.POST("/posts", api.CreatePost)
	routerGroup.PATCH("/posts/:id", api.UpdatePost)
	routerGroup.DELETE("/posts/:id", api.DeletePost)
	routerGroup.PATCH("/user", api.UpdateUser)
	routerGroup.GET("/test", func(context *gin.Context) {
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "authorised",
		})
	})

	routerGroup.GET("/logout", api.GetLogout)
}
