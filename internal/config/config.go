package config

import (
	"log/slog"
	"strings"
	"time"
)

type Config struct {
	Domains   []string
	Every     time.Duration
	Sleep     time.Duration
	Threshold []int

	logLevel  string
	logFormat string

	TelegramChat  int64
	TelegramToken string

	GotifyURL   string
	GotifyToken string

	MetricsEnabled bool
	MetricsAddress string
}

func New() *Config {
	return &Config{
		Sleep:     3 * time.Second,
		Threshold: []int{1, 7},

		logLevel:  strings.ToLower(slog.LevelInfo.String()),
		logFormat: FormatAuto.String(),

		MetricsAddress: ":9090",
	}
}
