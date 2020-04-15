package main

import (
	"log"
)

func main() {
	engine := setupBackend()
	port := getPort("3000")
	// Listen and Server in localhost
	log.Fatal(engine.Run("localhost:" + port))
}
