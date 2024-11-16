package integration

import (
	"context"
	"fmt"
	"log/slog"

	"gabe565.com/domain-watch/internal/config"
	"gabe565.com/domain-watch/internal/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	ChatID int64
	Bot    *tgbotapi.BotAPI
}

func (t *Telegram) Setup(_ context.Context, conf *config.Config) error {
	if t.ChatID = conf.TelegramChat; t.ChatID == 0 {
		return fmt.Errorf("telegram %w: chat ID", util.ErrNotConfigured)
	}

	return t.Login(conf.TelegramToken)
}

func (t *Telegram) Login(token string) error {
	if token == "" {
		return fmt.Errorf("telegram %w: token", util.ErrNotConfigured)
	}

	var err error
	t.Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	slog.Info("Connected to Telegram", "username", t.Bot.Self.UserName)
	return nil
}

func (t *Telegram) Send(_ context.Context, message string) error {
	if t.Bot == nil {
		return nil
	}

	payload := tgbotapi.NewMessage(t.ChatID, message)
	payload.ParseMode = tgbotapi.ModeMarkdown

	_, err := t.Bot.Send(payload)
	return err
}
