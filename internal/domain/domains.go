package domain

import log "github.com/sirupsen/logrus"

type Domains []Domain

func (d Domains) Tick() {
	for _, domain := range d {
		if err := domain.Run(); err != nil {
			log.WithError(err).Error("failed to fetch whois")
		}
	}
}
