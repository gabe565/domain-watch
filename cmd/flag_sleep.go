package cmd

import (
	"github.com/spf13/viper"
	"time"
)

func init() {
	Command.Flags().Duration("sleep", 3*time.Second, "sleep time between queries to avoid rate limits")
	if err := viper.BindPFlag("sleep", Command.Flags().Lookup("sleep")); err != nil {
		panic(err)
	}
	if err := Command.RegisterFlagCompletionFunc("sleep", noFileComp); err != nil {
		panic(err)
	}
}
