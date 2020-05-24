package app

import (
	. "ElasticJury/app/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"net/http"
)

type App struct {
	*gin.Engine
	db database
}

// Return a new app instance.
// This method will init the database and tables if the database is not found.
// Panics if an unknown db error occurs.
func NewApp(databaseName, password string) *App {
	// Setup db:
	db, err := newDatabase(databaseName, password)
	if err != nil {
		if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1049 { // unknown database
			fmt.Printf("Creating and reconnecting to %s", databaseName)
			// Create and use `ElasticJury`
			dbRoot, err := newDatabase("", password) // as root
			if err != nil {
				panic(err)
			}
			dbRoot.mustExec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8", databaseName))
			db, err = newDatabase(databaseName, password)
			if err != nil {
				panic(err)
			}
			db.mustExecScriptFile(InitTableScriptPath)
		} else {
			panic(err) // unknown err
		}
	}

	println("Database initialized.")

	// Setup router:
	// Disable Console Color
	// gin.DisableConsoleColor()
	// Release mode is faster
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	{
		// Ping test
		router.GET("/ping", func(context *gin.Context) {
			context.String(http.StatusOK, "pong")
		})
		// Retrieve case id by word, tag, law, judge
		router.POST("/search", db.makeSearchHandler())
		// Retrieve case info by case id
		router.GET("/info", db.makeCaseInfoHandler())
		// Retrieve case detail by one case id
		router.GET("/detail/:id", db.makeCaseDetailHandler())
	}

	return &App{
		Engine: router,
		db:     db,
	}
}
