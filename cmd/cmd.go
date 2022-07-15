package cmd

import (
	"errors"
	"github.com/gabe565/domain-watch/internal/config"
	"github.com/gabe565/domain-watch/internal/domain"
	"github.com/gabe565/domain-watch/internal/telegram"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

var Command = &cobra.Command{
	Use:     "domain-watch",
	PreRunE: preRun,
	RunE:    run,
}

var conf config.Config

func init() {
	Command.Flags().DurationVar(&conf.RunEvery, "every", 0, "enable cron mode and configure update interval")
	Command.Flags().DurationVar(&conf.Sleep, "sleep", 3*time.Second, "sleep time between queries to avoid rate limits")
	Command.Flags().StringVar(&conf.Token, "telegram-token", "", "Telegram token")
	Command.Flags().Int64Var(&conf.ChatId, "telegram-chat", 0, "Telegram chat ID")
	cobra.OnInitialize(initLog)
}

func preRun(cmd *cobra.Command, domainNames []string) (err error) {
	if completionFlag != "" {
		return completion(cmd, domainNames)
	}

	if conf.Token != "" {
		if conf.ChatId == 0 {
			return errors.New("telegram token flag requires --telegram-chat to be set")
		}

		if err := telegram.Login(conf.Token, conf.ChatId); err != nil {
			return err
		}
	}
	return nil
}

func run(cmd *cobra.Command, domainNames []string) (err error) {
	cmd.SilenceUsage = true

	domains := make(domain.Domains, 0, len(domainNames))
	for i, domainName := range domainNames {
		var sleep time.Duration
		if i != 0 {
			sleep = conf.Sleep
		}
		d := domain.Domain{
			Name:  domainName,
			Sleep: sleep,
		}
		domains = append(domains, d)
	}

	domains.Tick()

	if conf.RunEvery != 0 {
		log.Info("running as cron")

		c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
		_, err := c.AddFunc("@every "+conf.RunEvery.String(), domains.Tick)
		if err != nil {
			log.WithError(err).Error("failed to register job")
			return err
		}
		c.Run()
	}

	return nil
}
