package main

import (
	"github.com/gabe565/domain-expiration-notifier/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Command.Execute(); err != nil {
		log.Fatal(err)
	}
}
