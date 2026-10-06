package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	ConfigSyncsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "configsync_reconciliations_total",
			Help: "Total number of ConfigSync custom resources reconciled",
		},
	)
)

func init() {
	// Register custom metrics with controller-runtime's registry
	metrics.Registry.MustRegister(ConfigSyncsTotal)
}