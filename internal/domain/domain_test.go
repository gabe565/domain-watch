package domain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gabe565.com/domain-watch/internal/config"
	"gabe565.com/domain-watch/internal/integration"
	"gabe565.com/domain-watch/internal/integration/telegram"
	"github.com/go-telegram/bot"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomain_Whois(t *testing.T) {
	type fields struct {
		Name               string
		CurrWhois          whoisparser.WhoisInfo
		PrevWhois          *whoisparser.WhoisInfo
		TimeLeft           time.Duration
		TriggeredThreshold int
	}
	tests := []struct {
		name       string
		fields     fields
		wantDomain string
		wantErr    require.ErrorAssertionFunc
	}{
		{"example.com", fields{Name: "example.com"}, "example.com", require.NoError},
		{"a", fields{Name: "a"}, "", require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Domain{
				Name:               tt.fields.Name,
				CurrWhois:          tt.fields.CurrWhois,
				PrevWhois:          tt.fields.PrevWhois,
				TimeLeft:           tt.fields.TimeLeft,
				TriggeredThreshold: tt.fields.TriggeredThreshold,
			}
			got, err := d.Whois()
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.wantDomain, got.Domain.Domain)
			}
		})
	}
}

func TestDomain_NotifyThreshold(t *testing.T) {
	type fields struct {
		Name               string
		CurrWhois          whoisparser.WhoisInfo
		PrevWhois          *whoisparser.WhoisInfo
		TimeLeft           time.Duration
		TriggeredThreshold int
	}
	tests := []struct {
		name       string
		fields     fields
		wantNotify bool
		wantErr    require.ErrorAssertionFunc
	}{
		{"example.com 30d", fields{Name: "example.com", TimeLeft: 30 * 24 * time.Hour}, false, require.NoError},
		{"example.com 7d", fields{Name: "example.com", TimeLeft: 7 * 24 * time.Hour}, true, require.NoError},
		{"example.com 1d", fields{Name: "example.com", TimeLeft: 24 * time.Hour}, true, require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNotify := false
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/bot123/sendMessage", r.URL.Path)
				gotNotify = true
				resp := telegram.APIResponse{
					OK:     true,
					Result: json.RawMessage("{}"),
				}
				assert.NoError(t, json.NewEncoder(w).Encode(&resp))
			}))
			t.Cleanup(server.Close)
			tg := telegram.NewTestClient(t, bot.WithServerURL(server.URL))

			d := &Domain{
				conf: &config.Config{Threshold: []int{1, 7}},

				Name:               tt.fields.Name,
				CurrWhois:          tt.fields.CurrWhois,
				PrevWhois:          tt.fields.PrevWhois,
				TimeLeft:           tt.fields.TimeLeft,
				TriggeredThreshold: tt.fields.TriggeredThreshold,
			}
			tt.wantErr(t, d.NotifyThreshold(t.Context(), integration.Integrations{tg}))
			assert.Equal(t, tt.wantNotify, gotNotify)
		})
	}
}

func TestDomain_NotifyStatusChange(t *testing.T) {
	type fields struct {
		Name               string
		CurrWhois          whoisparser.WhoisInfo
		PrevWhois          *whoisparser.WhoisInfo
		TimeLeft           time.Duration
		TriggeredThreshold int
	}
	tests := []struct {
		name       string
		fields     fields
		wantNotify bool
		wantErr    require.ErrorAssertionFunc
	}{
		{"example.com no change", fields{Name: "example.com"}, false, require.NoError},
		{"example.com created status", fields{
			Name: "example.com",
			PrevWhois: &whoisparser.WhoisInfo{
				Domain: &whoisparser.Domain{Status: []string{}},
			},
			CurrWhois: whoisparser.WhoisInfo{
				Domain: &whoisparser.Domain{Status: []string{"a"}},
			},
		}, true, require.NoError},
		{"example.com removed status", fields{
			Name: "example.com",
			PrevWhois: &whoisparser.WhoisInfo{
				Domain: &whoisparser.Domain{Status: []string{"a"}},
			},
			CurrWhois: whoisparser.WhoisInfo{
				Domain: &whoisparser.Domain{Status: []string{}},
			},
		}, true, require.NoError},
		{"example.com changed status", fields{
			Name: "example.com",
			PrevWhois: &whoisparser.WhoisInfo{
				Domain: &whoisparser.Domain{Status: []string{"a"}},
			},
			CurrWhois: whoisparser.WhoisInfo{
				Domain: &whoisparser.Domain{Status: []string{"b"}},
			},
		}, true, require.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNotify := false
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/bot123/sendMessage", r.URL.Path)
				gotNotify = true
				resp := telegram.APIResponse{
					OK:     true,
					Result: json.RawMessage("{}"),
				}
				assert.NoError(t, json.NewEncoder(w).Encode(resp))
			}))
			t.Cleanup(server.Close)
			tg := telegram.NewTestClient(t, bot.WithServerURL(server.URL))

			d := &Domain{
				Name:               tt.fields.Name,
				CurrWhois:          tt.fields.CurrWhois,
				PrevWhois:          tt.fields.PrevWhois,
				TimeLeft:           tt.fields.TimeLeft,
				TriggeredThreshold: tt.fields.TriggeredThreshold,
			}
			tt.wantErr(t, d.NotifyStatusChange(t.Context(), integration.Integrations{tg}))
			assert.Equal(t, tt.wantNotify, gotNotify)
		})
	}
}
