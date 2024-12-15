//go:build wireinject
// +build wireinject

package di

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	handler_implementation "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/implementation"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/route"
	repository_implementation "github.com/ryvasa/go-super-farmer/internal/repository/implementation"
	usecase_implementation "github.com/ryvasa/go-super-farmer/internal/usecase/implementation"
	"github.com/ryvasa/go-super-farmer/pkg/auth/token"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/ryvasa/go-super-farmer/pkg/env"
	"github.com/ryvasa/go-super-farmer/utils"
)

var roleSet = wire.NewSet(
	repository_implementation.NewRoleRepository,
	usecase_implementation.NewRoleUsecase,
	handler_implementation.NewRoleHandler,
)

var userSet = wire.NewSet(
	repository_implementation.NewUserRepository,
	usecase_implementation.NewUserUsecase,
	handler_implementation.NewUserHandler,
)

var landSet = wire.NewSet(
	repository_implementation.NewLandRepository,
	usecase_implementation.NewLandUsecase,
	handler_implementation.NewLandHandler,
)

var authSet = wire.NewSet(
	usecase_implementation.NewAuthUsecase,
	handler_implementation.NewAuthHandler,
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
	repository_implementation.NewCommodityRepository,
	usecase_implementation.NewCommodityUsecase,
	handler_implementation.NewCommodityHandler,
)

var landCommoditySet = wire.NewSet(
	repository_implementation.NewLandCommodityRepository,
	usecase_implementation.NewLandCommodityUsecase,
	handler_implementation.NewLandCommodityHandler,
)

var priceSet = wire.NewSet(
	repository_implementation.NewPriceRepository,
	usecase_implementation.NewPriceUsecase,
	handler_implementation.NewPriceHandler,
)

var provinceSet = wire.NewSet(
	repository_implementation.NewProvinceRepository,
	usecase_implementation.NewProvinceUsecase,
	handler_implementation.NewProvinceHandler,
)

var citySet = wire.NewSet(
	repository_implementation.NewCityRepository,
	usecase_implementation.NewCityUsecase,
	handler_implementation.NewCityHandler,
)

var regionSet = wire.NewSet(
	repository_implementation.NewRegionRepository,
	usecase_implementation.NewRegionUsecase,
	handler_implementation.NewRegionHandler,
)

var priceHistorySet = wire.NewSet(
	repository_implementation.NewPriceHistoryRepository,
)

var demandSet = wire.NewSet(
	repository_implementation.NewDemandRepository,
	usecase_implementation.NewDemandUsecase,
	handler_implementation.NewDemandHandler,
)

var supplySet = wire.NewSet(
	repository_implementation.NewSupplyRepository,
	usecase_implementation.NewSupplyUsecase,
	handler_implementation.NewSupplyHandler,
)

var demandHistorySet = wire.NewSet(
	repository_implementation.NewDemandHistoryRepository,
)

var supplyHistorySet = wire.NewSet(
	repository_implementation.NewSupplyHistoryRepository,
)

var harvestSet = wire.NewSet(
	repository_implementation.NewHarvestRepository,
	usecase_implementation.NewHarvestUsecase,
	handler_implementation.NewHarvestHandler,
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
		harvestSet,
	)
	return nil, nil
}
