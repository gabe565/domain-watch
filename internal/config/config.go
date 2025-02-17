package config

import (
	"time"

	"gabe565.com/utils/slogx"
)

type Config struct {
	Domains   []string
	Every     time.Duration
	Sleep     time.Duration
	Threshold []int

	logLevel  slogx.Level
	logFormat slogx.Format

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

		logLevel:  slogx.LevelInfo,
		logFormat: slogx.FormatAuto,

		MetricsAddress: ":9090",
	}
}
