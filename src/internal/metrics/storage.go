package metrics

import (
	"fmt"
	"github.com/mrlhansen/idrac_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
)



type StorageMetricGroup struct {
    DriveInfo     *prometheus.Desc
    DriveHealth   *prometheus.Desc
    DriveCapacity *prometheus.Desc
}

func (metricGroup *StorageMetricGroup) GetMetricGroupType() MetricGroupType {
    return MetricGroupTypeStorage
}
func (metricGroup *StorageMetricGroup) IsEnabled(config *config.RootConfig) bool {
    return config.Collect.Storage
}

func (metricGroup *StorageMetricGroup) Describe(ch chan<- *prometheus.Desc) {
    ch <- metricGroup.DriveInfo
    ch <- metricGroup.DriveHealth
    ch <- metricGroup.DriveCapacity
}

func (mc *StorageMetricGroup) NewDriveInfo(id, name, manufacturer, model, serial, mediatype, protocol string, slot int) prometheus.Metric {
	var slotstr string

	if slot < 0 {
		slotstr = ""
	} else {
		slotstr = fmt.Sprint(slot)
	}

	return prometheus.MustNewConstMetric(
		mc.DriveInfo,
		prometheus.UntypedValue,
		1.0,
		id,
		manufacturer,
		mediatype,
		model,
		name,
		protocol,
		serial,
		slotstr,
	)
}

func (mc *StorageMetricGroup) NewDriveHealth(id, health string) prometheus.Metric {
	value := health2value(health)
	return prometheus.MustNewConstMetric(
		mc.DriveHealth,
		prometheus.GaugeValue,
		value,
		id,
		health,
	)
}

func (mc *StorageMetricGroup) NewDriveCapacity(id string, capacity int) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.DriveCapacity,
		prometheus.GaugeValue,
		float64(capacity),
		id,
	)
}


// Instance initialization
func NewStorageMetricGroup(prefix string) *StorageMetricGroup {
    return &StorageMetricGroup {
		DriveInfo: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "drive", "info"),
			"Information about disk drives",
			[]string{"id", "manufacturer", "mediatype", "model", "name", "protocol", "serial", "slot"}, nil,
		),
		DriveHealth: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "drive", "health"),
			"Health status for disk drives",
			[]string{"id", "status"}, nil,
		),
		DriveCapacity: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "drive", "capacity_bytes"),
			"Capacity of disk drives in bytes",
			[]string{"id"}, nil,
		),
	}
}
