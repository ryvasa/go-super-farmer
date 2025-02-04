package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	pb "github.com/ryvasa/go-super-farmer/proto/generated"
	"gorm.io/gorm"
)

type App struct {
	Router       *gin.Engine
	Env          *env.Env
	DB           *gorm.DB
	RabbitMQ     messages.RabbitMQ
	ReportClient pb.ReportServiceClient
}

func NewApp(
	router *gin.Engine,
	env *env.Env,
	db *gorm.DB,
	rabbitMQ messages.RabbitMQ,
	reportClient pb.ReportServiceClient,
) *App {
	return &App{
		Router:       router,
		Env:          env,
		DB:           db,
		RabbitMQ:     rabbitMQ,
		ReportClient: reportClient,
	}
}
