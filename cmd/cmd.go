package cmd

import (
	"context"
	"log/slog"
	"time"

	"gabe565.com/domain-watch/internal/config"
	"gabe565.com/domain-watch/internal/domain"
	"gabe565.com/domain-watch/internal/integration"
	"gabe565.com/domain-watch/internal/metrics"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra"
)

func New(opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "domain-watch [flags] domain...",
		DisableAutoGenTag: true,
		RunE:              run,
		ValidArgsFunction: cobra.NoFileCompletions,
	}
	cfg := config.New()
	cfg.RegisterFlags(cmd)
	config.RegisterCompletions(cmd)
	cmd.SetContext(config.NewContext(context.Background(), cfg))
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	conf, err := config.Load(cmd, args)
	if err != nil {
		return err
	}

	if conf.Completion != "" {
		return completion(cmd, conf.Completion)
	}

	slog.Info("Domain Watch", "version", cobrax.GetVersion(cmd), "commit", cobrax.GetCommit(cmd))

	integrations, err := integration.Setup(cmd.Context(), conf)
	if err != nil {
		return err
	}

	cmd.SilenceUsage = true

	if conf.MetricsEnabled {
		go func() {
			if err := metrics.Serve(conf); err != nil {
				slog.Error("Failed to serve metrics", "error", err)
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

	domains.Tick(cmd.Context(), integrations)

	if conf.Every != 0 {
		slog.Info("Running as cron", "interval", conf.Every)

		ticker := time.NewTicker(conf.Every)
		for range ticker.C {
			domains.Tick(cmd.Context(), integrations)
		}
	}

	return nil
}
