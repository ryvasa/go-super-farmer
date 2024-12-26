//go:build wireinject
// +build wireinject

package wire_excel

import (
	"github.com/google/wire"
	"github.com/ryvasa/go-super-farmer/cmd/report/app"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	report_handler "github.com/ryvasa/go-super-farmer/service_report/dilevery/http/handler"
	report_route "github.com/ryvasa/go-super-farmer/service_report/dilevery/http/routes"
	"github.com/ryvasa/go-super-farmer/service_report/repository"
	"github.com/ryvasa/go-super-farmer/service_report/usecase"
	"github.com/ryvasa/go-super-farmer/utils"
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
	usecase.NewReportUsecase,

	// Handler
	report_handler.NewReportHandler,

	// App
	app.NewApp,

	report_route.NewRoutes,
	report_handler.NewHandlers,

	utils.NewGlobFunc,
)

func InitializeReportApp() (*app.ReportApp, error) {
	wire.Build(allSet)
	return nil, nil
}
