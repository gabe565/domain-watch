package main

import (
	"github.com/gabe565/domain-watch/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	rootCmd := cmd.NewCommand()
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
