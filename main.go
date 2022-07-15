package main

import (
	"flag"
	"fmt"
	"github.com/araddon/dateparse"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/likexian/whois-go"
	whoisparser "github.com/likexian/whois-parser-go"
	"github.com/robfig/cron/v3"
	"log"
	"math"
	"os"
	"time"
)

var (
	bot        *tgbotapi.BotAPI
	chatId     int64
	quiet      bool
	whoisCache map[string]whoisparser.WhoisInfo
	runEvery   string
	token      string
	sleep      time.Duration
)

func usage() {
	fmt.Printf("Usage: %s [OPTIONS] domain ...\n", os.Args[0])
	flag.PrintDefaults()
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

func calcDaysUntil(date time.Time) int64 {
	return int64(math.Floor(time.Until(date).Hours() / 24))
}

func loopOverDomains(domains []string) {
	for key, domain := range domains {
		if key != 0 {
			time.Sleep(sleep)
		}

		parsedWhois, err := getParsedWhois(domain)
		if err != nil {
			log.Printf("error fetching %s: %v\n", domain, err)
			continue
		}

		if parsedWhois.Domain.ExpirationDate != "" {
			var date time.Time
			date, err = getDomainExpiration(parsedWhois)
			if err == nil {
				daysUntil := calcDaysUntil(date)
				if !quiet {
					log.Printf("%s %s %d", domain, date.Format("2006-01-02"), daysUntil)
				}
			}
		} else {
			if !quiet {
				log.Printf("%s does not have an expiration date", domain)
			}
		}

		if cachedWhois, ok := whoisCache[domain]; ok {
			notify(parsedWhois, cachedWhois)
		}

		if whoisCache != nil {
			whoisCache[domain] = parsedWhois
		}
	}
}

func main() {
	flag.StringVar(&runEvery, "every", "", "Will enable cron mode and configure update interval.")
	flag.BoolVar(&quiet, "q", false, "Run in quiet mode")
	flag.StringVar(&token, "telegram-token", "", "Telegram token to user on bot login.")
	flag.Int64Var(&chatId, "telegram-chat", 0, "Telegram chat/user ID.")
	flag.DurationVar(&sleep, "sleep", 3*time.Second, "Time to sleep between queries to avoid rate limits.")
	flag.Usage = usage
	flag.Parse()

	// Setup Telegram bot
	if token != "" && chatId != 0 {
		var err error
		bot, err = telegramLogin(token, chatId)
		if err != nil {
			log.Panicln(err)
		}
	}

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	// Pull domains from command line args
	domains := flag.Args()

	if runEvery != "" {
		whoisCache = make(map[string]whoisparser.WhoisInfo)

		log.Println("Initial run to fill cache")
		loopOverDomains(domains)

		log.Println("Running as cron")

		c := cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
		_, err := c.AddFunc(fmt.Sprintf("@every %s", runEvery), func() {
			loopOverDomains(domains)
		})
		if err != nil {
			log.Println(err)
			return
		}
		c.Run()
	} else {
		loopOverDomains(domains)
	}
}
