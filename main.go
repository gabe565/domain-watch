package main

import (
	"github.com/gabe565/domain-watch/cmd"
	log "github.com/sirupsen/logrus"
)

//go:generate go run ./internal/cmd/docs

func main() {
	if err := cmd.Command.Execute(); err != nil {
		log.Fatal(err)
	}
}
