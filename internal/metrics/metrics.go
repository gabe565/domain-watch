package metrics

import (
	"net/http"
	"time"

	"gabe565.com/domain-watch/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Serve(conf *config.Config) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:        conf.MetricsAddress,
		Handler:     mux,
		ReadTimeout: 10 * time.Second,
	}
	return server.ListenAndServe()
}
