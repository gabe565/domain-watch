package cmd

import "github.com/spf13/viper"

func init() {
	Command.Flags().Duration("every", 0, "enable cron mode and configure update interval")
	if err := viper.BindPFlag("every", Command.Flags().Lookup("every")); err != nil {
		panic(err)
	}
	if err := Command.RegisterFlagCompletionFunc("every", noFileComp); err != nil {
		panic(err)
	}
}
