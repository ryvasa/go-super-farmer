package handler

import handler_interface "github.com/ryvasa/go-super-farmer/service_api/delivery/http/handler/interface"

type Handlers struct {
	RoleHandler          handler_interface.RoleHandler
	UserHandler          handler_interface.UserHandler
	LandHandler          handler_interface.LandHandler
	AuthHandler          handler_interface.AuthHandler
	CommodityHandler     handler_interface.CommodityHandler
	LandCommodityHandler handler_interface.LandCommodityHandler
	PriceHandler         handler_interface.PriceHandler
	ProvinceHandler      handler_interface.ProvinceHandler
	CityHandler          handler_interface.CityHandler
	DemandHandler        handler_interface.DemandHandler
	SupplyHandler        handler_interface.SupplyHandler
	HarvestHandler       handler_interface.HarvestHandler
	SaleHandler          handler_interface.SaleHandler
}

func NewHandlers(
	roleHandler handler_interface.RoleHandler,
	userHandler handler_interface.UserHandler,
	landHandler handler_interface.LandHandler,
	authHandler handler_interface.AuthHandler, commodityHandler handler_interface.CommodityHandler, landCommodityHandler handler_interface.LandCommodityHandler,
	priceHandler handler_interface.PriceHandler,
	provinceHandler handler_interface.ProvinceHandler,
	cityHandler handler_interface.CityHandler,
	demandHandler handler_interface.DemandHandler,
	supplyHandler handler_interface.SupplyHandler,
	harvestHandler handler_interface.HarvestHandler,
	saleHandler handler_interface.SaleHandler,
) *Handlers {
	return &Handlers{
		RoleHandler:          roleHandler,
		UserHandler:          userHandler,
		LandHandler:          landHandler,
		AuthHandler:          authHandler,
		CommodityHandler:     commodityHandler,
		LandCommodityHandler: landCommodityHandler,
		PriceHandler:         priceHandler,
		ProvinceHandler:      provinceHandler,
		CityHandler:          cityHandler,
		DemandHandler:        demandHandler,
		SupplyHandler:        supplyHandler,
		HarvestHandler:       harvestHandler,
		SaleHandler:          saleHandler,
	}
}
