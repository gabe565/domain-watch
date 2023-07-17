package metrics

import (
	"net/http"

	"github.com/gabe565/domain-watch/internal/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Flags(cmd *cobra.Command) {
	cmd.Flags().Bool("metrics-enabled", false, "Enables Prometheus metrics API")
	if err := viper.BindPFlag("metrics.enabled", cmd.Flags().Lookup("metrics-enabled")); err != nil {
		panic(err)
	}
	if err := cmd.RegisterFlagCompletionFunc("metrics-enabled", util.NoFileComp); err != nil {
		panic(err)
	}

	cmd.Flags().String("metrics-address", ":9090", "Prometheus metrics API listen address")
	if err := viper.BindPFlag("metrics.address", cmd.Flags().Lookup("metrics-address")); err != nil {
		panic(err)
	}
	if err := cmd.RegisterFlagCompletionFunc("metrics-address", util.NoFileComp); err != nil {
		panic(err)
	}
}

func Serve(cmd *cobra.Command) error {
	if !viper.GetBool("metrics.enabled") {
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(viper.GetString("metrics.address"), mux)
}
