package config

import "context"

type ctxKey uint8

const configKey ctxKey = iota

func NewContext(ctx context.Context, conf *Config) context.Context {
	return context.WithValue(ctx, configKey, conf)
}

func FromContext(ctx context.Context) (*Config, bool) {
	conf, ok := ctx.Value(configKey).(*Config)
	return conf, ok
}
