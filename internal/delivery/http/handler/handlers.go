package handler

type Handlers struct {
	RoleHandler      RoleHandler
	UserHandler      UserHandler
	LandHandler      LandHandler
	AuthHandler      AuthHandler
	CommodityHandler CommodityHandler
}

func NewHandlers(roleHandler RoleHandler, userHandler UserHandler, landHandler LandHandler, authHandler AuthHandler, commodityHandler CommodityHandler) *Handlers {
	return &Handlers{
		RoleHandler:      roleHandler,
		UserHandler:      userHandler,
		LandHandler:      landHandler,
		AuthHandler:      authHandler,
		CommodityHandler: commodityHandler,
	}
}
