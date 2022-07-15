package cmd

func init() {
	Command.Flags().StringVar(&conf.Token, "telegram-token", "", "Telegram token")
	if err := Command.RegisterFlagCompletionFunc("telegram-token", noFileComp); err != nil {
		panic(err)
	}

	Command.Flags().Int64Var(&conf.ChatId, "telegram-chat", 0, "Telegram chat ID")
	if err := Command.RegisterFlagCompletionFunc("telegram-chat", noFileComp); err != nil {
		panic(err)
	}
}
