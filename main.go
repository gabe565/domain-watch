package main

import (
	"fmt"
	"github.com/araddon/dateparse"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/likexian/whois-go"
	whoisparser "github.com/likexian/whois-parser-go"
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

var bot *tgbotapi.BotAPI

var config Config

func init() {
	cmd.Flags().StringVar(&config.RunEvery, "every", "", "enable cron mode and configure update interval")
	cmd.Flags().DurationVar(&config.Sleep, "sleep", 3*time.Second, "sleep time between queries to avoid rate limits")
	cmd.Flags().StringVar(&config.Token, "telegram-token", "", "Telegram token")
	cmd.Flags().Int64Var(&config.ChatId, "telegram-chat", 0, "Telegram chat ID")
}

func getParsedWhois(domain string) (result whoisparser.WhoisInfo, err error) {
	rawWhois, err := whois.Whois(domain)
	if err != nil {
		return
	}
	result, err = whoisparser.Parse(rawWhois)
	return
}

func getDomainExpiration(parsedWhois whoisparser.WhoisInfo) (time.Time, error) {
	return dateparse.ParseStrict(parsedWhois.Domain.ExpirationDate)
}

func loopOverDomains(domains []string) {
	for key, domain := range domains {
		if key != 0 {
			time.Sleep(config.Sleep)
		}

		parsedWhois, err := getParsedWhois(domain)
		if err != nil {
			log.Printf("error fetching %s: %v\n", domain, err)
			continue
		}

		l := log.WithField("domain", domain)

		if parsedWhois.Domain.ExpirationDate != "" {
			var date time.Time
			date, err = getDomainExpiration(parsedWhois)
			if err != nil {
				l.WithError(err).Warn("failed to parse expiration date")
			} else {
				left := date.Sub(time.Now()).Truncate(24 * time.Hour)
				l.WithFields(log.Fields{
					"expires":   date,
					"days_left": left.Hours() / 24.0,
				}).Info("fetched whois")
			}
		} else {
			l.Info("domain does not have an expiration date")
		}

		if cachedWhois, ok := config.WhoisCache[domain]; ok {
			notify(parsedWhois, cachedWhois)
		}

		if config.WhoisCache != nil {
			config.WhoisCache[domain] = parsedWhois
		}
	}
}

func preRun(cmd *cobra.Command, domains []string) (err error) {
	if config.Token != "" && config.ChatId != 0 {
		var err error
		bot, err = telegramLogin(config.Token, config.ChatId)
		if err != nil {
			return err
		}
	}
	return nil
}

func run(cmd *cobra.Command, domains []string) (err error) {
	if config.RunEvery != "" {
		config.WhoisCache = make(map[string]whoisparser.WhoisInfo)

		log.Info("initial run to fill cache")
		loopOverDomains(domains)

		log.Info("running as cron")

		c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
		_, err := c.AddFunc(fmt.Sprintf("@every %s", config.RunEvery), func() {
			loopOverDomains(domains)
		})
		if err != nil {
			log.WithError(err).Error("failed to register job")
			return err
		}
		c.Run()
	} else {
		loopOverDomains(domains)
	}

	return nil
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
