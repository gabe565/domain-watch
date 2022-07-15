package cmd

func init() {
	Command.Flags().DurationVar(&conf.RunEvery, "every", 0, "enable cron mode and configure update interval")
	if err := Command.RegisterFlagCompletionFunc("every", noFileComp); err != nil {
		panic(err)
	}
}
