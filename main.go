package main

import (
	"ElasticJury/app"
	"log"
	"os"
)

func getPort(dft string) string {
	if port, ok := os.LookupEnv("PORT"); ok {
		return port
	} else {
		return dft
	}
}

func main() {
	searchEngine := app.NewApp()
	port := getPort("3000")
	// Listen and Server in localhost
	log.Fatal(searchEngine.Run("localhost:" + port))
}
