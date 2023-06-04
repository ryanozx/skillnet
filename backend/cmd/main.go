package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
	"gorm.io/gorm"
)

type serverConfig struct {
	db     *gorm.DB
	store  redis.Store
	router *gin.Engine
}

type envVars struct {
	webAppAddr   string
	webAppPort   string
	sessionKey   string
	redisHost    string
	redisPort    string
	redisMaxConn int
	redisSecret  string
}

func main() {
	env := retrieveEnvVars()
	serverConfig := initialiseProdServer(env)
	serverConfig.setupRoutes()
	serverConfig.runRouter(env)
	log.Println("Setup complete!")
}

func initialiseProdServer(env *envVars) *serverConfig {
	router := gin.Default()
	db := database.ConnectProdDatabase()
	server := serverConfig{
		router: router,
		db:     db,
	}
	server.setupRedis(env)
	return &server
}

func retrieveEnvVars() *envVars {
	webAppAddr := os.Getenv("WEBAPP_ADDRESS")
	webAppPort := os.Getenv("WEBAPP_PORT")
	sessionKey := os.Getenv("REDIS_SESSION_KEY")
	redisHost := os.Getenv("REDISHOST")
	redisPort := os.Getenv("REDISPORT")
	redisMaxConn, err := strconv.Atoi(os.Getenv("REDIS_MAX_CONNECTIONS"))
	if err != nil {
		panic(err)
	}
	redisSecret := os.Getenv("REDIS_SECRET_KEY")

	envVars := envVars{
		webAppAddr:   webAppAddr,
		webAppPort:   webAppPort,
		sessionKey:   sessionKey,
		redisHost:    redisHost,
		redisPort:    redisPort,
		redisMaxConn: redisMaxConn,
		redisSecret:  redisSecret,
	}
	return &envVars
}

func (server *serverConfig) setupRedis(env *envVars) {
	if env.sessionKey == "" {
		log.Fatalln("Set REDIS_SESSION_KEY to a secret string and try again")
	}
	redisAddress := fmt.Sprintf("%s:%s", env.redisHost, env.redisPort)
	const maxIdleConnections = 10
	const redisNetwork = "tcp"
	store, err := redis.NewStore(maxIdleConnections, redisNetwork, redisAddress, "", []byte(env.redisSecret))
	if err != nil {
		log.Fatal(err.Error())
	}
	server.store = store
}

func (server *serverConfig) runRouter(env *envVars) {
	routerAddress := fmt.Sprintf("%s:%s", env.webAppAddr, env.webAppPort)
	routerErr := server.router.Run(routerAddress)
	if routerErr != nil {
		return
	}
}
