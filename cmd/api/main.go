package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/cmd/api/pkg/wire"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
)

func main() {
	logrus.Log.Info("Starting API service...")
	app, err := wire.InitializeApp()
	if err != nil {
		log.Fatal(err)
		logrus.Log.Fatalf("failed to initialize app: %v", err)
	}
	logrus.Log.Info("API service started successfully")
	defer app.RabbitMQ.Close()
	app.Router.Use(gin.Recovery())
	app.Router.Use(gin.Logger())
	app.Router.Run(":" + app.Env.Server.Port)
}
