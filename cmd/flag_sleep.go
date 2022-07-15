package cmd

import (
	"time"
)

func init() {
	Command.Flags().DurationVar(&conf.Sleep, "sleep", 3*time.Second, "sleep time between queries to avoid rate limits")
	if err := Command.RegisterFlagCompletionFunc("sleep", noFileComp); err != nil {
		panic(err)
	}
}
