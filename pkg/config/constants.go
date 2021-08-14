package config

import (
	"time"
)

const ( // Controller / reconciliation
	// DefaultMaxReconciliationRequeues defines the maximum amount of times to requeue a reconciler
	DefaultMaxReconciliationRequeues = 3
	// DefaultReconciliationLoopDelayDuration defines the amount of time to wait between each reconciliation
	DefaultReconciliationLoopDelayDuration = 1 * time.Second
)

const ( // Cluster
	// DefaultClusterNodeAmount defines the default amount of nodes to provision for a cluster
	DefaultClusterNodeAmount = 2
)

const ( // Domain
	// DefaultSubdomainTTLSeconds defines the default time to live for a new domain
	DefaultSubdomainTTLSeconds = 300
)
