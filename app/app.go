package app

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

type App struct {
	*gin.Engine
	db *sql.DB
}

// Return a new app instance
func NewApp() *App {
	// Setup db:
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	// language=MySQL
	{
		// Create and use `ElasticJury`
		mustExec(db, "CREATE DATABASE IF NOT EXISTS ElasticJury DEFAULT CHARACTER SET utf8")
		mustExec(db, "USE ElasticJury")
		mustExecScriptFile(db, initTableScriptPath)
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
		router.GET("/search", MakeSearchHandler(db))
		// Other apis...
	}

	return &App{
		Engine: router,
		db:     db,
	}
}
