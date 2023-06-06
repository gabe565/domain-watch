package integration

import (
	"fmt"

	"github.com/gabe565/domain-watch/internal/util"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Telegram struct {
	ChatId int64
	Bot    *tgbotapi.BotAPI
}

func (t *Telegram) Flags(cmd *cobra.Command) error {
	cmd.Flags().String("telegram-token", "", "Telegram token")
	if err := viper.BindPFlag("telegram.token", cmd.Flags().Lookup("telegram-token")); err != nil {
		panic(err)
	}
	if err := cmd.RegisterFlagCompletionFunc("telegram-token", util.NoFileComp); err != nil {
		panic(err)
	}

	cmd.Flags().Int64("telegram-chat", 0, "Telegram chat ID")
	if err := viper.BindPFlag("telegram.chat", cmd.Flags().Lookup("telegram-chat")); err != nil {
		panic(err)
	}
	if err := cmd.RegisterFlagCompletionFunc("telegram-chat", util.NoFileComp); err != nil {
		panic(err)
	}

	cmd.MarkFlagsRequiredTogether("telegram-token", "telegram-chat")

	return nil
}

func (t *Telegram) Setup() error {
	token := viper.GetString("telegram.token")
	if token == "" {
		return fmt.Errorf("telegram %w: token", util.ErrNotConfigured)
	}

	t.ChatId = viper.GetInt64("telegram.chat")
	if t.ChatId == 0 {
		return fmt.Errorf("telegram %w: chat ID", util.ErrNotConfigured)
	}

	return t.Login(token)
}

func (t *Telegram) Login(token string) (err error) {
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
