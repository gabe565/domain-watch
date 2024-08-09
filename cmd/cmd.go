package cmd

import (
	"context"
	"time"

	"github.com/gabe565/domain-watch/internal/config"
	"github.com/gabe565/domain-watch/internal/domain"
	"github.com/gabe565/domain-watch/internal/integration"
	"github.com/gabe565/domain-watch/internal/metrics"
	"github.com/gabe565/domain-watch/internal/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "domain-watch [flags] domain...",
		DisableAutoGenTag: true,
		RunE:              run,
		ValidArgsFunction: util.NoFileComp,
	}
	cfg := config.New()
	cfg.RegisterFlags(cmd)
	config.RegisterCompletions(cmd)
	cmd.SetContext(config.NewContext(context.Background(), cfg))
	return cmd
}

func run(cmd *cobra.Command, args []string) (err error) {
	conf, err := config.Load(cmd, args)
	if err != nil {
		return err
	}

	if conf.Completion != "" {
		return completion(cmd, conf.Completion)
	}

	if err := integration.Setup(conf); err != nil {
		return err
	}

	cmd.SilenceUsage = true

	if conf.MetricsEnabled {
		go func() {
			if err := metrics.Serve(conf); err != nil {
				log.Error(err)
			}
		}()
	}

	domains := domain.Domains{
		Sleep:   conf.Sleep,
		Domains: make([]*domain.Domain, 0, len(conf.Domains)),
	}
	for _, domainName := range conf.Domains {
		domains.Add(domain.New(conf, domainName))
	}

	domains.Tick()

	if conf.Every != 0 {
		log.Info("running as cron")

		ticker := time.NewTicker(conf.Every)
		for range ticker.C {
			domains.Tick()
		}
	}

	return nil
}
