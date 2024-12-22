//go:build wireinject
// +build wireinject

package wire_excel

import (
	"github.com/google/wire"
	"github.com/ryvasa/go-super-farmer/cmd/excel/app"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	"github.com/ryvasa/go-super-farmer/service_excel/handler"
	"github.com/ryvasa/go-super-farmer/service_excel/repository"
	"github.com/ryvasa/go-super-farmer/service_excel/usecase"
)

var allSet = wire.NewSet(
	// Infrastructure
	env.LoadEnv,
	database.NewPostgres,
	messages.NewRabbitMQ,

	// Repository
	repository.NewReportRepositoryImpl,

	// Service
	usecase.NewExcelImpl,
	usecase.NewRabbitMQUsecase,

	// Handler
	handler.NewRabbitMQHandler,

	// App
	app.NewApp,
)

func InitializeExcelApp() (*app.ExcelApp, error) {
	wire.Build(allSet)
	return nil, nil
}
