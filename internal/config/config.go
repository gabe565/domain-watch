package config

import "time"

type Config struct {
	Completion string

	Domains   []string
	Every     time.Duration
	Sleep     time.Duration
	Threshold []int

	LogLevel  string
	LogFormat string

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

		LogLevel:  "info",
		LogFormat: "text",

		MetricsAddress: ":9090",
	}
}
