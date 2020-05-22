package app

import (
	. "ElasticJury/app/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type App struct {
	*gin.Engine
	db database
}

// Return a new app instance
func NewApp(password string) *App {
	// Setup db:
	db, err := newDatabase(password)
	if err != nil {
		panic(err)
	}
	// language=MySQL
	{
		// Create and use `ElasticJury`
		db.mustExec("CREATE DATABASE IF NOT EXISTS ElasticJury DEFAULT CHARACTER SET utf8")
		db.mustExec("USE ElasticJury")
		db.mustExecScriptFile(InitTableScriptPath)
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
		// Retrieve case id by word, tag, law, judge
		router.POST("/search", db.makeSearchHandler())
		// Retrieve case info by case id
		router.GET("/case", db.makeCaseInfoHandler())
	}

	return &App{
		Engine: router,
		db:     db,
	}
}
