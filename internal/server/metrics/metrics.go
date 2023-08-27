package metrics

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
)

func Register(router *mux.Router) {
	router.Handle("/metrics", promhttp.Handler())
}

func RegisterMetricsCollectors(globals *config.Globals, db kioradb.DB) {
	prometheus.MustRegister(NewTenantCountCollector(globals, db))
}
