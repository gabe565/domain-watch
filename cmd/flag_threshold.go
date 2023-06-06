package cmd

import (
	"github.com/gabe565/domain-watch/internal/util"
	"github.com/spf13/viper"
)

var defaultThreshold = []int{1, 7}

func init() {
	Command.Flags().IntSliceP("threshold", "t", defaultThreshold, "configure expiration notifications")
	if err := viper.BindPFlag("threshold", Command.Flags().Lookup("threshold")); err != nil {
		panic(err)
	}
	if err := Command.RegisterFlagCompletionFunc("threshold", util.NoFileComp); err != nil {
		panic(err)
	}
}
