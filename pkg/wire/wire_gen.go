// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/ryvasa/go-super-farmer/cmd/api/app"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/implementation"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/route"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/internal/repository/implementation"
	"github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	"github.com/ryvasa/go-super-farmer/pkg/auth/token"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/ryvasa/go-super-farmer/pkg/database/cache"
	"github.com/ryvasa/go-super-farmer/pkg/database/transaction"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
	"github.com/ryvasa/go-super-farmer/utils"
)

// Injectors from wire.go:

func InitializeApp() (*app.App, error) {
	envEnv, err := env.LoadEnv()
	if err != nil {
		return nil, err
	}
	db, err := database.NewPostgres(envEnv)
	if err != nil {
		return nil, err
	}
	roleRepository := repository_implementation.NewRoleRepository(db)
	roleUsecase := usecase_implementation.NewRoleUsecase(roleRepository)
	roleHandler := handler_implementation.NewRoleHandler(roleUsecase)
	userRepository := repository_implementation.NewUserRepository(db)
	hasher := utils.NewHasher()
	client := database.NewRedisClient(envEnv)
	cacheCache := cache.NewRedisCache(client)
	userUsecase := usecase_implementation.NewUserUsecase(userRepository, hasher, cacheCache)
	tokenToken := token.NewToken(envEnv)
	rabbitMQ, err := messages.NewRabbitMQ(envEnv)
	if err != nil {
		return nil, err
	}
	otp := utils.NewOTPGenerator()
	authUsecase := usecase_implementation.NewAuthUsecase(userRepository, tokenToken, hasher, rabbitMQ, cacheCache, otp)
	userHandler := handler_implementation.NewUserHandler(userUsecase, authUsecase)
	landRepository := repository_implementation.NewLandRepository(db)
	landUsecase := usecase_implementation.NewLandUsecase(landRepository, userRepository)
	authUtil := utils.NewAuthUtil()
	landHandler := handler_implementation.NewLandHandler(landUsecase, authUtil)
	authHandler := handler_implementation.NewAuthHandler(authUsecase)
	commodityRepository := repository_implementation.NewCommodityRepository(db)
	commodityUsecase := usecase_implementation.NewCommodityUsecase(commodityRepository, cacheCache)
	commodityHandler := handler_implementation.NewCommodityHandler(commodityUsecase)
	landCommodityRepository := repository_implementation.NewLandCommodityRepository(db)
	landCommodityUsecase := usecase_implementation.NewLandCommodityUsecase(landCommodityRepository, landRepository, commodityRepository, cacheCache)
	landCommodityHandler := handler_implementation.NewLandCommodityHandler(landCommodityUsecase)
	baseRepository := repository.NewBaseRepository(db)
	priceRepository := repository_implementation.NewPriceRepository(baseRepository)
	priceHistoryRepository := repository_implementation.NewPriceHistoryRepository(baseRepository)
	cityRepository := repository_implementation.NewCityRepository(db)
	transactionManager := transaction.NewTransactionManager(db)
	priceUsecase := usecase_implementation.NewPriceUsecase(priceRepository, priceHistoryRepository, cityRepository, commodityRepository, rabbitMQ, transactionManager, cacheCache)
	priceHandler := handler_implementation.NewPriceHandler(priceUsecase)
	provinceRepository := repository_implementation.NewProvinceRepository(db)
	provinceUsecase := usecase_implementation.NewProvinceUsecase(provinceRepository)
	provinceHandler := handler_implementation.NewProvinceHandler(provinceUsecase)
	cityUsecase := usecase_implementation.NewCityUsecase(cityRepository)
	cityHandler := handler_implementation.NewCityHandler(cityUsecase)
	demandRepository := repository_implementation.NewDemandRepository(baseRepository)
	demandHistoryRepository := repository_implementation.NewDemandHistoryRepository(baseRepository)
	demandUsecase := usecase_implementation.NewDemandUsecase(demandRepository, demandHistoryRepository, commodityRepository, cityRepository, transactionManager)
	demandHandler := handler_implementation.NewDemandHandler(demandUsecase)
	supplyRepository := repository_implementation.NewSupplyRepository(baseRepository)
	supplyHistoryRepository := repository_implementation.NewSupplyHistoryRepository(baseRepository)
	supplyUsecase := usecase_implementation.NewSupplyUsecase(supplyRepository, supplyHistoryRepository, commodityRepository, cityRepository, transactionManager)
	supplyHandler := handler_implementation.NewSupplyHandler(supplyUsecase)
	harvestRepository := repository_implementation.NewHarvestRepository(db)
	globFunc := utils.NewGlobFunc()
	harvestUsecase := usecase_implementation.NewHarvestUsecase(harvestRepository, cityRepository, landCommodityRepository, rabbitMQ, cacheCache, globFunc)
	harvestHandler := handler_implementation.NewHarvestHandler(harvestUsecase)
	handlers := handler.NewHandlers(roleHandler, userHandler, landHandler, authHandler, commodityHandler, landCommodityHandler, priceHandler, provinceHandler, cityHandler, demandHandler, supplyHandler, harvestHandler)
	engine := route.NewRouter(handlers)
	appApp := app.NewApp(engine, envEnv, db, rabbitMQ)
	return appApp, nil
}

// wire.go:

var tokenSet = wire.NewSet(token.NewToken)

var utilSet = wire.NewSet(utils.NewAuthUtil, utils.NewHasher, utils.NewOTPGenerator, utils.NewGlobFunc)

var repositorySet = wire.NewSet(repository.NewBaseRepository, repository_implementation.NewRoleRepository, repository_implementation.NewUserRepository, repository_implementation.NewLandRepository, repository_implementation.NewCommodityRepository, repository_implementation.NewLandCommodityRepository, repository_implementation.NewPriceRepository, repository_implementation.NewProvinceRepository, repository_implementation.NewCityRepository, repository_implementation.NewPriceHistoryRepository, repository_implementation.NewDemandRepository, repository_implementation.NewSupplyRepository, repository_implementation.NewDemandHistoryRepository, repository_implementation.NewSupplyHistoryRepository, repository_implementation.NewHarvestRepository)

var usecaseSet = wire.NewSet(usecase_implementation.NewRoleUsecase, usecase_implementation.NewUserUsecase, usecase_implementation.NewLandUsecase, usecase_implementation.NewAuthUsecase, usecase_implementation.NewCommodityUsecase, usecase_implementation.NewLandCommodityUsecase, usecase_implementation.NewPriceUsecase, usecase_implementation.NewProvinceUsecase, usecase_implementation.NewCityUsecase, usecase_implementation.NewDemandUsecase, usecase_implementation.NewSupplyUsecase, usecase_implementation.NewHarvestUsecase)

var handlerSet = wire.NewSet(handler_implementation.NewRoleHandler, handler_implementation.NewUserHandler, handler_implementation.NewLandHandler, handler_implementation.NewAuthHandler, handler_implementation.NewCommodityHandler, handler_implementation.NewLandCommodityHandler, handler_implementation.NewPriceHandler, handler_implementation.NewProvinceHandler, handler_implementation.NewCityHandler, handler_implementation.NewDemandHandler, handler_implementation.NewSupplyHandler, handler_implementation.NewHarvestHandler)

var rabbitMQSet = wire.NewSet(messages.NewRabbitMQ)

var cacheSet = wire.NewSet(cache.NewRedisCache)

var databaseSet = wire.NewSet(database.NewPostgres, database.NewRedisClient)

var txManagerSet = wire.NewSet(transaction.NewTransactionManager)
