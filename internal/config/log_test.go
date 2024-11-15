package config

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_LogLevel(t *testing.T) {
	type fields struct {
		logLevel string
	}
	tests := []struct {
		name    string
		fields  fields
		want    slog.Level
		wantErr require.ErrorAssertionFunc
	}{
		{"debug", fields{"debug"}, slog.LevelDebug, require.NoError},
		{"info", fields{"info"}, slog.LevelInfo, require.NoError},
		{"warning", fields{"warn"}, slog.LevelWarn, require.NoError},
		{"error", fields{"error"}, slog.LevelError, require.NoError},
		{"unknown", fields{""}, slog.LevelInfo, require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{logLevel: tt.fields.logLevel}
			got, err := c.LogLevel()
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_LogFormat(t *testing.T) {
	type fields struct {
		logFormat string
	}
	tests := []struct {
		name    string
		fields  fields
		want    LogFormat
		wantErr require.ErrorAssertionFunc
	}{
		{"default", fields{"auto"}, FormatAuto, require.NoError},
		{"color", fields{"color"}, FormatColor, require.NoError},
		{"plain", fields{"plain"}, FormatPlain, require.NoError},
		{"json", fields{"json"}, FormatJSON, require.NoError},
		{"unknown", fields{""}, FormatAuto, require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{logFormat: tt.fields.logFormat}
			got, err := c.LogFormat()
			tt.wantErr(t, err)
			require.IsType(t, tt.want, got)
			assert.Equal(t, tt.want, got)
		})
	}
}
