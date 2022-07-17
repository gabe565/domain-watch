package cmd

import "github.com/spf13/viper"

func init() {
	Command.Flags().String("telegram-token", "", "Telegram token")
	if err := viper.BindPFlag("telegram.token", Command.Flags().Lookup("telegram-token")); err != nil {
		panic(err)
	}
	if err := Command.RegisterFlagCompletionFunc("telegram-token", noFileComp); err != nil {
		panic(err)
	}

	Command.Flags().Int64("telegram-chat", 0, "Telegram chat ID")
	if err := viper.BindPFlag("telegram.chat", Command.Flags().Lookup("telegram-chat")); err != nil {
		panic(err)
	}
	if err := Command.RegisterFlagCompletionFunc("telegram-chat", noFileComp); err != nil {
		panic(err)
	}
}
