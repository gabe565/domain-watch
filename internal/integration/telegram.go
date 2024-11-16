package integration

import (
	"context"
	"fmt"
	"log/slog"

	"gabe565.com/domain-watch/internal/config"
	"gabe565.com/domain-watch/internal/util"
	"github.com/go-telegram/bot"
)

type Telegram struct {
	ChatID int64
	Bot    *bot.Bot
}

func (t *Telegram) Setup(ctx context.Context, conf *config.Config) error {
	if t.ChatID = conf.TelegramChat; t.ChatID == 0 {
		return fmt.Errorf("telegram %w: chat ID", util.ErrNotConfigured)
	}

	return t.Login(ctx, conf.TelegramToken)
}

func (t *Telegram) Login(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("telegram %w: token", util.ErrNotConfigured)
	}

	var err error
	t.Bot, err = bot.New(token, bot.WithSkipGetMe())
	if err != nil {
		return err
	}

	user, err := t.Bot.GetMe(ctx)
	if err != nil {
		return err
	}

	slog.Info("Connected to Telegram", "username", user.Username)
	return nil
}

func (t *Telegram) Send(ctx context.Context, message string) error {
	if t.Bot == nil {
		return nil
	}

	_, err := t.Bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    t.ChatID,
		Text:      message,
		ParseMode: "markdown",
	})
	return err
}
