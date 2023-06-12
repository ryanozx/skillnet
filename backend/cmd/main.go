/*
Contains functions to set up the server and run it.
*/
package main

import (
	"log"

	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
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
	db     *gorm.DB
	router *gin.Engine
	store  redis.Store
}

// Returns a server configuration with the production database (as defined
// by environmental variables) set as the database
func initialiseProdServer() *serverConfig {
	router := gin.Default()
	db := database.ConnectProdDatabase()
	store := setupRedis()
	server := serverConfig{
		db:     db,
		router: router,
		store:  store,
	}
	return &server
}

// Sets up the Redis store from environmental variables
func setupRedis() redis.Store {
	env := helpers.RetrieveRedisEnv()
	redisAddress := env.Address()
	const redisNetwork = "tcp"
	store, err := redis.NewStore(env.MaxConn, redisNetwork, redisAddress, "", []byte(env.Secret))
	if err != nil {
		log.Fatal(err.Error())
	}
	return store
}

// Runs the router - router has to be set up beforehand
func (s *serverConfig) runRouter() {
	env := helpers.RetrieveWebAppEnv()
	routerAddress := env.Address()
	err := s.router.Run(routerAddress)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
}
