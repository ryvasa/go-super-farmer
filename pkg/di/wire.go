//go:build wireinject
// +build wireinject

package di

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/ryvasa/go-super-farmer/config"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/user"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/route"
	repository "github.com/ryvasa/go-super-farmer/internal/repository/user"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
)

func InitializeRouter() (*gin.Engine, error) {
	wire.Build(
		config.ConnectDB,
		repository.NewUserRepository,
		usecase.NewUserUsecase,
		handler.NewUserHandler,
		route.NewRouter,
		wire.Bind(new(usecase.UserRepository), new(*repository.UserRepository)),
	)
	return nil, nil
}
