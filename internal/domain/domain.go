package domain

import (
	"fmt"
	"time"

	"gabe565.com/domain-watch/internal/config"
	"gabe565.com/domain-watch/internal/integration"
	"gabe565.com/domain-watch/internal/message"
	"github.com/araddon/dateparse"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/r3labs/diff/v3"
	log "github.com/sirupsen/logrus"
)

func New(conf *config.Config, name string) Domain {
	return Domain{conf: conf, Name: name}
}

type Domain struct {
	conf *config.Config

	Name               string
	CurrWhois          whoisparser.WhoisInfo
	PrevWhois          *whoisparser.WhoisInfo
	ExpiresAt          time.Time
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
		d.ExpiresAt, err = dateparse.ParseStrict(d.CurrWhois.Domain.ExpirationDate)
		if err != nil {
			d.TimeLeft = 0
			l.WithError(err).Warn("failed to parse expiration date")
		} else {
			d.TimeLeft = time.Until(d.ExpiresAt).Truncate(24 * time.Hour)
			l.WithFields(log.Fields{
				"expires":   d.ExpiresAt,
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
		for _, threshold := range d.conf.Threshold {
			if d.TriggeredThreshold != threshold && daysLeft <= threshold {
				msg := message.NewThresholdMessage(d.Name, daysLeft)
				integration.Send(msg)
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
			msg := message.NewStatusChangedMessage(d.Name, changes)
			integration.Send(msg)
		}
	}
	return nil
}

func (d *Domain) CheckNotifications() error {
	if err := d.NotifyThreshold(); err != nil {
		return err
	}
	if err := d.NotifyStatusChange(); err != nil {
		return err
	}
	return nil
}
