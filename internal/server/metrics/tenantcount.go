package metrics

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sinkingpoint/kiora/lib/kiora/config"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb"
	"github.com/sinkingpoint/kiora/lib/kiora/kioradb/query"
	"github.com/sinkingpoint/kiora/lib/kiora/model"

	"github.com/prometheus/client_golang/prometheus"
)

var _ = prometheus.Collector(&TenantCountCollector{})

var alertCountDesc = prometheus.NewDesc("kiora_alert_state", "Number of alerts in the system", []string{
	"tenant",
	"state",
}, nil)

// TenantCountCollector is a prometheus.Collector that collects the number of alerts in the system,
// broken down by tenant and state.
type TenantCountCollector struct {
	globals *config.Globals
	db      kioradb.DB
}

func NewTenantCountCollector(globals *config.Globals, db kioradb.DB) *TenantCountCollector {
	return &TenantCountCollector{
		globals: globals,
		db:      db,
	}
}

// Collect implements prometheus.Collector.
func (t *TenantCountCollector) Collect(ch chan<- prometheus.Metric) {
	type state struct {
		tenant config.Tenant
		state  model.AlertStatus
	}

	states := map[state]int64{}

	alerts := t.db.QueryAlerts(context.Background(), query.NewAlertQuery(query.MatchAll()))
	for _, alert := range alerts {
		tenant, err := t.globals.Tenanter.GetTenant(context.Background(), &alert)
		if err != nil {
			log.Debug().Err(err).Interface("alert", alert).Msg("Failed to get tenant")
			tenant = "error"
		}

		state := state{
			tenant: tenant,
			state:  alert.Status,
		}

		if count, ok := states[state]; ok {
			states[state] = count + 1
		} else {
			states[state] = 1
		}
	}

	for state, count := range states {
		ch <- prometheus.MustNewConstMetric(alertCountDesc,
			prometheus.GaugeValue,
			float64(count),
			string(state.tenant),
			string(state.state),
		)
	}
}

// Describe implements prometheus.Collector.
func (*TenantCountCollector) Describe(c chan<- *prometheus.Desc) {
	c <- alertCountDesc
}
