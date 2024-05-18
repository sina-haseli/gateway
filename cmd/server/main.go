// cmd/server/main.go
package main

import (
	"gateway/config"
	"gateway/internal/adapters/http"
	"gateway/internal/application"
	"gateway/internal/infrastructure"
	"gateway/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"time"
)

func serve() http.Handler {
	ca := config.InitializeConfig()

	services := infrastructure.NewServicesClient(infrastructure.DomainsToService(ca.Config.EntryPoints.Http.Domains, ca.Config.EntryPoints.Address))

	infra := infrastructure.NewInfrastructure(infrastructure.NewGRPCClient(ca.Config.GRPCClient.Host, ca.Config.GRPCClient.Port), services)
	app := application.NewService(infra, ca.Config.CircuitBreaker.FailureThreshold, ca.Config.CircuitBreaker.SuccessThreshold, time.Duration(ca.Config.CircuitBreaker.Timeout), services)
	//pkg.StartGRPCServer(ca.Config.GRPCClient.Host, ca.Config.GRPCClient.Port)
	return http.NewHandler(app)
}

func main() {
	logger.InitLogger()
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler := serve()

	// Routes
	e.GET("/components/", handler.GetRequest)
	e.GET("/grpc/", handler.GetGRPCRequest)

	// Start server
	e.Logger.Fatal(e.Start("0.0.0.0:9090"))
}
