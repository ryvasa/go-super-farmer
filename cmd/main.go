package main

import (
	"log"

	"github.com/ryvasa/go-super-farmer/pkg/di"
)

func main() {
	router, err := di.InitializeRouter()
	if err != nil {
		log.Fatalf("failed to initialize router: %v", err)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
