package metrics

import (
	"github.com/mrlhansen/idrac_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
)

type PowerMetricGroup struct {
    PowerSupplyOutputWatts      *prometheus.Desc
    PowerSupplyInputWatts       *prometheus.Desc
    PowerSupplyCapacityWatts    *prometheus.Desc
    PowerSupplyInputVoltage     *prometheus.Desc
    PowerSupplyEfficiencyPercent *prometheus.Desc
    PowerControlConsumedWatts   *prometheus.Desc
    PowerControlCapacityWatts   *prometheus.Desc
    PowerControlMinConsumedWatts *prometheus.Desc
    PowerControlMaxConsumedWatts *prometheus.Desc
    PowerControlAvgConsumedWatts *prometheus.Desc
    PowerControlInterval        *prometheus.Desc
}

func (metricGroup *PowerMetricGroup) GetMetricGroupType() MetricGroupType {
    return MetricGroupTypePower
}

func (metricGroup *PowerMetricGroup) IsEnabled(config *config.RootConfig) bool {
    return config.Collect.Power
}

func (metricGroup *PowerMetricGroup) Describe(ch chan<- *prometheus.Desc) {
    ch <- metricGroup.PowerSupplyOutputWatts
    ch <- metricGroup.PowerSupplyInputWatts
    ch <- metricGroup.PowerSupplyCapacityWatts
    ch <- metricGroup.PowerSupplyInputVoltage
    ch <- metricGroup.PowerSupplyEfficiencyPercent
    ch <- metricGroup.PowerControlConsumedWatts
    ch <- metricGroup.PowerControlCapacityWatts
    ch <- metricGroup.PowerControlMinConsumedWatts
    ch <- metricGroup.PowerControlMaxConsumedWatts
    ch <- metricGroup.PowerControlAvgConsumedWatts
    ch <- metricGroup.PowerControlInterval
}

func (mc *PowerMetricGroup) NewPowerSupplyInputWatts(value float64, id string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.PowerSupplyInputWatts,
		prometheus.GaugeValue,
		value,
		id,
	)
}

func (mc *PowerMetricGroup) NewPowerSupplyInputVoltage(value float64, id string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.PowerSupplyInputVoltage,
		prometheus.GaugeValue,
		value,
		id,
	)
}

func (mc *PowerMetricGroup) NewPowerSupplyOutputWatts(value float64, id string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.PowerSupplyOutputWatts,
		prometheus.GaugeValue,
		value,
		id,
	)
}

func (mc *PowerMetricGroup) NewPowerSupplyCapacityWatts(value float64, id string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.PowerSupplyCapacityWatts,
		prometheus.GaugeValue,
		value,
		id,
	)
}

func (mc *PowerMetricGroup) NewPowerSupplyEfficiencyPercent(value float64, id string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.PowerSupplyEfficiencyPercent,
		prometheus.GaugeValue,
		value,
		id,
	)
}

func (mc *PowerMetricGroup) NewPowerControlConsumedWatts(value float64, id, name string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.PowerControlConsumedWatts,
		prometheus.GaugeValue,
		value,
		id,
		name,
	)
}

func (mc *PowerMetricGroup) NewPowerControlCapacityWatts(value float64, id, name string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.PowerControlCapacityWatts,
		prometheus.GaugeValue,
		value,
		id,
		name,
	)
}

func (mc *PowerMetricGroup) NewPowerControlMinConsumedWatts(value float64, id, name string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.PowerControlMinConsumedWatts,
		prometheus.GaugeValue,
		value,
		id,
		name,
	)
}

func (mc *PowerMetricGroup) NewPowerControlMaxConsumedWatts(value float64, id, name string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.PowerControlMaxConsumedWatts,
		prometheus.GaugeValue,
		value,
		id,
		name,
	)
}

func (mc *PowerMetricGroup) NewPowerControlAvgConsumedWatts(value float64, id, name string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.PowerControlAvgConsumedWatts,
		prometheus.GaugeValue,
		value,
		id,
		name,
	)
}

func (mc *PowerMetricGroup) NewPowerControlInterval(interval int, id, name string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		mc.PowerControlInterval,
		prometheus.GaugeValue,
		float64(interval),
		id,
		name,
	)
}

// Instance initialization
func NewPowerMetricGroup(prefix string) *PowerMetricGroup {
    return &PowerMetricGroup {
		PowerSupplyOutputWatts: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "power_supply", "output_watts"),
			"Power supply output in watts",
			[]string{"id"}, nil,
		),
		PowerSupplyInputWatts: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "power_supply", "input_watts"),
			"Power supply input in watts",
			[]string{"id"}, nil,
		),
		PowerSupplyCapacityWatts: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "power_supply", "capacity_watts"),
			"Power supply capacity in watts",
			[]string{"id"}, nil,
		),
		PowerSupplyInputVoltage: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "power_supply", "input_voltage"),
			"Power supply input voltage",
			[]string{"id"}, nil,
		),
		PowerSupplyEfficiencyPercent: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "power_supply", "efficiency_percent"),
			"Power supply efficiency in percentage",
			[]string{"id"}, nil,
		),
		PowerControlConsumedWatts: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "power_control", "consumed_watts"),
			"Consumption of power control system in watts",
			[]string{"id", "name"}, nil,
		),
		PowerControlCapacityWatts: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "power_control", "capacity_watts"),
			"Capacity of power control system in watts",
			[]string{"id", "name"}, nil,
		),
		PowerControlMinConsumedWatts: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "power_control", "min_consumed_watts"),
			"Minimum consumption of power control system during the reported interval",
			[]string{"id", "name"}, nil,
		),
		PowerControlMaxConsumedWatts: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "power_control", "max_consumed_watts"),
			"Maximum consumption of power control system during the reported interval",
			[]string{"id", "name"}, nil,
		),
		PowerControlAvgConsumedWatts: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "power_control", "avg_consumed_watts"),
			"Average consumption of power control system during the reported interval",
			[]string{"id", "name"}, nil,
		),
		PowerControlInterval: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "power_control", "interval_in_minutes"),
			"Interval for measurements of power control system",
			[]string{"id", "name"}, nil,
		),
	}
}
