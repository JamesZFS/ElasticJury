package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// Server for hosting backend apis
func setupBackend() *gin.Engine {
	println("\n========== Setting up backend server ==========")
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()
	// Ping test
	router.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
	})
	// Other apis...

	return router
}

func getPort(dft string) string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return port
	} else {
		return dft
	}
}
