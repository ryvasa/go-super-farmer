//go:build wireinject
// +build wireinject

package di

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler/user"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/route"
	repository "github.com/ryvasa/go-super-farmer/internal/repository/user"
	usecase "github.com/ryvasa/go-super-farmer/internal/usecase/user"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/ryvasa/go-super-farmer/pkg/env"
)

func InitializeRouter() (*gin.Engine, error) {
	wire.Build(
		env.LoadEnv,
		database.ConnectDB,
		database.ProvideDSN,
		repository.NewUserRepository,
		usecase.NewUserUsecase,
		handler.NewUserHandler,
		route.NewRouter,
	)
	return nil, nil
}
