package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	var err error
	Command.PersistentFlags().String("log-level", "info", "log level (trace, debug, info, warning, error, fatal, panic)")
	err = Command.RegisterFlagCompletionFunc(
		"log-level",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				log.TraceLevel.String(),
				log.DebugLevel.String(),
				log.InfoLevel.String(),
				log.WarnLevel.String(),
				log.ErrorLevel.String(),
				log.FatalLevel.String(),
				log.PanicLevel.String(),
			}, cobra.ShellCompDirectiveNoFileComp
		})
	if err != nil {
		panic(err)
	}

	Command.PersistentFlags().String("log-format", "text", "log formatter (text, json)")
	err = Command.RegisterFlagCompletionFunc(
		"log-format",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"text", "json"}, cobra.ShellCompDirectiveNoFileComp
		})
	if err != nil {
		panic(err)
	}
}

func initLog() {
	logLevel, err := Command.PersistentFlags().GetString("log-level")
	if err != nil {
		panic(err)
	}
	parsedLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		log.WithField("log-level", logLevel).Warn("invalid log level. defaulting to info.")
		err = Command.PersistentFlags().Set("log-level", "info")
		if err != nil {
			panic(err)
		}
		parsedLevel = log.InfoLevel
	}
	log.SetLevel(parsedLevel)

	logFormat, err := Command.PersistentFlags().GetString("log-format")
	if err != nil {
		panic(err)
	}
	switch logFormat {
	case "text", "txt", "t":
		log.SetFormatter(&log.TextFormatter{})
	case "json", "j":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.WithField("log-format", logFormat).Warn("invalid log formatter. defaulting to text.")
		err = Command.PersistentFlags().Set("log-format", "text")
		if err != nil {
			panic(err)
		}
	}
}
