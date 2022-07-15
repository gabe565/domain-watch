package domain

import (
	"github.com/araddon/dateparse"
	"github.com/gabe565/domain-expiration-notifier/internal/telegram"
	"github.com/likexian/whois-go"
	whoisparser "github.com/likexian/whois-parser-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type Domain struct {
	Name  string
	Sleep time.Duration
	Last  *whoisparser.WhoisInfo
}

func (d Domain) Whois() (whoisparser.WhoisInfo, error) {
	raw, err := whois.Whois(d.Name)
	if err != nil {
		return whoisparser.WhoisInfo{}, err
	}

	return whoisparser.Parse(raw)
}

func (d Domain) Run() error {
	if d.Sleep != 0 {
		time.Sleep(d.Sleep)
	}

	w, err := d.Whois()
	if err != nil {
		return err
	}
	defer func() {
		d.Last = &w
	}()

	l := log.WithField("domain", d.Name)

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
		telegram.Notify(w, *d.Last)
	}

	return nil
}
