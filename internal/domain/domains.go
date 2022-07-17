package domain

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type Domains struct {
	Sleep   time.Duration
	Domains []Domain
}

func (d *Domains) Add(domain Domain) {
	d.Domains = append(d.Domains, domain)
}

func (d Domains) Tick() {
	for i, domain := range d.Domains {
		if i != 0 {
			time.Sleep(d.Sleep)
		}
		if err := domain.Run(); err != nil {
			log.WithError(err).Error("failed to fetch whois")
		}
	}
}
