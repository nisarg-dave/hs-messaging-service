package config

import "testing"

func TestLoad(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_USER", "postgres")
	t.Setenv("POSTGRES_PASSWORD", "secret")
	t.Setenv("POSTGRES_DB", "messages")
	t.Setenv("SERVER_PORT", "8080")

	cfg := Load()

	expectedDatabaseURL := "host=localhost user=postgres password=secret dbname=messages port=5432 sslmode=disable"
	if cfg.DatabaseURL != expectedDatabaseURL {
		t.Errorf("DatabaseURL = %q, want %q", cfg.DatabaseURL, expectedDatabaseURL)
	}

	if cfg.ServerPort != "8080" {
		t.Errorf("ServerPort = %q, want %q", cfg.ServerPort, "8080")
	}
}
