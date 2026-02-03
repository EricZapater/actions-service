package config

import (
	"os"
)

// Config guarda les variables carregades del fitxer .env
type Config struct {
    BackendUrl string
    ApiPort string
    RedisUrl string
    
    // Observability
    OtelEndpoint    string
    ServiceName     string
    ServiceVersion  string
    Environment     string
    LogLevel        string
}

// Càrrega les variables d'entorn del fitxer .env
func Load() (*Config, error) {
    // Carregar les variables d'entorn des del fitxer .env
    /*err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error al carregar el fitxer .env: %v", err)
        return nil, err
    }*/

    // Crear una instància de Config amb les variables d'entorn carregades
    config := &Config{
        BackendUrl: os.Getenv("BACKEND_URL"),
        ApiPort: os.Getenv("API_PORT"),
        RedisUrl: os.Getenv("REDIS_URL"),
        
        // Observability
        OtelEndpoint:   getEnvOrDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
        ServiceName:    getEnvOrDefault("OTEL_SERVICE_NAME", "actions-service"),
        ServiceVersion: getEnvOrDefault("OTEL_SERVICE_VERSION", "1.0.0"),
        Environment:    getEnvOrDefault("OTEL_ENVIRONMENT", "development"),
        LogLevel:       getEnvOrDefault("LOG_LEVEL", "info"),
    }

    // Retornar la configuració carregada
    return config, nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
