package blocklist

import (
	"github.com/coredns/coredns/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	pluginName = "blocklist"
)

var passedQueriesCount = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: pluginName,
	Name:      "queries_passed",
	Help:      "number of queries passed",
})

var blockedQueriesCount = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Subsystem: pluginName,
	Name:      "queries_blocked",
	Help:      "number of queries blocked",
})
