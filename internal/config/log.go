package config

import (
	log "github.com/sirupsen/logrus"
)

func InitLog(conf *Config) {
	parsedLevel, err := log.ParseLevel(conf.LogLevel)
	if err != nil {
		log.WithField("log-level", conf.LogLevel).Warn("invalid log level. defaulting to info.")
		conf.LogLevel = "info"
		parsedLevel = log.InfoLevel
	}
	log.SetLevel(parsedLevel)

	switch conf.LogFormat {
	case "text", "txt", "t":
		log.SetFormatter(&log.TextFormatter{})
	case "json", "j":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.WithField("log-format", conf.LogFormat).Warn("invalid log formatter. defaulting to text.")
		conf.LogFormat = "text"
		InitLog(conf)
	}
}
