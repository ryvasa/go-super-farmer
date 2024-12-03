package main

import (
	"log"
	"os"

	"github.com/ryvasa/go-super-farmer/pkg/di"
)

func main() {
	router, err := di.InitializeRouter()
	if err != nil {
		log.Fatalf("failed to initialize router: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	if err := router.Run(port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
