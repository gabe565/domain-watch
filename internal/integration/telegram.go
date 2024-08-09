package integration

import (
	"fmt"

	"github.com/gabe565/domain-watch/internal/config"
	"github.com/gabe565/domain-watch/internal/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

type Telegram struct {
	ChatId int64
	Bot    *tgbotapi.BotAPI
}

func (t *Telegram) Setup(conf *config.Config) error {
	if t.ChatId = conf.TelegramChat; t.ChatId == 0 {
		return fmt.Errorf("telegram %w: chat ID", util.ErrNotConfigured)
	}

	return t.Login(conf.TelegramToken)
}

func (t *Telegram) Login(token string) (err error) {
	if token == "" {
		return fmt.Errorf("telegram %w: token", util.ErrNotConfigured)
	}

	t.Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"username": t.Bot.Self.UserName,
	}).Info("connected to Telegram")
	return nil
}

func (t *Telegram) Send(message string) error {
	if t.Bot == nil {
		return nil
	}

	payload := tgbotapi.NewMessage(t.ChatId, message)
	payload.ParseMode = tgbotapi.ModeMarkdown

	_, err := t.Bot.Send(payload)
	return err
}
