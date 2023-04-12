package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/r3labs/diff/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Bot *tgbotapi.BotAPI

func Login(token string) (err error) {
	if token != "" {
		Bot, err = tgbotapi.NewBotAPI(token)
		if err != nil {
			return err
		}
		log.WithFields(log.Fields{
			"username": Bot.Self.UserName,
		}).Info("auth success")
	}
	return nil
}

func LoggedIn() bool {
	return Bot != nil
}

func Send(msg tgbotapi.MessageConfig) error {
	_, err := Bot.Send(msg)
	return err
}

func NewStatusChangedMessage(domain string, changes []diff.Change) (msg tgbotapi.MessageConfig) {
	var added, removed string
	for _, change := range changes {
		switch change.Type {
		case diff.UPDATE:
			removed += fmt.Sprintf("\n - %s", change.From)
			added += fmt.Sprintf("\n + %s", change.To)
		case diff.CREATE:
			added += fmt.Sprintf("\n + %s", change.To)
		case diff.DELETE:
			removed += fmt.Sprintf("\n - %s", change.From)
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

func NewThresholdMessage(domain string, timeLeft int) (msg tgbotapi.MessageConfig) {
	return tgbotapi.NewMessage(
		viper.GetInt64("telegram.chat"),
		fmt.Sprintf("%s will expire in %d days.", domain, timeLeft),
	)
}
