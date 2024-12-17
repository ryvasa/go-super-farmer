//go:build wireinject
// +build wireinject

package wire_excel

import (
	"github.com/google/wire"
	"github.com/ryvasa/go-super-farmer/cmd/excel/app"
	"github.com/ryvasa/go-super-farmer/cmd/excel/internal/handler"
	"github.com/ryvasa/go-super-farmer/cmd/excel/internal/repository"
	"github.com/ryvasa/go-super-farmer/cmd/excel/internal/usecase"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
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
