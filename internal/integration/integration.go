package integration

import (
	"context"

	"gabe565.com/domain-watch/internal/config"
)

type Integration interface {
	Name() string
	Setup(ctx context.Context, conf *config.Config) error
	Send(ctx context.Context, text string) error
}
