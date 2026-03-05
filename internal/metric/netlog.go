package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var netlogCreatedTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "ai_app_netlog_created_total",
		Help: "Количество успешно созданных записей netlog",
	},
	[]string{"app_name"},
)

func IncNetlogCreated(appName string) {
	netlogCreatedTotal.WithLabelValues(appName).Inc()
}
