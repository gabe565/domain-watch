package cmd

import (
	"time"

	"github.com/gabe565/domain-watch/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func registerSleepFlag(cmd *cobra.Command) {
	cmd.Flags().DurationP("sleep", "s", 3*time.Second, "sleep time between queries to avoid rate limits")
	if err := viper.BindPFlag("sleep", cmd.Flags().Lookup("sleep")); err != nil {
		panic(err)
	}
	if err := cmd.RegisterFlagCompletionFunc("sleep", util.NoFileComp); err != nil {
		panic(err)
	}
}
