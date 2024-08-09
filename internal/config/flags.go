package config

import (
	"github.com/spf13/cobra"
)

const (
	CompletionFlag = "completion"

	DomainsFlag   = "domains"
	EveryFlag     = "every"
	SleepFlag     = "sleep"
	ThresholdFlag = "threshold"

	LogFormatFlag = "log-format"
	LogLevelFlag  = "log-level"

	MetricsEnabledFlag = "metrics-enabled"
	MetricsAddressFlag = "metrics-address"

	TelegramChatFlag  = "telegram-chat"
	TelegramTokenFlag = "telegram-token"

	GotifyURLFlag   = "gotify-url"
	GotifyTokenFlag = "gotify-token"
)

func (c *Config) RegisterFlags(cmd *cobra.Command) {
	fs := cmd.Flags()
	fs.StringVar(&c.Completion, CompletionFlag, c.Completion, "Output command-line completion code for the specified shell. Can be 'bash', 'zsh', 'fish', or 'powershell'.")

	fs.StringSliceVar(&c.Domains, DomainsFlag, c.Domains, "List of domains to watch")
	fs.DurationVarP(&c.Every, EveryFlag, "e", c.Every, "enable cron mode and configure update interval")
	fs.DurationVarP(&c.Sleep, SleepFlag, "s", c.Sleep, "sleep time between queries to avoid rate limits")
	fs.IntSliceVarP(&c.Threshold, ThresholdFlag, "t", c.Threshold, "configure expiration notifications")

	fs.StringVarP(&c.LogLevel, LogLevelFlag, "l", c.LogLevel, "log level (trace, debug, info, warning, error, fatal, panic)")
	fs.StringVar(&c.LogFormat, LogFormatFlag, c.LogFormat, "log formatter (text, json)")

	fs.BoolVar(&c.MetricsEnabled, MetricsEnabledFlag, c.MetricsEnabled, "Enables Prometheus metrics API")
	fs.StringVar(&c.MetricsAddress, MetricsAddressFlag, c.MetricsAddress, "Prometheus metrics API listen address")

	fs.StringVar(&c.TelegramToken, TelegramTokenFlag, c.TelegramToken, "Telegram token")
	fs.Int64Var(&c.TelegramChat, TelegramChatFlag, c.TelegramChat, "Telegram chat ID")
	cmd.MarkFlagsRequiredTogether(TelegramTokenFlag, TelegramChatFlag)

	fs.StringVar(&c.GotifyURL, GotifyURLFlag, c.GotifyURL, "Gotify URL (include https:// and port if non-standard)")
	fs.StringVar(&c.GotifyToken, GotifyTokenFlag, c.GotifyToken, "Gotify app token")
	cmd.MarkFlagsRequiredTogether(GotifyURLFlag, GotifyTokenFlag)
}
