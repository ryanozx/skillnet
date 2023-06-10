package main

import (
	"log"

	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"github.com/ryanozx/skillnet/helpers"
	"gorm.io/gorm"
)

type serverConfig struct {
	db     *gorm.DB
	store  redis.Store
	router *gin.Engine
}

func main() {
	serverConfig := initialiseProdServer()
	serverConfig.setupRoutes()
	serverConfig.runRouter()
	log.Println("Setup complete!")
}

func initialiseProdServer() *serverConfig {
	router := gin.Default()
	db := database.ConnectProdDatabase()
	server := serverConfig{
		router: router,
		db:     db,
	}
	server.setupRedis()
	return &server
}

func (server *serverConfig) setupRedis() {
	env := helpers.RetrieveRedisEnv()
	redisAddress := env.Address()
	const redisNetwork = "tcp"
	store, err := redis.NewStore(env.MaxConn, redisNetwork, redisAddress, "", []byte(env.Secret))
	if err != nil {
		log.Fatal(err.Error())
	}
	server.store = store
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
