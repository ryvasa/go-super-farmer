package handler

type Handlers struct {
	RoleHandler          RoleHandler
	UserHandler          UserHandler
	LandHandler          LandHandler
	AuthHandler          AuthHandler
	CommodityHandler     CommodityHandler
	LandCommodityHandler LandCommodityHandler
	PriceHandler         PriceHandler
}

func NewHandlers(roleHandler RoleHandler, userHandler UserHandler, landHandler LandHandler, authHandler AuthHandler, commodityHandler CommodityHandler, landCommodityHandler LandCommodityHandler, priceHandler PriceHandler) *Handlers {
	return &Handlers{
		RoleHandler:          roleHandler,
		UserHandler:          userHandler,
		LandHandler:          landHandler,
		AuthHandler:          authHandler,
		CommodityHandler:     commodityHandler,
		LandCommodityHandler: landCommodityHandler,
		PriceHandler:         priceHandler,
	}
}
