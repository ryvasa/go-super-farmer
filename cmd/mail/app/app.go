package app

import (
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	mail_handler "github.com/ryvasa/go-super-farmer/service_mail/handler"
)

type MailApp struct {
	Env      *env.Env
	RabbitMQ messages.RabbitMQ
	Handler  mail_handler.MailHandler
}

func NewApp(
	env *env.Env,
	rabbitMQ messages.RabbitMQ,
	handler mail_handler.MailHandler,
) *MailApp {
	err := handler.ConsumerHandler()
	if err != nil {
		logrus.Log.Fatalf("failed to initiate consumer handler: %v", err)
	}
	return &MailApp{
		Env:      env,
		RabbitMQ: rabbitMQ,
		// Handler:  handler,
	}
}
