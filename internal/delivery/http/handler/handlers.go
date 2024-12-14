package handler

type Handlers struct {
	RoleHandler          RoleHandler
	UserHandler          UserHandler
	LandHandler          LandHandler
	AuthHandler          AuthHandler
	CommodityHandler     CommodityHandler
	LandCommodityHandler LandCommodityHandler
	PriceHandler         PriceHandler
	ProvinceHandler      ProvinceHandler
	CityHandler          CityHandler
	RegionHandler        RegionHandler
	DemandHandler        DemandHandler
	SupplyHandler        SupplyHandler
	HarvestHandler       HarvestHandler
}

func NewHandlers(roleHandler RoleHandler, userHandler UserHandler, landHandler LandHandler, authHandler AuthHandler, commodityHandler CommodityHandler, landCommodityHandler LandCommodityHandler, priceHandler PriceHandler, provinceHandler ProvinceHandler, cityHandler CityHandler, regionHandler RegionHandler, demandHandler DemandHandler, supplyHandler SupplyHandler, harvestHandler HarvestHandler) *Handlers {
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
		RegionHandler:        regionHandler,
		DemandHandler:        demandHandler,
		SupplyHandler:        supplyHandler,
		HarvestHandler:       harvestHandler,
	}
}
