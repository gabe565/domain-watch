package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/r3labs/diff"
	"reflect"
	"testing"
)

func TestNewStatusChangedMessage(t *testing.T) {
	type args struct {
		domain  string
		changes []diff.Change
	}
	tests := []struct {
		name    string
		args    args
		wantMsg tgbotapi.MessageConfig
	}{
		{
			"create",
			args{
				"example.com",
				[]diff.Change{
					{Type: diff.CREATE, To: "a"},
				},
			},
			tgbotapi.MessageConfig{
				ParseMode: tgbotapi.ModeMarkdown,
				Text:      "The statuses on example.com have changed. Here are the changes:\n```\n + a```",
			},
		},
		{
			"update",
			args{
				"example.com",
				[]diff.Change{
					{Type: diff.UPDATE, From: "a", To: "b"},
				},
			},
			tgbotapi.MessageConfig{
				ParseMode: tgbotapi.ModeMarkdown,
				Text:      "The statuses on example.com have changed. Here are the changes:\n```\n - a\n + b```",
			},
		},
		{
			"delete",
			args{
				"example.com",
				[]diff.Change{
					{Type: diff.DELETE, From: "a"},
				},
			},
			tgbotapi.MessageConfig{
				ParseMode: tgbotapi.ModeMarkdown,
				Text:      "The statuses on example.com have changed. Here are the changes:\n```\n - a```",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMsg := NewStatusChangedMessage(tt.args.domain, tt.args.changes); !reflect.DeepEqual(gotMsg, tt.wantMsg) {
				t.Errorf("NewStatusChangedMessage() = %v, want %v", gotMsg, tt.wantMsg)
			}
		})
	}
}

func TestNewThresholdMessage(t *testing.T) {
	type args struct {
		domain   string
		timeLeft int
	}
	tests := []struct {
		name    string
		args    args
		wantMsg tgbotapi.MessageConfig
	}{
		{
			"example.com 7d",
			args{"example.com", 7},
			tgbotapi.MessageConfig{
				Text: "example.com will expire in 7 days.",
			},
		},
		{
			"example.com 1d",
			args{"example.com", 1},
			tgbotapi.MessageConfig{
				Text: "example.com will expire in 1 days.",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMsg := NewThresholdMessage(tt.args.domain, tt.args.timeLeft); !reflect.DeepEqual(gotMsg, tt.wantMsg) {
				t.Errorf("NewThresholdMessage() = %v, want %v", gotMsg, tt.wantMsg)
			}
		})
	}
}
