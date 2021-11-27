package config

import (
	"fmt"
	"os"
	"path"
)

func GetAbsoluteXCTLDir() (string, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("acquiring home directory")
	}

	return path.Join(userDir, fmt.Sprintf(".%s", ApplicationName)), nil
}

// GetAbsoluteXCTLClusterDir returns the relevant cluster directory for cluserName in the xctl directory
func GetAbsoluteXCTLClusterDir(clusterName string) (string, error) {
	xctlDir, err := GetAbsoluteXCTLDir()
	if err != nil {
		return "", err
	}

	return path.Join(xctlDir, DefaultClustersDir, clusterName), nil
}

// GetAbsoluteKubeconfigPath knows where the cluster Kubeconfig file is
func GetAbsoluteKubeconfigPath(clusterName string) (string, error) {
	clusterDir, err := GetAbsoluteXCTLClusterDir(clusterName)
	if err != nil {
		return "", err
	}

	return path.Join(clusterDir, DefaultKubeconfigFilename), nil
}

func GetAbsoluteBinariesDir() (string, error) {
	xctlDir, err := GetAbsoluteXCTLDir()
	if err != nil {
		return "", fmt.Errorf("acquiring xctl directory: %w", err)
	}

	return path.Join(xctlDir, DefaultBinariesDir), nil
}

func IsDebugMode() bool {
	debugMode := os.Getenv("XCTL_DEBUG")

	return debugMode == "true"
}
