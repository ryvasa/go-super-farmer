package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	report_handler "github.com/ryvasa/go-super-farmer/service_report/dilevery/http/handler"
	"gorm.io/gorm"
)

type ReportApp struct {
	Router   *gin.Engine
	Env      *env.Env
	DB       *gorm.DB
	RabbitMQ messages.RabbitMQ
	Handler  *report_handler.Handlers
}

func NewApp(
	router *gin.Engine,
	env *env.Env,
	db *gorm.DB,
	rabbitMQ messages.RabbitMQ,
	handler *report_handler.Handlers,
) *ReportApp {
	return &ReportApp{
		Router:   router,
		Env:      env,
		DB:       db,
		RabbitMQ: rabbitMQ,
		Handler:  handler,
	}
}
