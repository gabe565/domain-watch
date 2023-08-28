package cmd

import (
	"errors"
	"time"

	"github.com/gabe565/domain-watch/internal/domain"
	"github.com/gabe565/domain-watch/internal/integration"
	"github.com/gabe565/domain-watch/internal/metrics"
	"github.com/gabe565/domain-watch/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "domain-watch [flags] domain...",
		DisableAutoGenTag: true,
		PreRunE:           preRun,
		RunE:              run,
		ValidArgsFunction: util.NoFileComp,
	}

	if err := integration.Flags(cmd); err != nil {
		panic(err)
	}
	metrics.Flags(cmd)
	registerCompletionFlag(cmd)
	registerEveryFlag(cmd)
	registerLogFlags(cmd)
	registerSleepFlag(cmd)
	registerThresholdFlag(cmd)

	return cmd
}

var domainNames []string

func preRun(cmd *cobra.Command, args []string) (err error) {
	initViper()
	initLog(cmd)

	if completionFlag != "" {
		return completion(cmd, domainNames)
	}

	if v := viper.GetStringSlice("domains"); v != nil {
		args = append(domainNames, v...)
	}
	domainNames = args

	if len(domainNames) == 0 {
		return errors.New("missing domain")
	}

	if err := integration.Setup(); err != nil {
		return err
	}

	return nil
}

func run(cmd *cobra.Command, _ []string) (err error) {
	cmd.SilenceUsage = true

	go func() {
		if err := metrics.Serve(cmd); err != nil {
			log.Error(err)
		}
	}()

	domains := domain.Domains{
		Sleep:   viper.GetDuration("sleep"),
		Domains: make([]*domain.Domain, 0, len(domainNames)),
	}
	for _, domainName := range domainNames {
		domains.Add(domain.Domain{Name: domainName})
	}

	domains.Tick()

	every := viper.GetDuration("every")
	if every != 0 {
		log.Info("running as cron")

		ticker := time.NewTicker(every)
		for range ticker.C {
			domains.Tick()
		}
	}

	return nil
}
