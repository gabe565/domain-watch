package cmd

import (
	"time"

	"github.com/gabe565/domain-watch/internal/util"
	"github.com/spf13/viper"
)

func init() {
	Command.Flags().DurationP("sleep", "s", 3*time.Second, "sleep time between queries to avoid rate limits")
	if err := viper.BindPFlag("sleep", Command.Flags().Lookup("sleep")); err != nil {
		panic(err)
	}
	if err := Command.RegisterFlagCompletionFunc("sleep", util.NoFileComp); err != nil {
		panic(err)
	}
}
