package domain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"gabe565.com/domain-watch/internal/config"
	"gabe565.com/domain-watch/internal/integration"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	whoisparser "github.com/likexian/whois-parser"
	log "github.com/sirupsen/logrus"
)

func TestDomain_Log(t *testing.T) {
	type fields struct {
		Name               string
		CurrWhois          whoisparser.WhoisInfo
		PrevWhois          *whoisparser.WhoisInfo
		TimeLeft           time.Duration
		TriggeredThreshold int
	}
	tests := []struct {
		name   string
		fields fields
		want   *log.Entry
	}{
		{"example.com", fields{Name: "example.com"}, log.WithField("domain", "example.com")},
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
			if got := d.Log(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Log() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		wantErr    bool
	}{
		{"example.com", fields{Name: "example.com"}, "example.com", false},
		{"a", fields{Name: "a"}, "", true},
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
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Whois() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if got.Domain.Domain != tt.wantDomain {
				t.Errorf("Whois() got domain = %v, want domain %v", got.Domain.Domain, tt.wantDomain)
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
		wantErr    bool
	}{
		{"example.com 30d", fields{Name: "example.com", TimeLeft: 30 * 24 * time.Hour}, false, false},
		{"example.com 7d", fields{Name: "example.com", TimeLeft: 7 * 24 * time.Hour}, true, false},
		{"example.com 1d", fields{Name: "example.com", TimeLeft: 24 * time.Hour}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := integration.TelegramTestSetup(); err != nil {
				t.Error(err)
				return
			}
			gotNotify := false
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot/sendMessage" {
					t.Errorf("Expected to request /bot/sendMessage, got: %s", r.URL.Path)
					return
				}
				gotNotify = true
				resp := tgbotapi.APIResponse{
					Ok:     true,
					Result: json.RawMessage("{}"),
				}
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			integration.Get("telegram").(*integration.Telegram).Bot.SetAPIEndpoint(server.URL + "/bot%s/%s")

			d := &Domain{
				conf: &config.Config{Threshold: []int{1, 7}},

				Name:               tt.fields.Name,
				CurrWhois:          tt.fields.CurrWhois,
				PrevWhois:          tt.fields.PrevWhois,
				TimeLeft:           tt.fields.TimeLeft,
				TriggeredThreshold: tt.fields.TriggeredThreshold,
			}
			if err := d.NotifyThreshold(); (err != nil) != tt.wantErr {
				t.Errorf("NotifyThreshold() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotNotify != tt.wantNotify {
				t.Errorf("NotifyThreshold() got notification = %t, want notification %t", gotNotify, tt.wantNotify)
			}
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
		wantErr    bool
	}{
		{"example.com no change", fields{Name: "example.com"}, false, false},
		{"example.com created status", fields{
			Name: "example.com",
			PrevWhois: &whoisparser.WhoisInfo{
				Domain: &whoisparser.Domain{Status: []string{}},
			},
			CurrWhois: whoisparser.WhoisInfo{
				Domain: &whoisparser.Domain{Status: []string{"a"}},
			},
		}, true, false},
		{"example.com removed status", fields{
			Name: "example.com",
			PrevWhois: &whoisparser.WhoisInfo{
				Domain: &whoisparser.Domain{Status: []string{"a"}},
			},
			CurrWhois: whoisparser.WhoisInfo{
				Domain: &whoisparser.Domain{Status: []string{}},
			},
		}, true, false},
		{"example.com changed status", fields{
			Name: "example.com",
			PrevWhois: &whoisparser.WhoisInfo{
				Domain: &whoisparser.Domain{Status: []string{"a"}},
			},
			CurrWhois: whoisparser.WhoisInfo{
				Domain: &whoisparser.Domain{Status: []string{"b"}},
			},
		}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := integration.TelegramTestSetup(); err != nil {
				t.Error(err)
				return
			}
			gotNotify := false
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/bot/sendMessage" {
					t.Errorf("Expected to request /bot/sendMessage, got: %s", r.URL.Path)
					return
				}
				gotNotify = true
				resp := tgbotapi.APIResponse{
					Ok:     true,
					Result: json.RawMessage("{}"),
				}
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					t.Error(err)
				}
			}))
			defer server.Close()
			integration.Get("telegram").(*integration.Telegram).Bot.SetAPIEndpoint(server.URL + "/bot%s/%s")

			d := &Domain{
				Name:               tt.fields.Name,
				CurrWhois:          tt.fields.CurrWhois,
				PrevWhois:          tt.fields.PrevWhois,
				TimeLeft:           tt.fields.TimeLeft,
				TriggeredThreshold: tt.fields.TriggeredThreshold,
			}
			if err := d.NotifyStatusChange(); (err != nil) != tt.wantErr {
				t.Errorf("NotifyStatusChange() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotNotify != tt.wantNotify {
				t.Errorf("NotifyStatusChange() got notification = %t, want notification %t", gotNotify, tt.wantNotify)
			}
		})
	}
}
