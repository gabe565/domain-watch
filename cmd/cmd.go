package cmd

import (
	"errors"
	"time"

	"github.com/gabe565/domain-watch/internal/domain"
	"github.com/gabe565/domain-watch/internal/telegram"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Command = &cobra.Command{
	Use:               "domain-watch [flags] domain...",
	DisableAutoGenTag: true,
	PreRunE:           preRun,
	RunE:              run,
	ValidArgsFunction: noFileComp,
}

func init() {
	cobra.OnInitialize(initViper, initLog)
}

var domainNames []string

func preRun(cmd *cobra.Command, args []string) (err error) {
	if completionFlag != "" {
		return completion(cmd, domainNames)
	}

	if v := viper.GetStringSlice("domains"); v != nil {
		args = append(domainNames, v...)
	}
	domainNames = args

	if len(domainNames) == 0 {
		return errors.New("missing domain")
	}

	token := viper.GetString("telegram.token")
	if token != "" {
		if !viper.IsSet("telegram.chat") {
			return errors.New("telegram token flag requires --telegram-chat to be set")
		}

		if err := telegram.Login(token); err != nil {
			return err
		}
	}
	return nil
}

func run(cmd *cobra.Command, _ []string) (err error) {
	cmd.SilenceUsage = true

	domains := domain.Domains{
		Sleep:   viper.GetDuration("sleep"),
		Domains: make([]*domain.Domain, 0, len(domainNames)),
	}
	for _, domainName := range domainNames {
		domains.Add(domain.Domain{Name: domainName})
	}

	domains.Tick()

	every := viper.GetDuration("every")
	if every != 0 {
		log.Info("running as cron")

		ticker := time.NewTicker(every)
		for range ticker.C {
			domains.Tick()
		}
	}

	return nil
}
