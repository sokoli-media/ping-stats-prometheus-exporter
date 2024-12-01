package prometheus_exporter

import (
	probing "github.com/prometheus-community/pro-bing"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log/slog"
	"time"
)

var pingLabels = []string{"remote_addr"}
var pingAvgReturnTime = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "ping_stats_avg_response_time"}, pingLabels)
var pingMinReturnTime = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "ping_stats_min_response_time"}, pingLabels)
var pingMaxReturnTime = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "ping_stats_max_response_time"}, pingLabels)
var pingSentPackages = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "ping_stats_sent_packages"}, pingLabels)
var pingPacketLoss = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "ping_stats_packet_loss_percentage"}, pingLabels)
var pingLastUpdate = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "ping_stats_last_update"}, pingLabels)

func ping(url string) (*probing.Statistics, error) {
	pinger, err := probing.NewPinger(url)
	if err != nil {
		return nil, err
	}
	pinger.Count = 3
	pinger.Timeout = 30 * time.Second

	err = pinger.Run()
	if err != nil {
		return nil, err
	}
	return pinger.Statistics(), nil
}

func pingAndReport(logger *slog.Logger, remoteAddress string) {
	logger.Info("collecting ping metrics", "remote_addr", remoteAddress)

	stats, err := ping(remoteAddress)
	if err != nil {
		logger.Error("couldn't ping", "remote_addr", remoteAddress, "error", err)
		return
	}

	labels := prometheus.Labels{"remote_addr": remoteAddress}
	pingAvgReturnTime.With(labels).Set(stats.AvgRtt.Seconds())
	pingMinReturnTime.With(labels).Set(stats.MinRtt.Seconds())
	pingMaxReturnTime.With(labels).Set(stats.MaxRtt.Seconds())
	pingSentPackages.With(labels).Set(float64(stats.PacketsSent))
	pingPacketLoss.With(labels).Set(stats.PacketLoss / 100)

	pingLastUpdate.With(labels).SetToCurrentTime()
}

func CollectPingStats(logger *slog.Logger, remoteAddresses []string) {
	ticker := time.NewTicker(15 * time.Second)
	for {
		<-ticker.C

		for _, address := range remoteAddresses {
			go pingAndReport(logger, address)
		}
	}
}
