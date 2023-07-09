/*
Contains functions to set up the server and run it.
*/
package main

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

func main() {
	serverConfig := initialiseProdServer()
	serverConfig.setupRoutes()
	serverConfig.runRouter()
	log.Println("Setup complete!")
}

// serverConfig contains the essentials to run the backend - a router,
// a Redis database for fast reads, and a database for persistent data
type serverConfig struct {
	db          *gorm.DB
	store       redis.Store
	redis       *goredis.Client
	router      *gin.Engine
	GoogleCloud *storage.Client
}

// Returns a server configuration with the production database (as defined
// by environmental variables) set as the database
func initialiseProdServer() *serverConfig {
	router := gin.Default()
	db := database.ConnectProdDatabase()
	store := setupStore()
	redis := setupRedis()
	GoogleCloud := setupGoogleCloud()
	server := serverConfig{
		db:          db,
		router:      router,
		store:       store,
		redis:       redis,
		GoogleCloud: GoogleCloud,
	}
	return &server
}

// Sets up the Redis store from environmental variables
func setupStore() redis.Store {
	env := helpers.RetrieveRedisEnv()
	redisAddress := env.Address()
	const redisNetwork = "tcp"
	store, err := redis.NewStore(env.MaxConn, redisNetwork, redisAddress, "", []byte(env.Secret))
	if err != nil {
		log.Fatal(err.Error())
	}
	return store
}

func setupRedis() *goredis.Client {
	env := helpers.RetrieveRedisEnv()
	redisAddress := env.Address()
	const redisNetwork = "tcp"
	client := goredis.NewClient(&goredis.Options{
		Addr:     redisAddress,
		Network:  redisNetwork,
		Password: "",
	})
	return client
}

func setupGoogleCloud() *storage.Client {
	ctx := context.Background()
	env := helpers.RetrieveGoogleCloudEnv()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(env.Filepath))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func (server *serverConfig) runRouter() {
	env := helpers.RetrieveWebAppEnv()
	routerAddress := env.Address()
	err := server.router.Run(routerAddress)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
}
