package config

import (
	"log/slog"
	"strings"

	"gabe565.com/utils/must"
	"github.com/spf13/cobra"
)

func RegisterCompletions(cmd *cobra.Command) {
	must.Must(cmd.RegisterFlagCompletionFunc(FlagCompletion, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return []string{"bash", "zsh", "fish", "powershell"}, cobra.ShellCompDirectiveNoFileComp
	}))

	must.Must(cmd.RegisterFlagCompletionFunc(FlagDomains, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagEvery, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagSleep, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagThreshold, cobra.NoFileCompletions))

	must.Must(cmd.RegisterFlagCompletionFunc(FlagMetricsEnabled, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagMetricsAddress, cobra.NoFileCompletions))

	must.Must(cmd.RegisterFlagCompletionFunc(FlagLogLevel,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{
				strings.ToLower(slog.LevelDebug.String()),
				strings.ToLower(slog.LevelInfo.String()),
				strings.ToLower(slog.LevelWarn.String()),
				strings.ToLower(slog.LevelError.String()),
			}, cobra.ShellCompDirectiveNoFileComp
		}))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagLogFormat,
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"text", "json"}, cobra.ShellCompDirectiveNoFileComp
		}))

	must.Must(cmd.RegisterFlagCompletionFunc(FlagTelegramToken, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagTelegramChat, cobra.NoFileCompletions))

	must.Must(cmd.RegisterFlagCompletionFunc(FlagGotifyURL, cobra.NoFileCompletions))
	must.Must(cmd.RegisterFlagCompletionFunc(FlagGotifyToken, cobra.NoFileCompletions))
}
