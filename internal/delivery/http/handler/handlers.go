package handler

type Handlers struct {
	RoleHandler RoleHandler
	UserHandler UserHandler
	LandHandler LandHandler
	AuthHandler AuthHandler
}

func NewHandlers(roleHandler RoleHandler, userHandler UserHandler, landHandler LandHandler, authHandler AuthHandler) *Handlers {
	return &Handlers{
		RoleHandler: roleHandler,
		UserHandler: userHandler,
		LandHandler: landHandler,
		AuthHandler: authHandler,
	}
}
