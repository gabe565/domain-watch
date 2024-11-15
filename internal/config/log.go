package config

import (
	"io"
	"log/slog"
	"strings"
	"time"

	"gabe565.com/utils/termx"
	"github.com/lmittmann/tint"
)

//go:generate go run github.com/dmarkham/enumer -type LogFormat -trimprefix Format -transform lower

type LogFormat uint8

const (
	FormatAuto LogFormat = iota
	FormatColor
	FormatPlain
	FormatJSON
)

func (c *Config) LogLevel() (slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(c.logLevel))
	if err != nil {
		level = slog.LevelInfo
	}
	return level, err
}

func (c *Config) LogFormat() (LogFormat, error) {
	var format LogFormat
	format, err := LogFormatString(c.logFormat)
	if err != nil {
		format = FormatAuto
	}
	return format, err
}

func (c *Config) InitLog(w io.Writer) {
	level, err := c.LogLevel()
	if err != nil {
		defer func(val string) {
			slog.Warn("Invalid log level. Defaulting to info.", "value", val)
		}(c.logLevel)
		c.logLevel = strings.ToLower(level.String())
	}

	format, err := c.LogFormat()
	if err != nil {
		defer func(val string) {
			slog.Warn("Invalid log format. Defaulting to auto.", "value", val)
		}(c.logFormat)
		c.logFormat = format.String()
	}

	InitLog(w, level, format)
}

func InitLog(w io.Writer, level slog.Level, format LogFormat) {
	switch format {
	case FormatJSON:
		slog.SetDefault(slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: level,
		})))
	default:
		var color bool
		switch format {
		case FormatAuto:
			color = termx.IsColor(w)
		case FormatColor:
			color = true
		}

		slog.SetDefault(slog.New(
			tint.NewHandler(w, &tint.Options{
				Level:      level,
				TimeFormat: time.DateTime,
				NoColor:    !color,
			}),
		))
	}
}
