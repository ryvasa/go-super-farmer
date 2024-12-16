package main

import (
	"log"

	"github.com/ryvasa/go-super-farmer/pkg/wire"
)

func main() {
	app, err := wire.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
