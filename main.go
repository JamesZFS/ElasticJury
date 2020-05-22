package main

import (
	"ElasticJury/app"
	"ElasticJury/app/common"
	"ElasticJury/app/natural"
	"log"
)

func main() {
	natural.Initialize()
	defer natural.Finalize()

	// Parse environment variables:
	password := common.GetEnvVar("PASSWORD", "")
	port := common.GetEnvVar("PORT", "8000")

	searchEngine := app.NewApp(password)
	// Listen and Server in localhost
	log.Fatal(searchEngine.Run("localhost:" + port))
}
