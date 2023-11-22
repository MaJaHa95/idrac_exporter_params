package metrics

import (
	"time"
	"github.com/mrlhansen/idrac_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
)


type IdracSelMetricGroup struct {
    SelEntry        *prometheus.Desc
}

func (metricGroup *IdracSelMetricGroup) GetMetricGroupType() MetricGroupType {
    return MetricGroupTypeIdracSel
}

func (metricGroup *IdracSelMetricGroup) IsEnabled(config *config.RootConfig) bool {
    return config.Collect.SEL
}

func (metricGroup *IdracSelMetricGroup) Describe(ch chan<- *prometheus.Desc) {
    ch <- metricGroup.SelEntry
}

func (mc *IdracSelMetricGroup) NewSelEntry(id string, message string, component string, severity string, created time.Time) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.SelEntry,
		prometheus.CounterValue,
		float64(created.Unix()),
		id,
		message,
		component,
		severity,
	)
}

// Instance initialization
func NewSelMetricGroup(prefix string) *IdracSelMetricGroup {
    return &IdracSelMetricGroup {
		SelEntry: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "sel", "entry"),
			"Entry from the system event log",
			[]string{"id", "message", "component", "severity"}, nil,
		),
	}
}
