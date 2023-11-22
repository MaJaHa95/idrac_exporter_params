package collector

import (
	"runtime"
	"strings"
	"sync"
	"errors"

	"github.com/mrlhansen/idrac_exporter/internal/config"
	"github.com/mrlhansen/idrac_exporter/internal/metrics"
	"github.com/mrlhansen/idrac_exporter/internal/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)


type CacheKey struct {
    Target string
    MetricGroupType metrics.MetricGroupType
}

var mu sync.Mutex
var collectors = map[CacheKey]*Collector{}

type MetricGroupRefresher[T metrics.MetricGroup] struct {
	metricGroup    T
	refresh        func(*Client, T, chan<- prometheus.Metric) error
}

type Collector struct {
	// Internal variables
	client     *Client
	registry   *prometheus.Registry
	collected  *sync.Cond
	collecting bool
	retries    uint
	errors     uint
	builder    *strings.Builder

	selectedMetricGroupType   metrics.MetricGroupType
	
	SystemMetricGroup         MetricGroupRefresher[*metrics.SystemMetricGroup]
	SensorsMetricGroup        MetricGroupRefresher[*metrics.SensorsMetricGroup]
	PowerMetricGroup          MetricGroupRefresher[*metrics.PowerMetricGroup]
	IdracSelMetricGroup       MetricGroupRefresher[*metrics.IdracSelMetricGroup]
	StorageMetricGroup    	  MetricGroupRefresher[*metrics.StorageMetricGroup]
	MemoryMetricGroup     	  MetricGroupRefresher[*metrics.MemoryMetricGroup]

	// Exporter
	ExporterBuildInfo         *prometheus.Desc
	ExporterScrapeErrorsTotal *prometheus.Desc
}

func NewCollector(metricGroupType metrics.MetricGroupType) *Collector {
	prefix := config.Config.MetricsPrefix

	collector := &Collector{
		ExporterBuildInfo: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "exporter", "build_info"),
			"Constant metric with build information for the exporter",
			nil, prometheus.Labels{
				"version":   version.Version,
				"revision":  version.Revision,
				"goversion": runtime.Version(),
			},
		),
		ExporterScrapeErrorsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(prefix, "exporter", "scrape_errors_total"),
			"Total number of errors encountered while scraping target",
			nil, nil,
		),
	}
	
	collector.SystemMetricGroup = MetricGroupRefresher[*metrics.SystemMetricGroup] {
		metricGroup: metrics.NewSystemMetricGroup(prefix),
		refresh: func(client *Client, metricGroup *metrics.SystemMetricGroup, ch chan<- prometheus.Metric) error {
			return client.RefreshSystem(metricGroup, ch)
		},
	}

	collector.SensorsMetricGroup = MetricGroupRefresher[*metrics.SensorsMetricGroup] {
		metricGroup: metrics.NewSensorsMetricGroup(prefix),
		refresh: func(client *Client, metricGroup *metrics.SensorsMetricGroup, ch chan<- prometheus.Metric) error {
			return client.RefreshSensors(metricGroup, ch)
		},
	}

	collector.PowerMetricGroup = MetricGroupRefresher[*metrics.PowerMetricGroup] {
		metricGroup: metrics.NewPowerMetricGroup(prefix),
		refresh: func(client *Client, metricGroup *metrics.PowerMetricGroup, ch chan<- prometheus.Metric) error {
			return client.RefreshPower(metricGroup, ch)
		},
	}

	collector.IdracSelMetricGroup = MetricGroupRefresher[*metrics.IdracSelMetricGroup] {
		metricGroup: metrics.NewSelMetricGroup(prefix),
		refresh: func(client *Client, metricGroup *metrics.IdracSelMetricGroup, ch chan<- prometheus.Metric) error {
			return client.RefreshIdracSel(metricGroup, ch)
		},
	}

	collector.StorageMetricGroup = MetricGroupRefresher[*metrics.StorageMetricGroup] {
		metricGroup: metrics.NewStorageMetricGroup(prefix),
		refresh: func(client *Client, metricGroup *metrics.StorageMetricGroup, ch chan<- prometheus.Metric) error {
			return client.RefreshStorage(metricGroup, ch)
		},
	}

	collector.MemoryMetricGroup = MetricGroupRefresher[*metrics.MemoryMetricGroup] {
		metricGroup: metrics.NewMemoryMetricGroup(prefix),
		refresh: func(client *Client, metricGroup *metrics.MemoryMetricGroup, ch chan<- prometheus.Metric) error {
			return client.RefreshMemory(metricGroup, ch)
		},
	}

	collector.builder = new(strings.Builder)
	collector.collected = sync.NewCond(new(sync.Mutex))
	collector.registry = prometheus.NewRegistry()
	collector.registry.Register(collector)

	collector.selectedMetricGroupType = metricGroupType

	return collector
}

