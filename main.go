package main

import (
	"flag"
	"fmt"
	"github.com/araddon/dateparse"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/likexian/whois-go"
	"github.com/likexian/whois-parser-go"
	"github.com/r3labs/diff"
	"github.com/robfig/cron/v3"
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
	whoisCache map[string]whoisparser.WhoisInfo
	runAsCron bool
	cronSpec string
	days int64
	token string
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

func daysUntil(date time.Time) int64 {
	return int64(math.Floor(time.Until(date).Hours() / 24))
}

func createMessage(domain string, changes []diff.Change) (msg tgbotapi.MessageConfig) {
	removed := ""
	added := ""
	for _, change := range changes {
		switch change.Type {
		case "update":
			removed += fmt.Sprintf("\n - %s", change.From)
			added += fmt.Sprintf("\n + %s", change.To)
			break
		case "create":
			added += fmt.Sprintf("\n + %s", change.To)
			break
		case "delete":
			removed += fmt.Sprintf("\n - %s", change.From)
			break
		}
	}
	msg = tgbotapi.NewMessage(
		chatId,
		fmt.Sprintf("The statuses on %s have changed. Here are the changes:\n```%s%s```", domain, removed, added),
	)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return
}

func notify(parsedWhois whoisparser.WhoisInfo, cachedWhois whoisparser.WhoisInfo) bool {
	if bot == nil {
		return false
	}
	changes, err := diff.Diff(cachedWhois.Domain.Status, parsedWhois.Domain.Status)
	if err != nil {
		return false
	}
	if len(changes) > 0 {
		msg := createMessage(parsedWhois.Domain.Domain, changes)
		_, _ = bot.Send(msg)
	}
	return true
}

func loopOverDomains(domains []string) {
	wg := &sync.WaitGroup{}
	wg.Add(len(domains))

	for _, domain := range domains {
		go func (domain string) {
			defer wg.Done()

			parsedWhois, err := getParsedWhois(domain)
			if err != nil {
				log.Println(err)
				return
			}
			date, err := getDomainExpiration(parsedWhois)
			if err != nil {
				log.Println(err)
				return
			}
			daysUntil := daysUntil(date)
			if !quiet {
				log.Printf("%s %s %d", domain, date.Format("2006-01-02"), daysUntil)
			}

			if cachedWhois, ok := whoisCache[domain]; ok {
				notify(parsedWhois, cachedWhois)
			}

			if runAsCron {
				whoisCache[domain] = parsedWhois
			}
		} (domain)
	}

	wg.Wait()
}

func main() {
	flag.Int64Var(&days, "days", 31, "Number of days within a notification is triggered.")
	flag.BoolVar(&runAsCron, "cron", false, "Whether to run daily as a cron or as a one-off.")
	flag.StringVar(&cronSpec, "cron-spec", "0 * * * *", "When to run update checks in cron format.")
	flag.BoolVar(&quiet, "q", false, "Run in quiet mode")
	flag.StringVar(&token, "telegram-token", "", "Telegram token to user on bot login.")
	flag.Int64Var(&chatId, "telegram-chat", 0, "Telegram chat/user ID.")
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

		whoisCache = make(map[string]whoisparser.WhoisInfo)

		c := cron.New()
		_, err := c.AddFunc(cronSpec, func() {
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
