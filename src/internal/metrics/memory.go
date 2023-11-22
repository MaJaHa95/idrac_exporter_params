package metrics

import (
	"fmt"
	
	"github.com/mrlhansen/idrac_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
)



type MemoryMetricGroup struct {
    MemoryModuleInfo          *prometheus.Desc
    MemoryModuleHealth        *prometheus.Desc
    MemoryModuleCapacity      *prometheus.Desc
    MemoryModuleSpeed         *prometheus.Desc
}

func (metricGroup *MemoryMetricGroup) GetMetricGroupType() MetricGroupType {
    return MetricGroupTypeMemory
}

func (metricGroup *MemoryMetricGroup) IsEnabled(config *config.RootConfig) bool {
    return config.Collect.Memory
}

func (metricGroup *MemoryMetricGroup) Describe(ch chan<- *prometheus.Desc) {
    ch <- metricGroup.MemoryModuleInfo
    ch <- metricGroup.MemoryModuleHealth
    ch <- metricGroup.MemoryModuleCapacity
    ch <- metricGroup.MemoryModuleSpeed
}


func (mc *MemoryMetricGroup) NewMemoryModuleInfo(id, name, manufacturer, memtype, serial, ecc string, rank int) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.MemoryModuleInfo,
		prometheus.UntypedValue,
		1.0,
		id,
		ecc,
		manufacturer,
		memtype,
		name,
		serial,
		fmt.Sprint(rank),
	)
}

func (mc *MemoryMetricGroup) NewMemoryModuleHealth(id, health string) prometheus.Metric {
	value := health2value(health)
	return prometheus.MustNewConstMetric(
		mc.MemoryModuleHealth,
		prometheus.GaugeValue,
		value,
		id,
		health,
	)
}

func (mc *MemoryMetricGroup) NewMemoryModuleCapacity(id string, capacity int) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.MemoryModuleCapacity,
		prometheus.GaugeValue,
		float64(capacity),
		id,
	)
}

func (mc *MemoryMetricGroup) NewMemoryModuleSpeed(id string, speed int) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.MemoryModuleSpeed,
		prometheus.GaugeValue,
		float64(speed),
		id,
	)
}

// Instance initialization
func NewMemoryMetricGroup(prefix string) *MemoryMetricGroup {
    return &MemoryMetricGroup {
		MemoryModuleInfo: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "memory_module", "info"),
			"Information about memory modules",
			[]string{"id", "ecc", "manufacturer", "type", "name", "serial", "rank"}, nil,
		),
		MemoryModuleHealth: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "memory_module", "health"),
			"Health status for memory modules",
			[]string{"id", "status"}, nil,
		),
		MemoryModuleCapacity: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "memory_module", "capacity_bytes"),
			"Capacity of memory modules in bytes",
			[]string{"id"}, nil,
		),
		MemoryModuleSpeed: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "memory_module", "speed_mhz"),
			"Speed of memory modules in Mhz",
			[]string{"id"}, nil,
		),
	}
}

