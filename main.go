package main

import (
	"gabe565.com/domain-watch/cmd"
	"gabe565.com/utils/cobrax"
	log "github.com/sirupsen/logrus"
)

var version string

func main() {
	root := cmd.New(cobrax.WithVersion(version))
	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}
