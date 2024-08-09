package metrics

import (
	"net/http"

	"github.com/gabe565/domain-watch/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(conf *config.Config) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(conf.MetricsAddress, mux)
}
