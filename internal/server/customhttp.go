package server

import (
	_ "bus-service/docs"
	"bus-service/internal/conf"
	"bus-service/internal/data"
	"bus-service/internal/route"
	"bus-service/pkg/customhttp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func NewCustomHttp(
	c *conf.Server,
	bus *route.BusRouter,
	keycloak *data.KeycloakAPI,
	route *route.RouteRouter,
	driver *route.DriverRoute,
	logger log.Logger) *customhttp.CustomHTTP {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Custom.Network != "" {
		opts = append(opts, http.Network(c.Custom.Network))
	}
	if c.Custom.Addr != "" {
		opts = append(opts, http.Address(c.Custom.Addr))
	}
	if c.Custom.Timeout != nil {
		opts = append(opts, http.Timeout(c.Custom.Timeout.AsDuration()))
	}
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"}
	config.AllowCredentials = true
	r.Use(cors.New(config))
	srv := http.NewServer(opts...)

	srv.HandlePrefix("/", r)
	return &customhttp.CustomHTTP{Http: srv}
}
