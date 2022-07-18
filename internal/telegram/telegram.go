package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	whoisparser "github.com/likexian/whois-parser-go"
	"github.com/r3labs/diff"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	bot *tgbotapi.BotAPI
)

func Login(token string) (err error) {
	if token != "" {
		bot, err = tgbotapi.NewBotAPI(token)
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"username": bot.Self.UserName,
		}).Info("auth success")
	}
	return nil
}

func CreateMessage(domain string, changes []diff.Change) (msg tgbotapi.MessageConfig) {
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
		viper.GetInt64("telegram.chat"),
		fmt.Sprintf(
			"The statuses on %s have changed. Here are the changes:\n```%s%s```",
			domain,
			removed,
			added,
		),
	)
	msg.ParseMode = tgbotapi.ModeMarkdown
	return
}

func Notify(parsedWhois whoisparser.WhoisInfo, cachedWhois whoisparser.WhoisInfo) bool {
	if bot == nil {
		return false
	}
	changes, err := diff.Diff(cachedWhois.Domain.Status, parsedWhois.Domain.Status)
	if err != nil {
		return false
	}
	if len(changes) > 0 {
		msg := CreateMessage(parsedWhois.Domain.Domain, changes)
		_, _ = bot.Send(msg)
	}
	return true
}
