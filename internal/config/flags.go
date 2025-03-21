package config

import (
	"strings"

	"gabe565.com/utils/slogx"
	"github.com/spf13/cobra"
)

const (
	FlagDomains   = "domains"
	FlagEvery     = "every"
	FlagSleep     = "sleep"
	FlagThreshold = "threshold"

	FlagLogFormat = "log-format"
	FlagLogLevel  = "log-level"

	FlagMetricsEnabled = "metrics-enabled"
	FlagMetricsAddress = "metrics-address"

	FlagTelegramChat  = "telegram-chat"
	FlagTelegramToken = "telegram-token"

	FlagGotifyURL   = "gotify-url"
	FlagGotifyToken = "gotify-token"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	fs := cmd.Flags()

	fs.StringSliceVar(&c.Domains, FlagDomains, c.Domains, "List of domains to watch")
	fs.DurationVarP(&c.Every, FlagEvery, "e", c.Every, "Enable cron mode and configure update interval")
	fs.DurationVarP(&c.Sleep, FlagSleep, "s", c.Sleep, "Sleep time between queries to avoid rate limits")
	fs.IntSliceVarP(&c.Threshold, FlagThreshold, "t", c.Threshold, "Configure expiration notifications")

	fs.VarP(&c.logLevel, FlagLogLevel, "l", "Log level (one of "+strings.Join(slogx.LevelStrings(), ", ")+")")
	fs.Var(&c.logFormat, FlagLogFormat, "Log formatter (one of "+strings.Join(slogx.FormatStrings(), ", ")+")")

	fs.BoolVar(&c.MetricsEnabled, FlagMetricsEnabled, c.MetricsEnabled, "Enables Prometheus metrics API")
	fs.StringVar(&c.MetricsAddress, FlagMetricsAddress, c.MetricsAddress, "Prometheus metrics API listen address")

	fs.StringVar(&c.TelegramToken, FlagTelegramToken, c.TelegramToken, "Telegram token")
	fs.Int64Var(&c.TelegramChat, FlagTelegramChat, c.TelegramChat, "Telegram chat ID")
	cmd.MarkFlagsRequiredTogether(FlagTelegramToken, FlagTelegramChat)

	fs.StringVar(&c.GotifyURL, FlagGotifyURL, c.GotifyURL, "Gotify URL (include https:// and port if non-standard)")
	fs.StringVar(&c.GotifyToken, FlagGotifyToken, c.GotifyToken, "Gotify app token")
	cmd.MarkFlagsRequiredTogether(FlagGotifyURL, FlagGotifyToken)
}
