package app

import (
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/logrus"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
)

type MailApp struct {
	Env      *env.Env
	RabbitMQ messages.RabbitMQ
	Handler  MailHandler
}

func NewApp(
	env *env.Env,
	rabbitMQ messages.RabbitMQ,
	handler MailHandler,
) *MailApp {
	defer rabbitMQ.Close()
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
