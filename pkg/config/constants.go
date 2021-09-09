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

const ( // Internal filesystem directories
	// DefaultAbsoluteRootPath defines the root of the internal filesystem
	DefaultAbsoluteRootPath = "/"
	// DefaultConfigDirName defines the default location for config files
	DefaultConfigDirName = "config"
	// DefaultKubeconfigFilename defines the name of the kubeconfig in the internal FS
	DefaultKubeconfigFilename = "kubeconfig.yaml"
	// DefaultManifestDir defines the name of the manifest directory
	DefaultManifestDir = "manifests"
	// DefaultClusterManifestFilename defines the name of the cluster manifest file
	DefaultClusterManifestFilename = "cluster.yaml"
)
