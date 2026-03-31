package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      AppConfig
	DB       DBConfig
	Redis    RedisConfig
	Keycloak KeycloakConfig
	Jaeger   JaegerConfig
}

type AppConfig struct {
	Env  string
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type RedisConfig struct {
	Enabled      bool
	Addr         string
	Password     string
	DB           int
	KeyPrefix    string
	AuthCacheTTL time.Duration
}

type KeycloakConfig struct {
	Url          string
	Realm        string
	ClientID     string
	ClientSecret string
}

type JaegerConfig struct {
	Endpoint string
}

func Load() *Config {
	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8081"),
			Env:  getEnv("APP_ENV", "development"),
		}, DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "users_db"),
			User:     getEnv("DB_USER", "ecommerce"),
			Password: getEnv("DB_PASSWORD", "ecommerce123"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		}, Redis: RedisConfig{
			Enabled:      getEnvBool("REDIS_ENABLED", true),
			Addr:         getEnv("REDIS_ADDR", "redis:6379"),
			Password:     getEnv("REDIS_PASSWORD", "redis123"),
			DB:           getEnvInt("REDIS_DB", 0),
			KeyPrefix:    getEnv("REDIS_KEY_PREFIX", "user-service"),
			AuthCacheTTL: getEnvDuration("REDIS_AUTH_CACHE_TTL", 5*time.Minute),
		},
		Keycloak: KeycloakConfig{
			Url:          getEnv("KEYCLOAK_URL", "http://keycloak:8180"),
			Realm:        getEnv("KEYCLOAK_REALM", "ecommerce"),
			ClientID:     getEnv("KEYCLOAK_CLIENT_ID", "ecommerce-service"),
			ClientSecret: getEnv("KEYCLOAK_CLIENT_SECRET", "service-secret"),
		},
		Jaeger: JaegerConfig{
			Endpoint: getEnv("JAEGER_ENDPOINT", "jaeger:4318"), // http:// olmadan
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
