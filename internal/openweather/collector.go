package openweather

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
)

var Namespace = "weather"
var Subsystem = "openweather"

type Collector struct {
	temperature *prometheus.Desc
	humidity    *prometheus.Desc
}

func NewCollector() *Collector {
	return &Collector{
		temperature: prometheus.NewDesc(prometheus.BuildFQName(Namespace, Subsystem, "temperature"),
			"The current temperature, in degrees Celsius", []string{"location"}, nil),
		humidity: prometheus.NewDesc(prometheus.BuildFQName(Namespace, Subsystem, "humidity"),
			"The current relative humidity percentage", []string{"location"}, nil),
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.temperature
	ch <- c.humidity
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	slog.Debug("metrics collected")
	ch <- prometheus.MustNewConstMetric(c.temperature, prometheus.GaugeValue, 28.5, "arvada")
	ch <- prometheus.MustNewConstMetric(c.humidity, prometheus.GaugeValue, 45.78, "arvada")
}
