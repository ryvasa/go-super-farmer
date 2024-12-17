package app

import (
	"log"

	"github.com/ryvasa/go-super-farmer/cmd/excel/internal/handler"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	"gorm.io/gorm"
)

type ExcelApp struct {
	Env      *env.Env
	DB       *gorm.DB
	RabbitMQ messages.RabbitMQ
	Handler  handler.RabbitMQHandler
}

func NewApp(
	env *env.Env,
	db *gorm.DB,
	rabbitMQ messages.RabbitMQ,
	handler handler.RabbitMQHandler,
) *ExcelApp {
	defer rabbitMQ.Close()
	err := handler.ConsumerHandler()
	if err != nil {
		log.Fatal(err)
	}
	return &ExcelApp{
		Env:      env,
		DB:       db,
		RabbitMQ: rabbitMQ,
		Handler:  handler,
	}
}
