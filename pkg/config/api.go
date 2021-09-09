package config

import "path"

// GetAbsoluteKubeconfigPath knows where the cluster Kubeconfig file is
func GetAbsoluteKubeconfigPath() string {
	return path.Join(DefaultAbsoluteRootPath, DefaultConfigDirName, DefaultKubeconfigFilename)
}

// GetAbsoluteInternalClusterManifestPath knows where the cluster manifest is
func GetAbsoluteInternalClusterManifestPath() string {
	return path.Join(DefaultAbsoluteRootPath, DefaultManifestDir, DefaultClusterManifestFilename)
}
