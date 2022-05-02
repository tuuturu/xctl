package config

import (
	"time"
)

// ApplicationName refers to the xctl project name
const ApplicationName = "xctl"

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

const ( // Cluster namespaces
	// DefaultOperationsNamespace defines the name of the operations namespace. This namespace contains operations
	// related resources. I.e.: ArgoCD, Dex, etc.
	DefaultOperationsNamespace = "operations"
	// DefaultMonitoringNamespace defines the name of the monitoring namespace. This namespace contains monitoring
	// related resources. I.e.: Grafana, Prometheus, Loki, Promtail, Tempo, etc.
	DefaultMonitoringNamespace = "monitoring"
)

const ( // Domain
	// DefaultSubdomainTTLSeconds defines the default time to live for a new domain
	DefaultSubdomainTTLSeconds = 300
)

const ( // Internal filesystem directories
	// DefaultEnvironmentsDir defines the folder containing environment directories
	DefaultEnvironmentsDir = "environments"
	// DefaultKubeconfigFilename defines the name of the kubeconfig in the internal FS
	DefaultKubeconfigFilename = "kubeconfig.yaml"
	// DefaultBinariesDir defines the name of the directory containing binaries
	DefaultBinariesDir = "binaries"
)
