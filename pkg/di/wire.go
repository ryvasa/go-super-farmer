//go:build wireinject
// +build wireinject

package di

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/route"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/pkg/auth/token"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/utils"
)

var roleSet = wire.NewSet(
	repository.NewRoleRepository,
	usecase.NewRoleUsecase,
	handler.NewRoleHandler,
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	usecase.NewUserUsecase,
	handler.NewUserHandler,
)

var landSet = wire.NewSet(
	repository.NewLandRepository,
	usecase.NewLandUsecase,
	handler.NewLandHandler,
)

var authSet = wire.NewSet(
	usecase.NewAuthUsecase,
	handler.NewAuthHandler,
)

var tokenSet = wire.NewSet(
	token.NewToken,
)

var authUtilSet = wire.NewSet(
	utils.NewAuthUtil,
)

var hashSet = wire.NewSet(
	utils.NewHasher,
)

var commoditySet = wire.NewSet(
	repository.NewCommodityRepository,
	usecase.NewCommodityUsecase,
	handler.NewCommodityHandler,
)

var landCommoditySet = wire.NewSet(
	repository.NewLandCommodityRepository,
	usecase.NewLandCommodityUsecase,
	handler.NewLandCommodityHandler,
)

var priceSet = wire.NewSet(
	repository.NewPriceRepository,
	usecase.NewPriceUsecase,
	handler.NewPriceHandler,
)

var provinceSet = wire.NewSet(
	repository.NewProvinceRepository,
	usecase.NewProvinceUsecase,
	handler.NewProvinceHandler,
)

var citySet = wire.NewSet(
	repository.NewCityRepository,
	usecase.NewCityUsecase,
	handler.NewCityHandler,
)

var regionSet = wire.NewSet(
	repository.NewRegionRepository,
	usecase.NewRegionUsecase,
	handler.NewRegionHandler,
)

var priceHistorySet = wire.NewSet(
	repository.NewPriceHistoryRepository,
)

var demandSet = wire.NewSet(
	repository.NewDemandRepository,
	usecase.NewDemandUsecase,
	handler.NewDemandHandler,
)

var supplySet = wire.NewSet(
	repository.NewSupplyRepository,
	usecase.NewSupplyUsecase,
	handler.NewSupplyHandler,
)

var demandHistorySet = wire.NewSet(
	repository.NewDemandHistoryRepository,
)

var supplyHistorySet = wire.NewSet(
	repository.NewSupplyHistoryRepository,
)

func InitializeRouter() (*gin.Engine, error) {
	wire.Build(
		env.LoadEnv,
		database.ConnectDB,
		database.ProvideDSN,
		handler.NewHandlers,
		route.NewRouter,
		roleSet,
		userSet,
		landSet,
		authSet,
		tokenSet,
		authUtilSet,
		commoditySet,
		landCommoditySet,
		priceSet,
		provinceSet,
		citySet,
		regionSet,
		priceHistorySet,
		hashSet,
		demandSet,
		supplySet,
		demandHistorySet,
		supplyHistorySet,
	)
	return nil, nil
}
