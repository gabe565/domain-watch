package cmd

import (
	"github.com/gabe565/domain-watch/internal/util"
	"github.com/spf13/viper"
)

func init() {
	Command.Flags().DurationP("every", "e", 0, "enable cron mode and configure update interval")
	if err := viper.BindPFlag("every", Command.Flags().Lookup("every")); err != nil {
		panic(err)
	}
	if err := Command.RegisterFlagCompletionFunc("every", util.NoFileComp); err != nil {
		panic(err)
	}
}
