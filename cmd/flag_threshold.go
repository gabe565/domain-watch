package cmd

import (
	"github.com/gabe565/domain-watch/internal/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var defaultThreshold = []int{1, 7}

func registerThresholdFlag(cmd *cobra.Command) {
	cmd.Flags().IntSliceP("threshold", "t", defaultThreshold, "configure expiration notifications")
	if err := viper.BindPFlag("threshold", cmd.Flags().Lookup("threshold")); err != nil {
		panic(err)
	}
	if err := cmd.RegisterFlagCompletionFunc("threshold", util.NoFileComp); err != nil {
		panic(err)
	}
}
