package piholedb

import (
	"github.com/coredns/coredns/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	subsystem = "piholedb"
)

var gravityLookupDuration = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: plugin.Namespace,
	Subsystem: subsystem,
	Name:      "lookup_duration",
	Help:      "how long it take to find record in gravity.db",
})
