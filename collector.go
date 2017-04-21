package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type collector struct {
	exporter     *Exporter
	driveshaft   *driveshaft
	up           *prometheus.Desc
	threadsGauge *prometheus.Desc
}

const metricsNamespace = "driveshaft"

// based on https://github.com/hnlq715/nginx-vts-exporter/
func newFuncMetric(metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(metricsNamespace, "", metricName),
		docString, labels, nil,
	)
}

func (e *Exporter) newCollector(g *driveshaft) *collector {
	return &collector{
		exporter:     e,
		driveshaft:   g,
		up:           newFuncMetric("up", "is driveshaft up", nil),
		threadsGauge: newFuncMetric("threads_count", "count of threads", []string{"function", "state"}),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.threadsGauge
}

func (c *collector) collectThreads(ch chan<- prometheus.Metric) {
	t, err := c.driveshaft.getThreads()
	if err != nil {
		c.exporter.logger.Error("failed to get driveshaft status", zap.Error(err))
		ch <- prometheus.MustNewConstMetric(
			c.up,
			prometheus.GaugeValue,
			float64(0),
		)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		c.up,
		prometheus.GaugeValue,
		float64(1),
	)

	for _, v := range t {
		ch <- prometheus.MustNewConstMetric(
			c.threadsGauge,
			prometheus.GaugeValue,
			float64(v.count),
			v.function, v.state)

	}
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.collectThreads(ch)
}
