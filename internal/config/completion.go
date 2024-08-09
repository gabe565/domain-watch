package config

import (
	"errors"

	"github.com/gabe565/domain-watch/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RegisterCompletions(cmd *cobra.Command) {
	if err := errors.Join(
		cmd.RegisterFlagCompletionFunc(CompletionFlag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"bash", "zsh", "fish", "powershell"}, cobra.ShellCompDirectiveNoFileComp
		}),

		cmd.RegisterFlagCompletionFunc(DomainsFlag, util.NoFileComp),
		cmd.RegisterFlagCompletionFunc(EveryFlag, util.NoFileComp),
		cmd.RegisterFlagCompletionFunc(SleepFlag, util.NoFileComp),
		cmd.RegisterFlagCompletionFunc(ThresholdFlag, util.NoFileComp),

		cmd.RegisterFlagCompletionFunc(MetricsEnabledFlag, util.NoFileComp),
		cmd.RegisterFlagCompletionFunc(MetricsAddressFlag, util.NoFileComp),

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

		cmd.RegisterFlagCompletionFunc(TelegramTokenFlag, util.NoFileComp),
		cmd.RegisterFlagCompletionFunc(TelegramChatFlag, util.NoFileComp),

		cmd.RegisterFlagCompletionFunc(GotifyURLFlag, util.NoFileComp),
		cmd.RegisterFlagCompletionFunc(GotifyTokenFlag, util.NoFileComp),
	); err != nil {
		panic(err)
	}
}
