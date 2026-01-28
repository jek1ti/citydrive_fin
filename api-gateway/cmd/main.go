package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jekiti/citydrive/api-gateway/internal/config"
	"github.com/jekiti/citydrive/api-gateway/internal/handler"
	"github.com/jekiti/citydrive/api-gateway/internal/middleware"
	"github.com/jekiti/citydrive/api-gateway/internal/service"
	"github.com/jekiti/citydrive/pkg/logger"
)

func main() {
	cfg := config.LoadGatewayConfig()
	if err := cfg.Validate(); err != nil {
		panic("Invalid configuration: " + err.Error())
	}
	log := logger.SetupLogger(cfg.App.LogLevel)

	log.Info("API Gateway starting...",
		"env", cfg.App.Env,
		"http_port", cfg.HTTP.Port,
		"log_level", cfg.App.LogLevel)

	log.Info("gateway config",
		"module", "bootstrap",
		"auth_target", cfg.GRPC.AuthAddr,
		"dial_timeout", cfg.GRPC.DialTimeout,
	)

	telemetryClient, err := service.NewTelemetryClient(&cfg.GRPC, log)
	if err != nil {
		log.Error("failed to create telemetry client:", "error", err)
		panic("telemetry client failed")
	}
	log.Info("TelemetryClient created successful")
	defer telemetryClient.Close()

	authClient, err := service.NewAuthClient(&cfg.GRPC)
	if err != nil {
		log.Error("failed to create Auth client:", "error", err)
		panic("Auth client failed")
	}
	log.Info("AuthClient created successful")
	defer authClient.Close()

	adminClient, err := service.NewAdminClient(&cfg.GRPC)
	if err != nil {
		log.Error("failed to create admin client:", "error", err)
		panic("admin client failed")
	}
	log.Info("AdminClient created successful")
	defer adminClient.Close()

	carHandler := handler.NewTelemetryHandler(telemetryClient)
	authHandler := handler.NewAuthHandler(authClient)
	adminHandler := handler.NewAdminHandler(adminClient)
	router := gin.Default()

	router.Use(middleware.TracingMiddleware(cfg.Tracing.HeaderName))

	carInfoGroup := router.Group("/api/v1")
	{
		carInfoGroup.Use(middleware.RequireCarAuth(cfg.JWT.CarSecretKey))
		carInfoGroup.PUT("/car-info", carHandler.PutCarInfo)
	}
	authGroup := router.Group("/v1/user")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/register", authHandler.Register)
	}

	adminGroup := router.Group("/api/v1/cars")
	{
		adminGroup.Use(middleware.RequireAuth(cfg.JWT.SecretKey))
		adminGroup.GET("/now", adminHandler.GetCarsNow)
		adminGroup.GET("/:id", adminHandler.GetCar)
		adminGroup.GET("/history", adminHandler.GetCarsHistory)
		adminGroup.GET("/:id/history", adminHandler.GetCarHistory)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "api-gateway",
		})
	})

	log.Info("Starting HTTP server", "port", cfg.HTTP.Port)

	if err := router.Run(":" + cfg.HTTP.Port); err != nil {
		log.Error("Failed to start HTTP server", "error", err)
		panic("HTTP server failed: " + err.Error())
	}
}
