package integration

import (
	"errors"
	"log/slog"

	"gabe565.com/domain-watch/internal/config"
	"gabe565.com/domain-watch/internal/util"
)

type Integration interface {
	Setup(*config.Config) error
	Send(string) error
}

var Default = map[string]Integration{
	"telegram": &Telegram{},
	"gotify":   &Gotify{},
}

func Setup(conf *config.Config) error {
	var configured uint8

	for _, integration := range Default {
		err := integration.Setup(conf)
		if err != nil {
			if errors.Is(err, util.ErrNotConfigured) {
				continue
			}
			return err
		}
		configured += 1
	}

	if configured == 0 {
		slog.Warn("No integrations were configured")
	}

	return nil
}

func Send(message string) {
	for name, integration := range Default {
		if err := integration.Send(message); err != nil {
			slog.Error("Failed to send message", "integration", name, "error", err)
		}
	}
}

func Get(key string) Integration {
	return Default[key]
}
