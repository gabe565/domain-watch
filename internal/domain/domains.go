package domain

import (
	"context"
	"time"

	"gabe565.com/domain-watch/internal/integration"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//nolint:gochecknoglobals
var (
	lastTickMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "domain_watch",
		Name:      "last_fetch_seconds",
		Help:      "Unix timestamp for when the last fetch occurred.",
	})

	domainCountMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "domain_watch",
		Name:      "domains",
		Help:      "Number of domains that are being watched.",
	})

	successMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "domain_watch",
		Name:      "update_success",
		Help:      "Whether the last fetch succeeded.",
	}, []string{"domain"})

	expirationMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "domain_watch",
		Name:      "expires_timestamp_seconds",
		Help:      "Unix timestamp for when the domain will expire.",
	}, []string{"domain"})
)

type Domains struct {
	Sleep   time.Duration
	Domains []*Domain
}

func (d *Domains) Add(domain Domain) {
	domainCountMetric.Add(1)
	d.Domains = append(d.Domains, &domain)
}

func (d Domains) Tick(ctx context.Context, integrations integration.Integrations) {
	defer func() {
		lastTickMetric.SetToCurrentTime()
	}()

	for i, domain := range d.Domains {
		if i != 0 {
			time.Sleep(d.Sleep)
		}
		if err := domain.Run(ctx, integrations); err == nil {
			successMetric.With(prometheus.Labels{"domain": domain.Name}).Set(1)
		} else {
			successMetric.With(prometheus.Labels{"domain": domain.Name}).Set(0)
			domain.Log().Error("Domain update failed", "error", err)
		}
		if domain.ExpiresAt.Unix() > 0 {
			expirationMetric.With(prometheus.Labels{"domain": domain.Name}).Set(float64(domain.ExpiresAt.Unix()))
		}
	}
}