func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.ExporterBuildInfo
	ch <- collector.ExporterScrapeErrorsTotal
	
	collector.SystemMetricGroup.metricGroup.Describe(ch)
	collector.SensorsMetricGroup.metricGroup.Describe(ch)
	collector.PowerMetricGroup.metricGroup.Describe(ch)
	collector.IdracSelMetricGroup.metricGroup.Describe(ch)
	collector.StorageMetricGroup.metricGroup.Describe(ch)
	collector.MemoryMetricGroup.metricGroup.Describe(ch)
}

func tryRefresh[T metrics.MetricGroup](collector *Collector, metricGroup MetricGroupRefresher[T], ch chan<- prometheus.Metric) error {
	if collector.selectedMetricGroupType == metricGroup.metricGroup.GetMetricGroupType() {
		if !metricGroup.metricGroup.IsEnabled(&config.Config) {
			return errors.New("The requested metric group isn't enabled")
		}
	} else if collector.selectedMetricGroupType != metrics.MetricGroupTypeAny || !metricGroup.metricGroup.IsEnabled(&config.Config) {
		return nil
	}

	return metricGroup.refresh(collector.client, metricGroup.metricGroup, ch)
}

func (collector *Collector) Collect(ch chan<- prometheus.Metric) {

	if err := tryRefresh(collector, collector.SystemMetricGroup, ch); err != nil {
        collector.errors++
    }
	
	if err := tryRefresh(collector, collector.SensorsMetricGroup, ch); err != nil {
        collector.errors++
    }
	
	if err := tryRefresh(collector, collector.PowerMetricGroup, ch); err != nil {
        collector.errors++
    }
	
	if err := tryRefresh(collector, collector.IdracSelMetricGroup, ch); err != nil {
        collector.errors++
    }
	
	if err := tryRefresh(collector, collector.StorageMetricGroup, ch); err != nil {
        collector.errors++
    }
	
	if err := tryRefresh(collector, collector.MemoryMetricGroup, ch); err != nil {
        collector.errors++
    }

	ch <- prometheus.MustNewConstMetric(collector.ExporterBuildInfo, prometheus.UntypedValue, 1)
	ch <- prometheus.MustNewConstMetric(collector.ExporterScrapeErrorsTotal, prometheus.GaugeValue, float64(collector.errors))
}

func (collector *Collector) Gather() (string, error) {
	collector.collected.L.Lock()

	// If a collection is already in progress wait for it to complete and return the cached data
	if collector.collecting {
		collector.collected.Wait()
		metrics := collector.builder.String()
		collector.collected.L.Unlock()
		return metrics, nil
	}

	// Set collecting to true and let other goroutines enter in critical section
	collector.collecting = true
	collector.collected.L.Unlock()

	// Defer set collecting to false and wake waiting goroutines
	defer func() {
		collector.collected.L.Lock()
		collector.collected.Broadcast()
		collector.collecting = false
		collector.collected.L.Unlock()
	}()

	// Collect metrics
	collector.builder.Reset()

	m, err := collector.registry.Gather()
	if err != nil {
		return "", err
	}

	for i := range m {
		expfmt.MetricFamilyToText(collector.builder, m[i])
	}

	return collector.builder.String(), nil
}

// Resets an existing collector of the given target
func Reset(target string, metricGroupType metrics.MetricGroupType) {
	key := CacheKey{Target: target, MetricGroupType: metricGroupType}

	mu.Lock()
	_, ok := collectors[key]
	if ok {
		delete(collectors, key)
	}
	mu.Unlock()
}

func GetCollector(target string, metricGroupType metrics.MetricGroupType) (*Collector, error) {
	key := CacheKey{Target: target, MetricGroupType: metricGroupType}

	mu.Lock()
	collector, ok := collectors[key]
	if !ok {
		collector = NewCollector(metricGroupType)
		collectors[key] = collector
	}
	mu.Unlock()

	// Do not act concurrently on the same host
	collector.collected.L.Lock()
	defer collector.collected.L.Unlock()

	// Find (potentially cached) Redfish client
	if collector.client == nil {
		client, err := GetClient(target)

		if err != nil {
			return nil, err
		}

		collector.client = client
	}

	return collector, nil
}
