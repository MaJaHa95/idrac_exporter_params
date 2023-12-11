package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"errors"

	"github.com/mrlhansen/idrac_exporter/internal/collector"
	"github.com/mrlhansen/idrac_exporter/internal/logging"
	"github.com/mrlhansen/idrac_exporter/internal/metrics"
	"github.com/mrlhansen/idrac_exporter/internal/config"
)

const (
	contentTypeHeader     = "Content-Type"
	contentEncodingHeader = "Content-Encoding"
	acceptEncodingHeader  = "Accept-Encoding"
)

var gzipPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

func getTargetParam(req *http.Request) (string, error) {
	target := req.URL.Query().Get("target")

	forcedTarget := config.Config.SingleHost
	
	if target == "" {
		if forcedTarget == "" {
			logging.Errorf(nil, "Received request from %s without 'target' parameter", req.Host)

			return "", errors.New("Query parameter 'target' is mandatory")
		} else {
			return forcedTarget, nil
		}
	} else if forcedTarget != "" && target != forcedTarget {
		logging.Errorf(nil, "Received request from %s with invalid 'target' parameter. Expected '%s', got '%s'", req.Host, forcedTarget, target)

		return "", fmt.Errorf("Query parameter 'target' must be either omitted or set to %s", forcedTarget)
	}

	return target, nil
}

func getMetricGroupName(val metrics.MetricGroupType) (string, error) {
	switch val {
	case metrics.MetricGroupTypeAny:
		return "any", nil
	case metrics.MetricGroupTypeSystem:
		return "system", nil
	case metrics.MetricGroupTypeSensors:
		return "sensors", nil
	case metrics.MetricGroupTypePower:
		return "power", nil
	case metrics.MetricGroupTypeIdracSel:
		return "sel", nil
	case metrics.MetricGroupTypeStorage:
		return "storage", nil
	case metrics.MetricGroupTypeMemory:
		return "memory", nil
	default:
		return "", fmt.Errorf("Unrecognized metric group type: '%s'", val)
	}
}

func getMetricGroupParam(req *http.Request) (metrics.MetricGroupType, error) {
	metric := req.URL.Query().Get("metric_group")

	switch metric {
	case "any":
		return metrics.MetricGroupTypeAny, nil
	case "":
		return metrics.MetricGroupTypeAny, nil
	case "system":
		return metrics.MetricGroupTypeSystem, nil
	case "sensors":
		return metrics.MetricGroupTypeSensors, nil
	case "power":
		return metrics.MetricGroupTypePower, nil
	case "sel":
		return metrics.MetricGroupTypeIdracSel, nil
	case "storage":
		return metrics.MetricGroupTypeStorage, nil
	case "memory":
		return metrics.MetricGroupTypeMemory, nil
	default:
		return metrics.MetricGroupTypeAny, fmt.Errorf("Unrecognized value for query parameter 'metric': '%s'", metric)
	}
}

func HealthHandler(rsp http.ResponseWriter, req *http.Request) {
	// just return a simple 200 for now
}

func ResetHandler(rsp http.ResponseWriter, req *http.Request) {
	target, err := getTargetParam(req)
	if err != nil {
		http.Error(rsp, err.Error(), http.StatusBadRequest)
		return
	}

	metricGroup, err := getMetricGroupParam(req)
	if err != nil {
		http.Error(rsp, err.Error(), http.StatusBadRequest)
		return
	}
	
	logging.Debugf("Handling reset-request from %s for host %s", req.Host, target)

	collector.Reset(target, metricGroup)
}


func MetricsHandler(rsp http.ResponseWriter, req *http.Request) {
	target, err := getTargetParam(req)
	if err != nil {
		http.Error(rsp, err.Error(), http.StatusBadRequest)
		return
	}

	metricGroup, err := getMetricGroupParam(req)
	if err != nil {
		http.Error(rsp, err.Error(), http.StatusBadRequest)
		return
	}

	handleRequestLogSuffix := ""
	if metricGroup == metrics.MetricGroupTypeAny {
		handleRequestLogSuffix = " (all metric groups)"
	} else {
		metricGroupName, err := getMetricGroupName(metricGroup)
		if err != nil {
			logging.Errorf(err, "Error finding metric group name")

			handleRequestLogSuffix = ""
		} else {
			handleRequestLogSuffix = fmt.Sprintf(" (%s)", metricGroupName)
		}
	}
	logging.Debugf("Handling request from %s for host %s%s", req.Host, target, handleRequestLogSuffix)

	c, err := collector.GetCollector(target, metricGroup)
	if err != nil {
		errorMsg := fmt.Sprintf("Error instantiating metrics collector for host %s", target)
		logging.Error(err, errorMsg)
		http.Error(rsp, errorMsg, http.StatusInternalServerError)
		return
	}

	logging.Debugf("Collecting metrics for host %s", target)

	metrics, err := c.Gather()
	if err != nil {
		errorMsg := fmt.Sprintf("Error collecting metrics for host %s", target)
		logging.Error(err, errorMsg)
		http.Error(rsp, errorMsg, http.StatusInternalServerError)
		return
	}

	logging.Debugf("Metrics for host %s collected", target)

	header := rsp.Header()
	header.Set(contentTypeHeader, "text/plain")

	// Code inspired by the official Prometheus metrics http handler
	w := io.Writer(rsp)
	if gzipAccepted(req.Header) {
		header.Set(contentEncodingHeader, "gzip")
		gz := gzipPool.Get().(*gzip.Writer)
		defer gzipPool.Put(gz)

		gz.Reset(w)
		defer gz.Close()

		w = gz
	}

	fmt.Fprint(w, metrics)
}

// gzipAccepted returns whether the client will accept gzip-encoded content.
func gzipAccepted(header http.Header) bool {
	a := header.Get(acceptEncodingHeader)
	parts := strings.Split(a, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "gzip" || strings.HasPrefix(part, "gzip;") {
			return true
		}
	}
	return false
}
