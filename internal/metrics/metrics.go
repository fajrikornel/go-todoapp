package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var registry *prometheus.Registry
var metricsMap map[string]prometheus.Collector

type MetricType interface {
	GetCollectorInitializer(name, help string, labels []string) func() prometheus.Collector
}

type mappedMetric struct {
	MetricName               string
	Help                     string
	PrometheusLabels         prometheus.Labels
	CollectorInitializerFunc func() prometheus.Collector
}

func init() {
	registry = prometheus.DefaultRegisterer.(*prometheus.Registry)
	metricsMap = make(map[string]prometheus.Collector)
}

func GetRegisterer() prometheus.Registerer {
	return registry
}

func GetGatherer() prometheus.Gatherer {
	return registry
}
