package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	whoisparser "github.com/likexian/whois-parser-go"
	"github.com/r3labs/diff"
	log "github.com/sirupsen/logrus"
)

func telegramLogin(token string, chatId int64) (bot *tgbotapi.BotAPI, err error) {
	if token != "" && chatId != 0 {
		bot, err = tgbotapi.NewBotAPI(token)
		if err != nil {
			return
		}
		log.WithFields(log.Fields{
			"username": bot.Self.UserName,
			"chat":     chatId,
		}).Info("auth success", bot.Self.UserName)
	}
	return
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
		config.ChatId,
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
