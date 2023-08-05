package metrics

import "github.com/prometheus/client_golang/prometheus"

var registry *prometheus.Registry
var metricsMap map[string]prometheus.Collector

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
