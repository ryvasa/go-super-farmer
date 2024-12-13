package route

import (
	"fmt"

	"github.com/gin-gonic/gin"
	handler "github.com/ryvasa/go-super-farmer/internal/delivery/http/handler"
	"github.com/ryvasa/go-super-farmer/internal/delivery/http/middleware"
	"github.com/ryvasa/go-super-farmer/pkg/auth/casbin"
	"github.com/ryvasa/go-super-farmer/pkg/auth/token"
	"github.com/ryvasa/go-super-farmer/pkg/env"
)

func NewRouter(handler *handler.Handlers) *gin.Engine {
	r := gin.Default()

	public := r.Group("/api")
	protected := r.Group("/api")

	env, err := env.LoadEnv()
	if err != nil {
		fmt.Printf("failed to load env: %v", err)
	}

	modelPath := "./pkg/auth/casbin/model.conf"
	policyPath := "./pkg/auth/casbin/policy.csv"

	enforcer, err := casbin.Init(modelPath, policyPath)
	if err != nil {
		fmt.Printf("failed to initialize casbin: %v", err)
	}

	token := token.NewToken(env)
	authMiddleware := middleware.NewAuthMiddleware(token)
	autzMiddleware := middleware.NewAutzMiddleware(enforcer)

	protected.Use(authMiddleware.Handle())
	protected.Use(autzMiddleware.Handle())

	RoleRoutes(public, protected, handler.RoleHandler)
	UserRoutes(public, protected, handler.UserHandler)
	LandRoutes(public, protected, handler.LandHandler)
	AuthRoutes(public, handler.AuthHandler)
	CommodityRoutes(public, protected, handler.CommodityHandler)
	LandCommodityRoutes(public, protected, handler.LandCommodityHandler)
	PriceRoutes(public, protected, handler.PriceHandler)
	ProvinceRoute(public, protected, handler.ProvinceHandler)
	CityRoute(public, protected, handler.CityHandler)
	RegionRoute(public, protected, handler.RegionHandler)
	DemandRoutes(public, protected, handler.DemandHandler)
	SupplyRoutes(public, protected, handler.SupplyHandler)
	return r

}
