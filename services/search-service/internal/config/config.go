package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	App     AppConfig
	Elastic ElasticConfig
	Redis   RedisConfig
}

type AppConfig struct {
	Port string
}

type ElasticConfig struct {
	URL     string
	Index   string
	Timeout time.Duration
}

type RedisConfig struct {
	Enabled   bool
	Addr      string
	Password  string
	DB        int
	KeyPrefix string
	CacheTTL  time.Duration
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8086"),
		},
		Elastic: ElasticConfig{
			URL:     getEnv("ELASTICSEARCH_URL", "http://elasticsearch:9200"),
			Index:   getEnv("ELASTICSEARCH_INDEX", "products"),
			Timeout: getEnvDuration("ELASTICSEARCH_TIMEOUT", 5*time.Second),
		},
		Redis: RedisConfig{
			Enabled:   getEnvBool("REDIS_ENABLED", true),
			Addr:      getEnv("REDIS_ADDR", "redis:6379"),
			Password:  getEnv("REDIS_PASSWORD", "redis123"),
			DB:        getEnvInt("REDIS_DB", 0),
			KeyPrefix: getEnv("REDIS_KEY_PREFIX", "search-service"),
			CacheTTL:  getEnvDuration("REDIS_CACHE_TTL", 2*time.Minute),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}
