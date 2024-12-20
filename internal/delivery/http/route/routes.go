package route

import (
	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/middleware"
	"github.com/ryvasa/go-super-farmer/pkg/auth/casbin"
	"github.com/ryvasa/go-super-farmer/pkg/auth/token"
	"github.com/ryvasa/go-super-farmer/pkg/env"
)

type Router interface {
	Register(public, protected *gin.RouterGroup)
}

func NewRouter(handlers *handler.Handlers) *gin.Engine {
	r := gin.Default()

	public := r.Group("/api")
	protected := r.Group("/api")

	env, err := env.LoadEnv()
	if err != nil {
		panic(err)
	}

	// Setup middleware
	enforcer, err := casbin.Init("./pkg/auth/casbin/model.conf", "./pkg/auth/casbin/policy.csv")
	if err != nil {
		panic(err)
	}

	tokenService := token.NewToken(env)
	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	autzMiddleware := middleware.NewAutzMiddleware(enforcer)

	protected.Use(authMiddleware.Handle())
	protected.Use(autzMiddleware.Handle())

	routes := []Router{
		NewAuthRoute(handlers.AuthHandler),
		NewUserRoute(handlers.UserHandler),
		NewLandRoute(handlers.LandHandler),
		NewCommodityRoute(handlers.CommodityHandler),
		NewPriceRoute(handlers.PriceHandler),
		NewProvinceRoute(handlers.ProvinceHandler),
		NewCityRoute(handlers.CityHandler),
		NewDemandRoute(handlers.DemandHandler),
		NewSupplyRoute(handlers.SupplyHandler),
		NewHarvestRoute(handlers.HarvestHandler),
		NewLandCommodityRoute(handlers.LandCommodityHandler),
		NewRoleRoute(handlers.RoleHandler),
	}

	// Register all routes
	for _, router := range routes {
		router.Register(public, protected)
	}

	return r
}
