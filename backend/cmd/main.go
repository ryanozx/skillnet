package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/ryanozx/skillnet/database"
)

func main() {
	router := gin.Default()
	database.ConnectDatabase()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(config))

	sessionKey := os.Getenv("REDIS_SESSION_KEY")
	RedisHost := os.Getenv("REDISHOST")
	RedisPort := os.Getenv("REDISPORT")
	if sessionKey == "" {
		log.Fatal("ERROR: Set REDIS_SESSION_KEY to a secret string and try again")
	}
	store, redisErr := redis.NewStore(10, "tcp", RedisHost+":"+RedisPort, "", []byte("secret"))
	if redisErr != nil {
		log.Fatal(redisErr.Error())
	}

	registerRoutes(router, store)

	routerErr := router.Run("0.0.0.0:8080")
	if routerErr != nil {
		return
	}
	log.Output(1, "Setup complete!")
}
