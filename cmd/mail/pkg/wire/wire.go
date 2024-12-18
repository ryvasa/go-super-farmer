//go:build wireinject
// +build wireinject

package wire_excel

import (
	"github.com/google/wire"
	"github.com/ryvasa/go-super-farmer/cmd/mail/app"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
)

var allSet = wire.NewSet(
	env.LoadEnv,
	messages.NewRabbitMQ,
	app.NewMailService,
	app.NewMailHandler,
	app.NewApp,
)

func InitializeExcelApp() (*app.MailApp, error) {
	wire.Build(allSet)
	return nil, nil
}
