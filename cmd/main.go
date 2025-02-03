package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/pkg/wire"
)

func main() {
	logrus.Log.Info("Starting API service...")

	// Initialize app
	app, err := wire.InitializeApp()
	if err != nil {
		logrus.Log.Fatalf("Failed to initialize app: %v", err)
	}

	// Setup Gin
	app.Router.Use(gin.Recovery())
	app.Router.Use(gin.Logger())

	// Start server in goroutine
	go func() {
		if err := app.Router.Run(":" + app.Env.Server.Port); err != nil {
			logrus.Log.Fatalf("Failed to start server: %v", err)
		}
	}()

	logrus.Log.Infof("API service started on port %s", app.Env.Server.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logrus.Log.Info("Shutting down server...")

	// Cleanup
	if app.RabbitMQ != nil {
		app.RabbitMQ.Close()
	}

	logrus.Log.Info("Server exited properly")
}
