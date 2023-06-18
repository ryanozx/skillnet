/*
Contains functions to set up the router and register various routes. Modify this
file to add additional routes where necessary.
*/
package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/controllers"
	"github.com/ryanozx/skillnet/helpers"
	"github.com/ryanozx/skillnet/middleware"
)

// Sets up the routes on the server router
func (s *serverConfig) setupRoutes() {
	log.Println("Setting up routes...")
	s.configureCors()
	s.router.Use(sessions.Sessions("skillnet", s.store))

	routerGroup := s.RouterGroups()
	apiEnv := controllers.CreateAPIEnv(s.db, s.GoogleCloud)

	// Register routes - routes are grouped by features for greater
	// modularity
	setupPostAPI(routerGroup, apiEnv)
	setupUserAPI(routerGroup, apiEnv)
	setupAuthAPI(routerGroup, apiEnv)
	setupPhotoAPI(routerGroup, apiEnv)
}

// Sets up CORS to allow the frontend app to access resources
func (s *serverConfig) configureCors() {
	// Get address of frontend app from environmental variables
	env := helpers.RetrieveClientEnv()
	localClientAddress := env.Address()

	// Set up configuration and apply it to the router
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{localClientAddress}
	corsConfig.AllowCredentials = true
	s.router.Use(cors.New(corsConfig))
}

// Initialises a RouterGroups object from the server router
func (s *serverConfig) RouterGroups() *RouterGroups {
	// All publicly accessible routes are prefixed with "/",
	// while all non-publicly accessible routes are prefixed with "/auth"
	publicGroup := s.router.Group("/")
	privateGroup := s.router.Group("/auth")

	// add middleware - for instance, middleware to check that the user
	// has a valid session in order to access non-publicly accessible routes
	privateGroup.Use(middleware.AuthRequired)

	routerGroup := RouterGroups{
		public:  publicGroup,
		private: privateGroup,
	}
	return &routerGroup
}

type RouterGroups struct {
	public  *gin.RouterGroup
	private *gin.RouterGroup
}

func (rg *RouterGroups) Public() *gin.RouterGroup {
	return rg.public
}

func (rg *RouterGroups) Private() *gin.RouterGroup {
	return rg.private
}

// At present we use only two router groups as only the AuthRequired middleware is
// used, hence there needs to be a distinction between public and private routes.
// Should any subset of routes require additional middleware, the router groups can
// be added
type RouterGrouper interface {
	Public() *gin.RouterGroup
	Private() *gin.RouterGroup
}

// Sets up Post API
func setupPostAPI(rg RouterGrouper, api PostAPIer) {
	api.InitialisePostHandler()
	registerPostRoutes(rg, api)
}

// PostAPIer is an interface that describes the methods required to implement
// CRUD for Posts
type PostAPIer interface {
	InitialisePostHandler()
	// Generates post feed
	GetPosts(*gin.Context)
	// Returns a specific post
	GetPostByID(*gin.Context)
	CreatePost(*gin.Context)
	UpdatePost(*gin.Context)
	DeletePost(*gin.Context)
}

func registerPostRoutes(rg RouterGrouper, api PostAPIer) {
	// Public routes
	rg.Public().GET("/posts", api.GetPosts)
	rg.Public().GET("/posts/:id", api.GetPostByID)

	// Private routes
	rg.Private().POST("/posts", api.CreatePost)
	rg.Private().PATCH("/posts/:id", api.UpdatePost)
	rg.Private().DELETE("/posts/:id", api.DeletePost)
}

// Sets up User API
func setupUserAPI(rg RouterGrouper, api UserAPIer) {
	api.InitialiseUserHandler()
	registerUserRoutes(rg, api)
}

// UserAPIer is an interface that describes the methods required to implement
// CRUD for Users
type UserAPIer interface {
	InitialiseUserHandler()
	// Returns user profile as seen by visitor
	GetProfile(*gin.Context)
	// Returns own user profile with private information
	GetSelfProfile(*gin.Context)
	CreateUser(*gin.Context)
	UpdateUser(*gin.Context)
}

func registerUserRoutes(rg RouterGrouper, api UserAPIer) {
	rg.Public().GET("/users/:username", api.GetProfile)
	rg.Public().POST("/signup", api.CreateUser)

	rg.Private().GET("/user", api.GetSelfProfile)
	rg.Private().PATCH("/user", api.UpdateUser)
}

// Sets up Auth API
func setupAuthAPI(rg RouterGrouper, api AuthAPIer) {
	api.InitialiseAuthHandler()
	registerAuthRoutes(rg, api)
}

// AuthAPIer is an interface that describes the methods required to implement
// authentication
type AuthAPIer interface {
	InitialiseAuthHandler()
	GetLogin(*gin.Context)
	PostLogin(*gin.Context)
	PostLogout(*gin.Context)
}

func registerAuthRoutes(rg RouterGrouper, api AuthAPIer) {
	rg.Public().GET("/login", api.GetLogin)
	rg.Public().POST("/login", api.PostLogin)

	rg.Private().POST("/logout", api.PostLogout)
}

func setupPhotoAPI(rg RouterGrouper, api PhotoAPIer) {
	// api.InitialisePhotoHandler()
	registerPhotoRoutes(rg, api)
}

type PhotoAPIer interface {
	PostUserPicture(*gin.Context)
}

func registerPhotoRoutes(rg RouterGrouper, api PhotoAPIer) {
	// rg.Public().POST("/testupload", api.PostUserPicture)

	rg.Private().POST("/user/photo", api.PostUserPicture)
}
