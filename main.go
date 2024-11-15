package main

import (
	"log/slog"
	"os"

	"gabe565.com/domain-watch/cmd"
	"gabe565.com/domain-watch/internal/config"
	"gabe565.com/utils/cobrax"
)

var version string

func main() {
	config.InitLog(os.Stderr, slog.LevelInfo, config.FormatAuto)
	root := cmd.New(cobrax.WithVersion(version))
	if err := root.Execute(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
