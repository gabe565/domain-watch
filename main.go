package main

import (
	"github.com/gabe565/domain-expiration-notifier/internal/config"
	"github.com/gabe565/domain-expiration-notifier/internal/domain"
	"github.com/gabe565/domain-expiration-notifier/internal/telegram"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

var cmd = &cobra.Command{
	Use:     "domain-expiration-notifier",
	PreRunE: preRun,
	RunE:    run,
}

var conf config.Config

func init() {
	cmd.Flags().StringVar(&conf.RunEvery, "every", "", "enable cron mode and configure update interval")
	cmd.Flags().DurationVar(&conf.Sleep, "sleep", 3*time.Second, "sleep time between queries to avoid rate limits")
	cmd.Flags().StringVar(&conf.Token, "telegram-token", "", "Telegram token")
	cmd.Flags().Int64Var(&conf.ChatId, "telegram-chat", 0, "Telegram chat ID")
}

func preRun(cmd *cobra.Command, domainNames []string) (err error) {
	if conf.Token != "" && conf.ChatId != 0 {
		if err := telegram.Login(conf.Token, conf.ChatId); err != nil {
			return err
		}
	}
	return nil
}

func run(cmd *cobra.Command, domainNames []string) (err error) {
	domains := make(domain.Domains, 0, len(domainNames))
	for _, domainName := range domainNames {
		domains = append(domains, domain.Domain{Name: domainName})
	}

	domains.Tick()

	if conf.RunEvery != "" {
		log.Info("running as cron")

		c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
		_, err := c.AddFunc("@every "+conf.RunEvery, domains.Tick)
		if err != nil {
			log.WithError(err).Error("failed to register job")
			return err
		}
		c.Run()
	}

	return nil
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
