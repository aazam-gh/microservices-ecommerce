package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/yasinbozat/ecommerce-platform/services/search-service/internal/cache"
	"github.com/yasinbozat/ecommerce-platform/services/search-service/internal/config"
	"github.com/yasinbozat/ecommerce-platform/services/search-service/internal/handler"
	"github.com/yasinbozat/ecommerce-platform/services/search-service/internal/service"
)

func main() {
	cfg := config.Load()

	redisClient, err := config.NewRedis(cfg)
	if err != nil {
		log.Printf("redis unavailable: %v", err)
	}
	if redisClient != nil {
		defer redisClient.Close()
	}

	esClient, err := config.NewElasticsearch(cfg)
	if err != nil {
		log.Fatalf("cannot create elasticsearch client: %v", err)
	}

	searchCache := cache.NewRedisSearchCache(redisClient, cfg.Redis.KeyPrefix)
	searchService := service.NewSearchService(service.NewElasticsearchClient(esClient), searchCache, cfg)
	searchHandler := handler.NewSearchHandler(searchService)

	app := fiber.New()
	app.Get("/search", searchHandler.Search)

	app.Listen(":" + cfg.App.Port)
}
