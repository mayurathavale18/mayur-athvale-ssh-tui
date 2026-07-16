package config

import (
	"os"
)

type Config struct {
	Host       string
	Port       string
	HostKeyDir string
	DBPath     string
	// TGBotToken/TGContactChatID wire the Contact tab's form to a Telegram
	// bot's sendMessage API -- same bot as jobwatch (if you're reusing one),
	// but a separate chat id so portfolio-visitor messages don't mix with
	// job-tracker notifications. Contact submissions fail with a clear
	// in-TUI message if either is unset.
	TGBotToken      string
	TGContactChatID string
}

func Load() Config {
	return Config{
		Host:            getEnv("SSH_HOST", "0.0.0.0"),
		Port:            getEnv("SSH_PORT", "22"),
		HostKeyDir:      getEnv("HOST_KEY_DIR", ".ssh"),
		DBPath:          getEnv("DB_PATH", "data/analytics.db"),
		TGBotToken:      getEnv("TG_BOT_TOKEN", ""),
		TGContactChatID: getEnv("TG_CONTACT_CHAT_ID", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
