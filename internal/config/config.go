package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Token        string
	ChatID       int64
	TickMs       int64
	RetryEveryMs int64
	MaxAttempts  int
}

func Load() (*Config, error) {
	token := os.Getenv("TG_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("TG_TOKEN not set")
	}
	
	chatIDStr := os.Getenv("TG_CHAT_ID")
	if chatIDStr == "" {
		return nil, fmt.Errorf("TG_CHAT_ID not set")
	}

	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid TG_CHAT_ID: %w", err)
	}
	
	return &Config{
		Token:        token,
		ChatID:       chatID,
		TickMs:       30 * 1000,        // 30 detik
		RetryEveryMs: 5 * 60 * 1000,    // 5 menit
		MaxAttempts:  6,
	}, nil
}
