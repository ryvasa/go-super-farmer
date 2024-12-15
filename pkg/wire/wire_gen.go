// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/ryvasa/go-super-farmer/cmd/app"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/implementation"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/route"
	"github.com/ryvasa/go-super-farmer/internal/repository/implementation"
	"github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	"github.com/ryvasa/go-super-farmer/pkg/auth/token"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/ryvasa/go-super-farmer/pkg/env"
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
	userUsecase := usecase_implementation.NewUserUsecase(userRepository, hasher)
	userHandler := handler_implementation.NewUserHandler(userUsecase)
	landRepository := repository_implementation.NewLandRepository(db)
	landUsecase := usecase_implementation.NewLandUsecase(landRepository, userRepository)
	authUtil := utils.NewAuthUtil()
	landHandler := handler_implementation.NewLandHandler(landUsecase, authUtil)
	tokenToken := token.NewToken(envEnv)
	authUsecase := usecase_implementation.NewAuthUsecase(userRepository, tokenToken, hasher)
	authHandler := handler_implementation.NewAuthHandler(authUsecase)
	commodityRepository := repository_implementation.NewCommodityRepository(db)
	commodityUsecase := usecase_implementation.NewCommodityUsecase(commodityRepository)
	commodityHandler := handler_implementation.NewCommodityHandler(commodityUsecase)
	landCommodityRepository := repository_implementation.NewLandCommodityRepository(db)
	landCommodityUsecase := usecase_implementation.NewLandCommodityUsecase(landCommodityRepository, landRepository, commodityRepository)
	landCommodityHandler := handler_implementation.NewLandCommodityHandler(landCommodityUsecase)
	priceRepository := repository_implementation.NewPriceRepository(db)
	priceHistoryRepository := repository_implementation.NewPriceHistoryRepository(db)
	regionRepository := repository_implementation.NewRegionRepository(db)
	priceUsecase := usecase_implementation.NewPriceUsecase(priceRepository, priceHistoryRepository, regionRepository, commodityRepository)
	priceHandler := handler_implementation.NewPriceHandler(priceUsecase)
	provinceRepository := repository_implementation.NewProvinceRepository(db)
	provinceUsecase := usecase_implementation.NewProvinceUsecase(provinceRepository)
	provinceHandler := handler_implementation.NewProvinceHandler(provinceUsecase)
	cityRepository := repository_implementation.NewCityRepository(db)
	cityUsecase := usecase_implementation.NewCityUsecase(cityRepository)
	cityHandler := handler_implementation.NewCityHandler(cityUsecase)
	regionUsecase := usecase_implementation.NewRegionUsecase(regionRepository, cityRepository, provinceRepository)
	regionHandler := handler_implementation.NewRegionHandler(regionUsecase)
	demandRepository := repository_implementation.NewDemandRepository(db)
	demandHistoryRepository := repository_implementation.NewDemandHistoryRepository(db)
	demandUsecase := usecase_implementation.NewDemandUsecase(demandRepository, demandHistoryRepository, commodityRepository, regionRepository)
	demandHandler := handler_implementation.NewDemandHandler(demandUsecase)
	supplyRepository := repository_implementation.NewSupplyRepository(db)
	supplyHistoryRepository := repository_implementation.NewSupplyHistoryRepository(db)
	supplyUsecase := usecase_implementation.NewSupplyUsecase(supplyRepository, supplyHistoryRepository, commodityRepository, regionRepository)
	supplyHandler := handler_implementation.NewSupplyHandler(supplyUsecase)
	harvestRepository := repository_implementation.NewHarvestRepository(db)
	harvestUsecase := usecase_implementation.NewHarvestUsecase(harvestRepository, regionRepository, landCommodityRepository)
	harvestHandler := handler_implementation.NewHarvestHandler(harvestUsecase)
	handlers := handler.NewHandlers(roleHandler, userHandler, landHandler, authHandler, commodityHandler, landCommodityHandler, priceHandler, provinceHandler, cityHandler, regionHandler, demandHandler, supplyHandler, harvestHandler)
	engine := route.NewRouter(handlers)
	appApp := app.NewApp(engine, envEnv, db)
	return appApp, nil
}

// wire.go:

var tokenSet = wire.NewSet(token.NewToken)

var utilSet = wire.NewSet(utils.NewAuthUtil, utils.NewHasher)

var repositorySet = wire.NewSet(repository_implementation.NewRoleRepository, repository_implementation.NewUserRepository, repository_implementation.NewLandRepository, repository_implementation.NewCommodityRepository, repository_implementation.NewLandCommodityRepository, repository_implementation.NewPriceRepository, repository_implementation.NewProvinceRepository, repository_implementation.NewCityRepository, repository_implementation.NewRegionRepository, repository_implementation.NewPriceHistoryRepository, repository_implementation.NewDemandRepository, repository_implementation.NewSupplyRepository, repository_implementation.NewDemandHistoryRepository, repository_implementation.NewSupplyHistoryRepository, repository_implementation.NewHarvestRepository)

var usecaseSet = wire.NewSet(usecase_implementation.NewRoleUsecase, usecase_implementation.NewUserUsecase, usecase_implementation.NewLandUsecase, usecase_implementation.NewAuthUsecase, usecase_implementation.NewCommodityUsecase, usecase_implementation.NewLandCommodityUsecase, usecase_implementation.NewPriceUsecase, usecase_implementation.NewProvinceUsecase, usecase_implementation.NewCityUsecase, usecase_implementation.NewRegionUsecase, usecase_implementation.NewDemandUsecase, usecase_implementation.NewSupplyUsecase, usecase_implementation.NewHarvestUsecase)

var handlerSet = wire.NewSet(handler_implementation.NewRoleHandler, handler_implementation.NewUserHandler, handler_implementation.NewLandHandler, handler_implementation.NewAuthHandler, handler_implementation.NewCommodityHandler, handler_implementation.NewLandCommodityHandler, handler_implementation.NewPriceHandler, handler_implementation.NewProvinceHandler, handler_implementation.NewCityHandler, handler_implementation.NewRegionHandler, handler_implementation.NewDemandHandler, handler_implementation.NewSupplyHandler, handler_implementation.NewHarvestHandler)