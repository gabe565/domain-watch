package domain

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/gabe565/domain-watch/internal/telegram"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser-go"
	"github.com/r3labs/diff"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type Domain struct {
	Name               string
	CurrWhois          whoisparser.WhoisInfo
	PrevWhois          *whoisparser.WhoisInfo
	TimeLeft           time.Duration
	TriggeredThreshold int
}

func (d Domain) Whois() (whoisparser.WhoisInfo, error) {
	raw, err := whois.Whois(d.Name)
	if err != nil {
		return whoisparser.WhoisInfo{}, err
	}

	return whoisparser.Parse(raw)
}

func (d Domain) Log() *log.Entry {
	return log.WithField("domain", d.Name)
}

func (d *Domain) Run() (err error) {
	d.CurrWhois, err = d.Whois()
	if err != nil {
		return fmt.Errorf("failed to fetch whois: %w", err)
	}
	defer func() {
		d.PrevWhois = &d.CurrWhois
	}()

	l := d.Log()

	if d.CurrWhois.Domain.ExpirationDate != "" {
		var date time.Time
		date, err = dateparse.ParseStrict(d.CurrWhois.Domain.ExpirationDate)
		if err != nil {
			d.TimeLeft = 0
			l.WithError(err).Warn("failed to parse expiration date")
		} else {
			d.TimeLeft = time.Until(date).Truncate(24 * time.Hour)
			l.WithFields(log.Fields{
				"expires":   date,
				"days_left": d.TimeLeft.Hours() / 24.0,
			}).Info("fetched whois")
		}
	} else {
		l.Info("domain does not have an expiration date")
	}

	if err := d.CheckNotifications(); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (d *Domain) NotifyThreshold() error {
	if d.TimeLeft != 0 {
		daysLeft := int(d.TimeLeft.Hours() / 24)
		for _, threshold := range viper.GetIntSlice("threshold") {
			if d.TriggeredThreshold != threshold && daysLeft <= threshold {
				msg := telegram.NewThresholdMessage(d.Name, daysLeft)
				if err := telegram.Send(msg); err != nil {
					return err
				}
				d.TriggeredThreshold = threshold
				break
			}
		}
	}
	return nil
}

func (d *Domain) NotifyStatusChange() error {
	if d.PrevWhois != nil {
		changes, err := diff.Diff(d.PrevWhois.Domain.Status, d.CurrWhois.Domain.Status)
		if err != nil {
			return err
		}

		if len(changes) > 0 {
			msg := telegram.NewStatusChangedMessage(d.Name, changes)
			if err := telegram.Send(msg); err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *Domain) CheckNotifications() error {
	if !telegram.LoggedIn() {
		return nil
	}
	if err := d.NotifyThreshold(); err != nil {
		return err
	}
	if err := d.NotifyStatusChange(); err != nil {
		return err
	}
	return nil
}
