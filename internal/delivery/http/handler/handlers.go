package handler

type Handlers struct {
	RoleHandler RoleHandler
	UserHandler UserHandler
}

func NewHandlers(roleHandler RoleHandler, userHandler UserHandler) *Handlers {
	return &Handlers{
		RoleHandler: roleHandler,
		UserHandler: userHandler,
	}
}
