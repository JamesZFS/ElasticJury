package app

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"strings"
)

type App struct {
	*gin.Engine
	db *sql.DB
}

/// Return a new app instance
func NewApp() *App {
	// Setup db:
	db, err1 := sql.Open("mysql", dataBaseName)
	{
		if err1 != nil {
			panic(err1)
		}
		bytes, err2 := ioutil.ReadFile(initScriptPath)
		if err2 != nil {
			panic(err2)
		}
		scripts := strings.Split(string(bytes), ";")
		for _, script := range scripts {
			if strings.TrimSpace(script) == "" {
				continue
			}
			if _, err3 := db.Exec(script); err3 != nil { // execute init sql script
				panic(err3)
			}
		}
		println("Database initialized.")
	}

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
