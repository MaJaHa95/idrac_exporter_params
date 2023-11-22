package metrics

import (
	"strings"
	"github.com/mrlhansen/idrac_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
)

type SystemMetricGroup struct {
	SystemPowerOn      *prometheus.Desc
	SystemHealth       *prometheus.Desc
	SystemIndicatorLED *prometheus.Desc
	SystemMemorySize   *prometheus.Desc
	SystemCpuCount     *prometheus.Desc
	SystemBiosInfo     *prometheus.Desc
	SystemMachineInfo  *prometheus.Desc
}

func (metricGroup *SystemMetricGroup) GetMetricGroupType() MetricGroupType {
    return MetricGroupTypeSystem
}
func (metricGroup *SystemMetricGroup)IsEnabled(config *config.RootConfig) bool {
	return config.Collect.System
}
func (metricGroup *SystemMetricGroup) Describe(ch chan<- *prometheus.Desc) {
	ch <- metricGroup.SystemPowerOn
	ch <- metricGroup.SystemHealth
	ch <- metricGroup.SystemIndicatorLED
	ch <- metricGroup.SystemMemorySize
	ch <- metricGroup.SystemCpuCount
	ch <- metricGroup.SystemBiosInfo
	ch <- metricGroup.SystemMachineInfo
}

func (mc *SystemMetricGroup) NewSystemPowerOn(state string) prometheus.Metric {
	var value float64
	if state == "On" {
		value = 1
	}
	return prometheus.MustNewConstMetric(
		mc.SystemPowerOn,
		prometheus.GaugeValue,
		value,
	)
}

func (mc *SystemMetricGroup) NewSystemHealth(health string) prometheus.Metric {
	value := health2value(health)
	return prometheus.MustNewConstMetric(
		mc.SystemHealth,
		prometheus.GaugeValue,
		value,
		health,
	)
}

func (mc *SystemMetricGroup) NewSystemIndicatorLED(state string) prometheus.Metric {
	var value float64
	if state != "Off" {
		value = 1
	}
	return prometheus.MustNewConstMetric(
		mc.SystemIndicatorLED,
		prometheus.GaugeValue,
		value,
		state,
	)
}

func (mc *SystemMetricGroup) NewSystemMemorySize(memory float64) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.SystemMemorySize,
		prometheus.GaugeValue,
		memory,
	)
}

func (mc *SystemMetricGroup) NewSystemCpuCount(cpus int, model string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.SystemCpuCount,
		prometheus.GaugeValue,
		float64(cpus),
		strings.TrimSpace(model),
	)
}

func (mc *SystemMetricGroup) NewSystemBiosInfo(version string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.SystemBiosInfo,
		prometheus.UntypedValue,
		1.0,
		version,
	)
}

func (mc *SystemMetricGroup) NewSystemMachineInfo(manufacturer, model, serial, sku string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.SystemMachineInfo,
		prometheus.UntypedValue,
		1.0,
		manufacturer,
		model,
		serial,
		sku,
	)
}

func NewSystemMetricGroup(prefix string) *SystemMetricGroup {
    return &SystemMetricGroup {
		SystemPowerOn: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "system", "power_on"),
			"Power state of the system",
			nil, nil,
		),
		SystemHealth: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "system", "health"),
			"Health status of the system",
			[]string{"status"}, nil,
		),
		SystemIndicatorLED: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "system", "indicator_led_on"),
			"Indicator LED state of the system",
			[]string{"state"}, nil,
		),
		SystemMemorySize: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "system", "memory_size_bytes"),
			"Total memory size of the system in bytes",
			nil, nil,
		),
		SystemCpuCount: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "system", "cpu_count"),
			"Total number of CPUs in the system",
			[]string{"model"}, nil,
		),
		SystemBiosInfo: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "system", "bios_info"),
			"Information about the BIOS",
			[]string{"version"}, nil,
		),
		SystemMachineInfo: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "system", "machine_info"),
			"Information about the machine",
			[]string{"manufacturer", "model", "serial", "sku"}, nil,
		),
	}
}
