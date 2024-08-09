package config

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RegisterCompletions(cmd *cobra.Command) {
	if err := errors.Join(
		cmd.RegisterFlagCompletionFunc(CompletionFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"bash", "zsh", "fish", "powershell"}, cobra.ShellCompDirectiveNoFileComp
		}),

		cmd.RegisterFlagCompletionFunc(DomainsFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(EveryFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(SleepFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(ThresholdFlag, cobra.NoFileCompletions),

		cmd.RegisterFlagCompletionFunc(MetricsEnabledFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(MetricsAddressFlag, cobra.NoFileCompletions),

		cmd.RegisterFlagCompletionFunc(LogLevelFlag,
			func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return []string{
					log.TraceLevel.String(),
					log.DebugLevel.String(),
					log.InfoLevel.String(),
					log.WarnLevel.String(),
					log.ErrorLevel.String(),
					log.FatalLevel.String(),
					log.PanicLevel.String(),
				}, cobra.ShellCompDirectiveNoFileComp
			}),
		cmd.RegisterFlagCompletionFunc(LogFormatFlag,
			func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				return []string{"text", "json"}, cobra.ShellCompDirectiveNoFileComp
			}),

		cmd.RegisterFlagCompletionFunc(TelegramTokenFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(TelegramChatFlag, cobra.NoFileCompletions),

		cmd.RegisterFlagCompletionFunc(GotifyURLFlag, cobra.NoFileCompletions),
		cmd.RegisterFlagCompletionFunc(GotifyTokenFlag, cobra.NoFileCompletions),
	); err != nil {
		panic(err)
	}
}
