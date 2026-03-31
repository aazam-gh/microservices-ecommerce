package main

import (
	"context"
	"log"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/cache"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/config"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/handler"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/middleware"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/internal/service"
	"github.com/yasinbozat/ecommerce-platform/services/user-service/repository/postgres"
)

func main() {
	cfg := config.Load()
	db, err := config.NewDatabase(cfg)
	if err != nil {
		panic(err)
	}
	redisClient, err := config.NewRedis(cfg)
	if err != nil {
		log.Printf("redis unavailable, continuing without auth cache: %v", err)
	}
	if redisClient != nil {
		defer redisClient.Close()
	}
	tp, err := config.NewTracer(cfg)
	if err != nil {
		panic(err)
	}
	defer tp.Shutdown(context.Background())

	userRepo := postgres.NewUserRepository(db)
	addressRepo := postgres.NewAddressRepository(db)
	authCache := cache.NewRedisAuthCache(redisClient, cfg.Redis.KeyPrefix)

	userService := service.NewUserService(userRepo, addressRepo)
	authService := service.NewAuthService(userRepo, cfg, authCache)

	userHandler := handler.NewUserHandler(userService)
	addressHandler := handler.NewAddressHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	app := fiber.New()
	app.Use(logger.New())
	app.Use(otelfiber.Middleware())

	prometheus := fiberprometheus.New("user-service")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// internal route (no middleware)
	app.Get("/internal/auth/validate", authHandler.Validate)

	// api routes (middleware)
	api := app.Group("/api/v1", middleware.KeycloakMiddleware())

	// user routes
	api.Get("/users/me", userHandler.GetProfile)
	api.Put("/users/me", userHandler.UpdateProfile)

	// address routes
	api.Get("/users/me/addresses", addressHandler.List)
	api.Post("/users/me/addresses", addressHandler.Create)
	api.Put("/users/me/addresses/:id", addressHandler.Update)
	api.Delete("/users/me/addresses/:id", addressHandler.Delete)
	api.Patch("/users/me/addresses/:id/default", addressHandler.SetDefault)

	// health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Listen(":" + cfg.App.Port)
}
