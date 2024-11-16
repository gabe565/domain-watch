package config

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/spf13/cobra"
)

func RegisterCompletions(cmd *cobra.Command) {
	if err := errors.Join(
		cmd.RegisterFlagCompletionFunc(FlagCompletion, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"bash", "zsh", "fish", "powershell"}, cobra.ShellCompDirectiveNoFileComp
		}),

		cmd.RegisterFlagCompletionFunc(FlagDomains, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(FlagEvery, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(FlagSleep, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(FlagThreshold, cobra.NoFileCompletions),

		cmd.RegisterFlagCompletionFunc(FlagMetricsEnabled, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(FlagMetricsAddress, cobra.NoFileCompletions),

		cmd.RegisterFlagCompletionFunc(FlagLogLevel,
			func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
				return []string{
					strings.ToLower(slog.LevelDebug.String()),
					strings.ToLower(slog.LevelInfo.String()),
					strings.ToLower(slog.LevelWarn.String()),
					strings.ToLower(slog.LevelError.String()),
				}, cobra.ShellCompDirectiveNoFileComp
			}),
		cmd.RegisterFlagCompletionFunc(FlagLogFormat,
			func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
				return []string{"text", "json"}, cobra.ShellCompDirectiveNoFileComp
			}),

		cmd.RegisterFlagCompletionFunc(FlagTelegramToken, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(FlagTelegramChat, cobra.NoFileCompletions),

		cmd.RegisterFlagCompletionFunc(FlagGotifyURL, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(FlagGotifyToken, cobra.NoFileCompletions),
	); err != nil {
		panic(err)
	}
}
