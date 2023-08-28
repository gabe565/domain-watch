package cmd

import (
	"github.com/gabe565/domain-watch/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func registerEveryFlag(cmd *cobra.Command) {
	cmd.Flags().DurationP("every", "e", 0, "enable cron mode and configure update interval")
	if err := viper.BindPFlag("every", cmd.Flags().Lookup("every")); err != nil {
		panic(err)
	}
	if err := cmd.RegisterFlagCompletionFunc("every", util.NoFileComp); err != nil {
		panic(err)
	}
}
