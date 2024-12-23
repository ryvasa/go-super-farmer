//go:build wireinject
// +build wireinject

package wire_excel

import (
	"github.com/google/wire"
	"github.com/ryvasa/go-super-farmer/cmd/mail/app"
	mail_handler "github.com/ryvasa/go-super-farmer/service_mail/handler"
	mail_usecase "github.com/ryvasa/go-super-farmer/service_mail/usecase"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
)

var allSet = wire.NewSet(
	env.LoadEnv,
	messages.NewRabbitMQ,
	mail_usecase.NewMailUsecase,
	mail_handler.NewMailHandler,
	app.NewApp,
)

func InitializeExcelApp() (*app.MailApp, error) {
	wire.Build(allSet)
	return nil, nil
}
