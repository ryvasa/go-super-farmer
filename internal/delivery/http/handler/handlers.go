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
}

func NewHandlers(roleHandler RoleHandler, userHandler UserHandler, landHandler LandHandler, authHandler AuthHandler, commodityHandler CommodityHandler, landCommodityHandler LandCommodityHandler, priceHandler PriceHandler, provinceHandler ProvinceHandler, cityHandler CityHandler, regionHandler RegionHandler) *Handlers {
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
	}
}
