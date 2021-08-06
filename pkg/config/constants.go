package config

import (
	"errors"
	"time"
)

const (
	// DefaultClusterNodeAmount defines the default amount of nodes to provision for a cluster
	DefaultClusterNodeAmount = 2
	// DefaultMaxReconciliationRequeues defines the maximum amount of times to requeue a reconciler
	DefaultMaxReconciliationRequeues = 3
	// DefaultReconciliationLoopDelayDuration defines the amount of time to wait between each reconciliation
	DefaultReconciliationLoopDelayDuration = 1 * time.Second
)

// ErrNotFound indicates that something is missing
var ErrNotFound = errors.New("not found")
