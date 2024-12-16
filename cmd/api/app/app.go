package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/internal/delivery/rabbitmq"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"gorm.io/gorm"
)

type App struct {
	Router    *gin.Engine
	Env       *env.Env
	DB        *gorm.DB
	Publisher rabbitmq.Publisher
}

func NewApp(router *gin.Engine, env *env.Env, db *gorm.DB, publisher rabbitmq.Publisher) *App {
	return &App{
		Router:    router,
		Env:       env,
		DB:        db,
		Publisher: publisher,
	}
}
func (a *App) Start() error {
	return a.Router.Run(a.Env.Server.Port)
}
