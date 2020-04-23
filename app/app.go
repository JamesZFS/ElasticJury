package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type App struct {
	*gin.Engine
	db database
}

// Return a new app instance
func NewApp() *App {
	// Setup db:
	db, err := newDatabase()
	if err != nil {
		panic(err)
	}
	// language=MySQL
	{
		// Create and use `ElasticJury`
		db.mustExec("CREATE DATABASE IF NOT EXISTS ElasticJury DEFAULT CHARACTER SET utf8")
		db.mustExec("USE ElasticJury")
		db.mustExecScriptFile(initTableScriptPath)
	}
	println("Database initialized.")

	// Setup router:
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()
	{
		// Ping test
		router.GET("/ping", func(context *gin.Context) {
			context.String(http.StatusOK, "pong")
		})
		// Retrieve api:
		router.GET("/search", db.makeSearchHandler())
		// Other apis...
	}

	return &App{
		Engine: router,
		db:     db,
	}
}
