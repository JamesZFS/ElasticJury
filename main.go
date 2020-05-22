package main

import (
	"ElasticJury/app"
	"log"
)

func main() {
	// Parse environment variables:
	password := app.GetEnvVar("PASSWORD", "")
	port := app.GetEnvVar("PORT", "8000")

	searchEngine := app.NewApp(password)
	// Listen and Server in localhost
	log.Fatal(searchEngine.Run("localhost:" + port))
}
