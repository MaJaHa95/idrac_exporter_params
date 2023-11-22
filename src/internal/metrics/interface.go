package metrics

import (
	"github.com/mrlhansen/idrac_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricGroupType int

const (
	MetricGroupTypeAny MetricGroupType = iota
	MetricGroupTypeSystem
	MetricGroupTypeSensors
	MetricGroupTypePower
	MetricGroupTypeIdracSel
	MetricGroupTypeStorage
	MetricGroupTypeMemory
)

type MetricGroup interface {
	GetMetricGroupType() MetricGroupType

	IsEnabled(config *config.RootConfig) bool
    Describe(ch chan<- *prometheus.Desc)
}
