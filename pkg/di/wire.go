//go:build wireinject
// +build wireinject

package di

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/route"
	"github.com/ryvasa/go-super-farmer/internal/repository"
	"github.com/ryvasa/go-super-farmer/internal/usecase"
	"github.com/ryvasa/go-super-farmer/pkg/auth/token"
	"github.com/ryvasa/go-super-farmer/pkg/database"
	"github.com/ryvasa/go-super-farmer/pkg/env"
)

var roleSet = wire.NewSet(
	repository.NewRoleRepository,
	usecase.NewRoleUsecase,
	handler.NewRoleHandler,
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	usecase.NewUserUsecase,
	handler.NewUserHandler,
)

var landSet = wire.NewSet(
	repository.NewLandRepository,
	usecase.NewLandUsecase,
	handler.NewLandHandler,
)

var authSet = wire.NewSet(
	usecase.NewAuthUsecase,
	handler.NewAuthHandler,
)

var utilSet = wire.NewSet(
	utils.NewTokenUtil
)

// var tokenSet = wire.NewSet(
// 	token.NewToken,
// )

// var authUtilSet = wire.NewSet(
// 	utils.NewAuthUtil, // gunakan konstruktor yang sesuai
// )

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
		// authUtilSet,
	)
	return nil, nil
}
