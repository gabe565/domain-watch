package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func registerLogFlags(cmd *cobra.Command) {
	var err error
	cmd.Flags().StringP("log-level", "l", "info", "log level (trace, debug, info, warning, error, fatal, panic)")
	if err := viper.BindPFlag("log.level", cmd.Flags().Lookup("log-level")); err != nil {
		panic(err)
	}
	err = cmd.RegisterFlagCompletionFunc(
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

	cmd.Flags().String("log-format", "text", "log formatter (text, json)")
	if err := viper.BindPFlag("log.format", cmd.Flags().Lookup("log-format")); err != nil {
		panic(err)
	}
	err = cmd.RegisterFlagCompletionFunc(
		"log-format",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{"text", "json"}, cobra.ShellCompDirectiveNoFileComp
		})
	if err != nil {
		panic(err)
	}
}

func initLog(cmd *cobra.Command) {
	logLevel := viper.GetString("log.level")
	parsedLevel, err := log.ParseLevel(logLevel)
	if err != nil {
		log.WithField("log-level", logLevel).Warn("invalid log level. defaulting to info.")
		err = cmd.Flags().Set("log-level", "info")
		if err != nil {
			panic(err)
		}
		parsedLevel = log.InfoLevel
	}
	log.SetLevel(parsedLevel)

	logFormat := viper.GetString("log.format")
	switch logFormat {
	case "text", "txt", "t":
		log.SetFormatter(&log.TextFormatter{})
	case "json", "j":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.WithField("log-format", logFormat).Warn("invalid log formatter. defaulting to text.")
		err = cmd.Flags().Set("log-format", "text")
		if err != nil {
			panic(err)
		}
	}
}
