package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Setenv("BACKEND_URL", "http://backend")
	t.Setenv("API_PORT", "8080")
	t.Setenv("REDIS_URL", "localhost:6379")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.BackendUrl != "http://backend" {
		t.Fatalf("unexpected backend url: %s", cfg.BackendUrl)
	}
	if cfg.ApiPort != "8080" {
		t.Fatalf("unexpected api port: %s", cfg.ApiPort)
	}
	if cfg.RedisUrl != "localhost:6379" {
		t.Fatalf("unexpected redis url: %s", cfg.RedisUrl)
	}

	os.Unsetenv("BACKEND_URL")
	os.Unsetenv("API_PORT")
	os.Unsetenv("REDIS_URL")

	cfg, err = Load()
	if err != nil {
		t.Fatalf("expected no error on missing env, got %v", err)
	}

	if cfg.BackendUrl != "" || cfg.ApiPort != "" || cfg.RedisUrl != "" {
		t.Fatalf("expected empty config when env not set, got %+v", cfg)
	}
}


