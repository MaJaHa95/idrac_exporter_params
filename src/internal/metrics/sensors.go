package metrics

import (
	"github.com/mrlhansen/idrac_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
)



type SensorsMetricGroup struct {
    SensorsTemperature *prometheus.Desc
    SensorsFanSpeed    *prometheus.Desc
}

func (metricGroup *SensorsMetricGroup) GetMetricGroupType() MetricGroupType {
    return MetricGroupTypeSensors
}

func (metricGroup *SensorsMetricGroup) IsEnabled(config *config.RootConfig) bool {
    return config.Collect.Sensors
}

func (metricGroup *SensorsMetricGroup) Describe(ch chan<- *prometheus.Desc) {
    ch <- metricGroup.SensorsTemperature
    ch <- metricGroup.SensorsFanSpeed
}


func (mc *SensorsMetricGroup) NewSensorsTemperature(temperature float64, id, name, units string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.SensorsTemperature,
		prometheus.GaugeValue,
		temperature,
		id,
		name,
		units,
	)
}

func (mc *SensorsMetricGroup) NewSensorsFanSpeed(speed float64, id, name, units string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.SensorsFanSpeed,
		prometheus.GaugeValue,
		speed,
		id,
		name,
		units,
	)
}

// Instance initialization
func NewSensorsMetricGroup(prefix string) *SensorsMetricGroup {
    return &SensorsMetricGroup {
		SensorsTemperature: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "sensors", "temperature"),
			"Sensors reporting temperature measurements",
			[]string{"id", "name", "units"}, nil,
		),
		SensorsFanSpeed: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "sensors", "fan_speed"),
			"Sensors reporting fan speed measurements",
			[]string{"id", "name", "units"}, nil,
		),
	}
}
