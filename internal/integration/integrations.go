package integration

import (
	"context"
	"errors"
	"log/slog"
	"slices"

	"gabe565.com/domain-watch/internal/config"
	"gabe565.com/domain-watch/internal/integration/gotify"
	"gabe565.com/domain-watch/internal/integration/telegram"
	"gabe565.com/domain-watch/internal/util"
)

type Integrations []Integration

func All() Integrations {
	return Integrations{
		&telegram.Telegram{},
		&gotify.Gotify{},
	}
}

func Setup(ctx context.Context, conf *config.Config) (Integrations, error) {
	all := All()
	integrations := make(Integrations, 0, len(all))
	for _, integration := range all {
		if err := integration.Setup(ctx, conf); err != nil {
			if errors.Is(err, util.ErrNotConfigured) {
				continue
			}
			return nil, err
		}
		integrations = append(integrations, integration)
	}

	if len(integrations) == 0 {
		slog.Warn("No integrations were configured")
	}

	integrations = slices.Clip(integrations)
	return integrations, nil
}

func (i Integrations) Send(ctx context.Context, message string) {
	for _, integration := range i {
		if err := integration.Send(ctx, message); err != nil {
			slog.Error("Failed to send message", "integration", integration.Name(), "error", err)
		}
	}
}
