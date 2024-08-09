package integration

import (
	"errors"

	"github.com/gabe565/domain-watch/internal/config"
	"github.com/gabe565/domain-watch/internal/util"
	log "github.com/sirupsen/logrus"
)

type Integration interface {
	Setup(*config.Config) error
	Send(string) error
}

var Default = map[string]Integration{
	"telegram": &Telegram{},
	"gotify":   &Gotify{},
}

func Setup(conf *config.Config) error {
	var configured uint8

	for _, integration := range Default {
		err := integration.Setup(conf)
		if err != nil {
			if errors.Is(err, util.ErrNotConfigured) {
				continue
			}
			return err
		}
		configured += 1
	}

	if configured == 0 {
		log.Warn("no integrations were configured")
	}

	return nil
}

func Send(message string) {
	for name, integration := range Default {
		if err := integration.Send(message); err != nil {
			log.WithField("integration", name).Error(err)
		}
	}
}

func Get(key string) Integration {
	return Default[key]
}
