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

var (
	// ErrNotFound indicates that something is missing
	ErrNotFound = errors.New("not found")
	// ErrTimeout indicates that something exceeded its deadline
	ErrTimeout = errors.New("timeout")
)
