package prometheus_exporter

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
)

func RunHTTPServer(logger *slog.Logger, remoteAddresses []string) {
	go CollectPingStats(logger, remoteAddresses)

	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		logger.Error("failed to run http server", "error", err)
		return
	}
}
