package config

import (
	"os"
)

type Config struct {
	Host       string
	Port       string
	HostKeyDir string
	DBPath     string
}

func Load() Config {
	return Config{
		Host:       getEnv("SSH_HOST", "0.0.0.0"),
		Port:       getEnv("SSH_PORT", "22"),
		HostKeyDir: getEnv("HOST_KEY_DIR", ".ssh"),
		DBPath:     getEnv("DB_PATH", "data/analytics.db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
