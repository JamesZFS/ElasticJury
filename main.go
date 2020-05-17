package main

import (
	"ElasticJury/app"
	"log"
)

func main() {
	searchEngine := app.NewApp()
	port := app.GetEnvVar("PORT", "8000")
	// Listen and Server in localhost
	log.Fatal(searchEngine.Run("localhost:" + port))
}
