package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config guarda les variables carregades del fitxer .env
type Config struct {
    BackendUrl string
    ApiPort string
}

// Càrrega les variables d'entorn del fitxer .env
func Load() (*Config, error) {
    // Carregar les variables d'entorn des del fitxer .env
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error al carregar el fitxer .env: %v", err)
        return nil, err
    }

    // Crear una instància de Config amb les variables d'entorn carregades
    config := &Config{
        BackendUrl: os.Getenv("BACKEND_URL"),
        ApiPort: os.Getenv("API_PORT"),
    }

    // Retornar la configuració carregada
    return config, nil
}
