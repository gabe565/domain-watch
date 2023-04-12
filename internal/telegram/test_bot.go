package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func InitTestBot() (err error) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/bot/getMe" {
			var buf bytes.Buffer
			u := tgbotapi.User{
				IsBot:         true,
				FirstName:     "Bot",
				UserName:      "Bot",
				CanJoinGroups: true,
			}
			if err := json.NewEncoder(&buf).Encode(u); err != nil {
				panic(err)
			}

			resp := tgbotapi.APIResponse{
				Ok:     true,
				Result: json.RawMessage(buf.Bytes()),
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				panic(err)
			}
		}
	}))
	defer server.Close()

	Bot, err = tgbotapi.NewBotAPIWithAPIEndpoint("", server.URL+"/bot%s/%s")
	if err != nil {
		return err
	}

	return nil
}
