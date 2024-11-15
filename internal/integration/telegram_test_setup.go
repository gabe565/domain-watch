package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TelegramTestSetup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/bot/getMe" {
			var buf bytes.Buffer
			u := tgbotapi.User{
				IsBot:         true,
				FirstName:     "Bot",
				UserName:      "Bot",
				CanJoinGroups: true,
			}
			assert.NoError(t, json.NewEncoder(&buf).Encode(u))

			resp := tgbotapi.APIResponse{
				Ok:     true,
				Result: json.RawMessage(buf.Bytes()),
			}
			assert.NoError(t, json.NewEncoder(w).Encode(resp))
		}
	}))
	t.Cleanup(server.Close)

	telegram := Get("telegram").(*Telegram)
	var err error
	telegram.Bot, err = tgbotapi.NewBotAPIWithAPIEndpoint("", server.URL+"/bot%s/%s")
	require.NoError(t, err)
}
