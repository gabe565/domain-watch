package main

import (
	"flag"
	"fmt"
	"github.com/araddon/dateparse"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/likexian/whois-go"
	"github.com/likexian/whois-parser-go"
	cron "github.com/robfig/cron/v3"
	"log"
	"math"
	"os"
	"sync"
	"time"
)

var (
	bot *tgbotapi.BotAPI
	chatId int64
	quiet bool
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] domain ...\n", os.Args[0])
	flag.PrintDefaults()
}

func getDomainExpiration(domain string) (result time.Time, err error) {
	rawWhois, err := whois.Whois(domain)
	if err != nil {
		return time.Now(), err
	}
	parsedWhois, err := whoisparser.Parse(rawWhois)
	if err != nil {
		return time.Now(), err
	}
	result, err = dateparse.ParseStrict(parsedWhois.Domain.ExpirationDate )
	if err != nil {
		return time.Now(), err
	}
	return result, nil
}


func daysUntil(date time.Time) (days int64) {
	return int64(math.Floor(time.Until(date).Hours() / 24))
}

func loopOverDomains(domains []string, daysUntilNotify int64) {
	wg := &sync.WaitGroup{}
	wg.Add(len(domains))

	for _, domain := range domains {
		go func (domain string) {
			defer wg.Done()

			date, err := getDomainExpiration(domain)
			if err != nil {
				log.Println(err)
				return
			}
			daysUntil := daysUntil(date)
			if !quiet {
				log.Printf("%s\t| %s\t| %d Days\n", domain, date.Format("2006-01-02"), daysUntil)
			}
			if bot != nil && daysUntil <= daysUntilNotify {
				msg := tgbotapi.NewMessage(
					chatId,
					fmt.Sprintf("The domain registration at %s will expire in %d days.", domain, daysUntil),
				)
				_, _ = bot.Send(msg)
			}
		} (domain)
	}

	wg.Wait()
}

func main() {
	var days int64
	flag.Int64Var(&days, "days", 31, "Number of days within a notification is triggered.")

	var runAsCron bool
	flag.BoolVar(&runAsCron, "cron", false, "Whether to run daily as a cron or as a one-off.")
	var cronSpec string
	flag.StringVar(&cronSpec, "cron-spec", "0 9 * * *", "When to run update checks in cron format.")

	flag.BoolVar(&quiet, "q", false, "Run in quiet mode")

	var token string
	flag.StringVar(&token, "telegram-token", "", "Telegram token to user on bot login.")
	flag.Int64Var(&chatId, "telegram-id", 0, "Telegram chat/user ID.")

	flag.Usage = usage
	flag.Parse()

	// Setup Telegram bot
	if token != "" && chatId != 0 {
		var err error
		bot, err = tgbotapi.NewBotAPI(token)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Authorized on account %s", bot.Self.UserName)
		log.Printf("Sending notifications to chat ID %d", chatId)
	}

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	// Pull domains from command line args
	domains := flag.Args()

	if runAsCron {
		log.Printf("Running as cron")
		c := cron.New()
		c.AddFunc(cronSpec, func() {
			loopOverDomains(domains, days)
		})
		c.Run()
	} else {
		loopOverDomains(domains, days)
	}
}
