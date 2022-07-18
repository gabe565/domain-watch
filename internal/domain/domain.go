package domain

import (
	"github.com/araddon/dateparse"
	"github.com/gabe565/domain-watch/internal/telegram"
	"github.com/likexian/whois-go"
	whoisparser "github.com/likexian/whois-parser-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type Domain struct {
	Name string
	Last *whoisparser.WhoisInfo
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

func (d *Domain) Run() error {
	w, err := d.Whois()
	if err != nil {
		return err
	}
	defer func() {
		d.Last = &w
	}()

	l := d.Log()

	if w.Domain.ExpirationDate != "" {
		var date time.Time
		date, err = dateparse.ParseStrict(w.Domain.ExpirationDate)
		if err != nil {
			l.WithError(err).Warn("failed to parse expiration date")
		} else {
			left := date.Sub(time.Now()).Truncate(24 * time.Hour)
			l.WithFields(log.Fields{
				"expires":   date,
				"days_left": left.Hours() / 24.0,
			}).Info("fetched whois")
		}
	} else {
		l.Info("domain does not have an expiration date")
	}

	if d.Last != nil {
		if err := telegram.Notify(w, *d.Last); err != nil {
			return err
		}
	}

	return nil
}
