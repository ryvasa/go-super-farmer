package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
	Env    *env.Env
	DB     *gorm.DB
}

func NewApp(router *gin.Engine, env *env.Env, db *gorm.DB) *App {
	return &App{
		Router: router,
		Env:    env,
		DB:     db,
	}
}
func (a *App) Start() error {
	return a.Router.Run(a.Env.Server.Port)
}
