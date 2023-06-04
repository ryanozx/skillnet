package main

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(t *testing.T) {
	t.Parallel()
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
