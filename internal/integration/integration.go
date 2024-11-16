package integration

import (
	"context"
	"errors"
	"log/slog"

	"gabe565.com/domain-watch/internal/config"
	"gabe565.com/domain-watch/internal/integration/gotify"
	"gabe565.com/domain-watch/internal/integration/telegram"
	"gabe565.com/domain-watch/internal/util"
)

type Integration interface {
	Setup(ctx context.Context, conf *config.Config) error
	Send(ctx context.Context, text string) error
}

type Integrations map[string]Integration

func Default() Integrations {
	return map[string]Integration{
		"telegram": &telegram.Telegram{},
		"gotify":   &gotify.Gotify{},
	}
}

func Setup(ctx context.Context, conf *config.Config) (Integrations, error) {
	var configured uint8

	integrations := Default()

	for _, integration := range integrations {
		err := integration.Setup(ctx, conf)
		if err != nil {
			if errors.Is(err, util.ErrNotConfigured) {
				continue
			}
			return nil, err
		}
		configured++
	}

	if configured == 0 {
		slog.Warn("No integrations were configured")
	}

	return integrations, nil
}

func (i Integrations) Send(ctx context.Context, message string) {
	for name, integration := range i {
		if err := integration.Send(ctx, message); err != nil {
			slog.Error("Failed to send message", "integration", name, "error", err)
		}
	}
}
