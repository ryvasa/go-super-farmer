package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	"gorm.io/gorm"
)

type App struct {
	Router   *gin.Engine
	Env      *env.Env
	DB       *gorm.DB
	RabbitMQ messages.RabbitMQ
}

func NewApp(router *gin.Engine, env *env.Env, db *gorm.DB, rabbitMQ messages.RabbitMQ) *App {
	return &App{
		Router:   router,
		Env:      env,
		DB:       db,
		RabbitMQ: rabbitMQ,
	}
}
func (a *App) Start() error {
	defer a.RabbitMQ.Close()

	a.Router.Use(gin.Recovery())
	a.Router.Use(gin.Logger())

	return a.Router.Run(a.Env.Server.Port)
}
