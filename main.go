package main

import (
	"log/slog"
	"os"
	"ping-stats-prometheus-exporter/prometheus_exporter"
	"strings"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	remoteAddrsToPing := strings.Split(os.Getenv("PING_WORTHY_REMOTE_ADDRESSES"), ",")
	if len(remoteAddrsToPing) == 0 {
		logger.Error("no remote address found to be pinged")
		return
	}

	prometheus_exporter.RunHTTPServer(logger, remoteAddrsToPing)
}
