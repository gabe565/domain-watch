package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type APIResponse struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result,omitempty"`
	Description string          `json:"description,omitempty"`
	ErrorCode   int             `json:"error_code,omitempty"`
	Parameters  struct {
		RetryAfter      int `json:"retry_after,omitempty"`
		MigrateToChatID int `json:"migrate_to_chat_id,omitempty"`
	} `json:"parameters,omitempty"`
}

func NewTestClient(t *testing.T, opts ...bot.Option) *Telegram {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/bot123/getMe" {
			var buf bytes.Buffer
			u := models.User{
				IsBot:         true,
				FirstName:     "Bot",
				Username:      "Bot",
				CanJoinGroups: true,
			}
			assert.NoError(t, json.NewEncoder(&buf).Encode(u))

			resp := APIResponse{
				OK:     true,
				Result: json.RawMessage(buf.Bytes()),
			}
			assert.NoError(t, json.NewEncoder(w).Encode(resp))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	telegram := &Telegram{}
	var err error
	telegram.Bot, err = bot.New("123", bot.WithServerURL(server.URL))
	require.NoError(t, err)

	for _, opt := range opts {
		opt(telegram.Bot)
	}
	return telegram
}
