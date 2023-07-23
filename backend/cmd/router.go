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
	"github.com/redis/go-redis/v9"
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
	apiEnv := &controllers.APIEnv{
		DB:          s.db,
		GoogleCloud: s.GoogleCloud,
		NotifRedis:  s.notifRedis,
	}

	// Sets the ClientAddress and BackendAddress global variables in the models package so that the env file
	// does not need to be read everytime we require the client address or backend address
	helpers.SetModelClientAddress()
	helpers.SetModelBackendAddress()

	// Register routes - routes are grouped by features for greater
	// modularity
	setupPostAPI(routerGroup, apiEnv)
	setupUserAPI(routerGroup, apiEnv)
	setupAuthAPI(routerGroup, apiEnv)
	setupPhotoAPI(routerGroup, apiEnv)
	setupLikeAPI(routerGroup, apiEnv, s.likesRedis)
	setupCommentAPI(routerGroup, apiEnv, s.commentsRedis)
	setupNotificationAPI(routerGroup, apiEnv, s.notifRedis)
	setupCommunityAPI(routerGroup, apiEnv)
	setupProjectAPI(routerGroup, apiEnv)
	setupSearchAPI(routerGroup, apiEnv)
}

// Sets up CORS to allow the frontend app to access resources
func (s *serverConfig) configureCors() {
	// Get address of frontend app from environmental variables
	env := helpers.RetrieveClientEnv()
	localClientAddress := env.Host

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
	const postPathWithID = helpers.PostPath + "/:" + helpers.PostIDKey

	// Private routes
	rg.Private().GET(helpers.PostPath, api.GetPosts)
	rg.Private().GET(postPathWithID, api.GetPostByID)
	rg.Private().POST(helpers.PostPath, api.CreatePost)
	rg.Private().PATCH(postPathWithID, api.UpdatePost)
	rg.Private().DELETE(postPathWithID, api.DeletePost)
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
	const userPath = "/user"
	rg.Public().POST("/signup", api.CreateUser)

	rg.Private().GET("/users/:username", api.GetProfile)
	rg.Private().GET(userPath, api.GetSelfProfile)
	rg.Private().PATCH(userPath, api.UpdateUser)
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

func setupLikeAPI(rg RouterGrouper, api LikeAPIer, client *redis.Client) {
	api.InitialiseLikeHandler(client)
	registerLikeRoutes(rg, api)
}

type LikeAPIer interface {
	InitialiseLikeHandler(*redis.Client)
	PostLike(*gin.Context)
	DeleteLike(*gin.Context)
}

func registerLikeRoutes(rg RouterGrouper, api LikeAPIer) {
	const likePathWithID = "/likes/:" + helpers.PostIDKey

	rg.Private().POST(likePathWithID, api.PostLike)
	rg.Private().DELETE(likePathWithID, api.DeleteLike)
}

func setupCommentAPI(rg RouterGrouper, api CommentAPIer, client *redis.Client) {
	api.InitialiseCommentHandler(client)
	registerCommentRoutes(rg, api)
}

type CommentAPIer interface {
	InitialiseCommentHandler(*redis.Client)
	CreateComment(*gin.Context)
	// Generates comment feed
	GetComments(*gin.Context)
	UpdateComment(*gin.Context)
	DeleteComment(*gin.Context)
}

func registerCommentRoutes(rg RouterGrouper, api CommentAPIer) {
	const commentRouteWithID = helpers.CommentPath + "/:" + helpers.CommentIDKey

	// Private routes
	rg.Private().GET(helpers.CommentPath, api.GetComments)
	rg.Private().POST(helpers.CommentPath, api.CreateComment)
	rg.Private().PATCH(commentRouteWithID, api.UpdateComment)
	rg.Private().DELETE(commentRouteWithID, api.DeleteComment)
}

func setupNotificationAPI(rg RouterGrouper, api NotificationAPIer, client *redis.Client) {
	api.InitialiseNotificationHandler(client)
	registerNotificationRoutes(rg, api)
}

type NotificationAPIer interface {
	InitialiseNotificationHandler(*redis.Client)
	GetNotifications(*gin.Context)
}

func registerNotificationRoutes(rg RouterGrouper, api NotificationAPIer) {
	rg.Private().GET("/notifications", api.GetNotifications)
}

type CommunityAPIer interface {
	InitialiseCommunityHandler()
	CreateCommunity(*gin.Context)
	GetCommunities(*gin.Context)
	GetCommunityByName(*gin.Context)
	UpdateCommunity(*gin.Context)
}

func setupCommunityAPI(rg RouterGrouper, api CommunityAPIer) {
	api.InitialiseCommunityHandler()
	registerCommunityRoutes(rg, api)
}

func registerCommunityRoutes(rg RouterGrouper, api CommunityAPIer) {
	const communityPathWithName = helpers.CommunityPath + "/:" + helpers.CommunityNameKey
	rg.Private().GET(helpers.CommunityPath, api.GetCommunities)
	rg.Private().GET(communityPathWithName, api.GetCommunityByName)
	rg.Private().POST(helpers.CommunityPath, api.CreateCommunity)
	rg.Private().PATCH(communityPathWithName, api.UpdateCommunity)
}

type ProjectAPIer interface {
	InitialiseProjectHandler()
	CreateProject(*gin.Context)
	DeleteProject(*gin.Context)
	GetProjects(*gin.Context)
	GetProjectByID(*gin.Context)
	UpdateProject(*gin.Context)
}

func setupProjectAPI(rg RouterGrouper, api ProjectAPIer) {
	api.InitialiseProjectHandler()
	registerProjectRoutes(rg, api)
}

func registerProjectRoutes(rg RouterGrouper, api ProjectAPIer) {
	const projectPathWithID = helpers.ProjectPath + "/:" + helpers.ProjectIDKey
	rg.Private().GET(helpers.ProjectPath, api.GetProjects)
	rg.Private().GET(projectPathWithID, api.GetProjectByID)
	rg.Private().POST(helpers.ProjectPath, api.CreateProject)
	rg.Private().DELETE(projectPathWithID, api.DeleteProject)
	rg.Private().PATCH(projectPathWithID, api.UpdateProject)
}

func setupSearchAPI(rg RouterGrouper, api SearchAPIer) {
	registerSearchRoutes(rg, api)
}

type SearchAPIer interface {
	GetSearchResults(*gin.Context)
}

func registerSearchRoutes(rg RouterGrouper, api SearchAPIer) {
	rg.Private().GET("/search", api.GetSearchResults)
}
