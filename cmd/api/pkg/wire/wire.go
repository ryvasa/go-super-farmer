//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/ryvasa/go-super-farmer/cmd/api/app"
	"github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler"
	handler_implementation "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/implementation"
	"github.com/ryvasa/go-super-farmer/service_api/delivery/http/route"
	"github.com/ryvasa/go-super-farmer/service_api/repository"
	repository_implementation "github.com/ryvasa/go-super-farmer/service_api/repository/implementation"
	usecase_implementation "github.com/ryvasa/go-super-farmer/service_api/usecase/implementation"
	"github.com/ryvasa/go-super-farmer/pkg/auth/token"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/ryvasa/go-super-farmer/pkg/database/cache"
	"github.com/ryvasa/go-super-farmer/pkg/database/transaction"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	"github.com/ryvasa/go-super-farmer/utils"
)

var tokenSet = wire.NewSet(
	token.NewToken,
)

var utilSet = wire.NewSet(
	utils.NewAuthUtil,
	utils.NewHasher,
	utils.NewOTPGenerator,
	utils.NewGlobFunc,
)

var repositorySet = wire.NewSet(
	repository.NewBaseRepository,
	repository_implementation.NewRoleRepository,
	repository_implementation.NewUserRepository,
	repository_implementation.NewLandRepository,
	repository_implementation.NewCommodityRepository,
	repository_implementation.NewLandCommodityRepository,
	repository_implementation.NewPriceRepository,
	repository_implementation.NewProvinceRepository,
	repository_implementation.NewCityRepository,
	repository_implementation.NewPriceHistoryRepository,
	repository_implementation.NewDemandRepository,
	repository_implementation.NewSupplyRepository,
	repository_implementation.NewDemandHistoryRepository,
	repository_implementation.NewSupplyHistoryRepository,
	repository_implementation.NewHarvestRepository,
)

var usecaseSet = wire.NewSet(
	usecase_implementation.NewRoleUsecase,
	usecase_implementation.NewUserUsecase,
	usecase_implementation.NewLandUsecase,
	usecase_implementation.NewAuthUsecase,
	usecase_implementation.NewCommodityUsecase,
	usecase_implementation.NewLandCommodityUsecase,
	usecase_implementation.NewPriceUsecase,
	usecase_implementation.NewProvinceUsecase,
	usecase_implementation.NewCityUsecase,
	usecase_implementation.NewDemandUsecase,
	usecase_implementation.NewSupplyUsecase,
	usecase_implementation.NewHarvestUsecase,
)

var handlerSet = wire.NewSet(
	handler_implementation.NewRoleHandler,
	handler_implementation.NewUserHandler,
	handler_implementation.NewLandHandler,
	handler_implementation.NewAuthHandler,
	handler_implementation.NewCommodityHandler,
	handler_implementation.NewLandCommodityHandler,
	handler_implementation.NewPriceHandler,
	handler_implementation.NewProvinceHandler,
	handler_implementation.NewCityHandler,
	handler_implementation.NewDemandHandler,
	handler_implementation.NewSupplyHandler,
	handler_implementation.NewHarvestHandler,
)

var rabbitMQSet = wire.NewSet(
	messages.NewRabbitMQ,
)

var cacheSet = wire.NewSet(
	cache.NewRedisCache,
)

var databaseSet = wire.NewSet(
	database.NewPostgres,
	database.NewRedisClient,
)

var txManagerSet = wire.NewSet(
	transaction.NewTransactionManager,
)

func InitializeApp() (*app.App, error) {
	wire.Build(
		env.LoadEnv,
		handler.NewHandlers,
		route.NewRouter,
		app.NewApp,
		databaseSet,
		rabbitMQSet,
		tokenSet,
		utilSet,
		repositorySet,
		usecaseSet,
		handlerSet,
		cacheSet,
		txManagerSet,
	)
	return nil, nil
}