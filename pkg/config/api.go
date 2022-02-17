package config

import (
	"fmt"
	"os"
	"path"
)

// GetAbsoluteXCTLDir returns the main xctl configuration directory, usually lives in the users home directory
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

	return path.Join(xctlDir, DefaultEnvironmentsDir, clusterName), nil
}

// GetAbsoluteKubeconfigPath knows where the cluster Kubeconfig file is
func GetAbsoluteKubeconfigPath(clusterName string) (string, error) {
	clusterDir, err := GetAbsoluteXCTLClusterDir(clusterName)
	if err != nil {
		return "", err
	}

	return path.Join(clusterDir, DefaultKubeconfigFilename), nil
}

// GetAbsoluteBinariesDir returns the absolute path of the directory containing downloaded binaries
func GetAbsoluteBinariesDir() (string, error) {
	xctlDir, err := GetAbsoluteXCTLDir()
	if err != nil {
		return "", fmt.Errorf("acquiring xctl directory: %w", err)
	}

	return path.Join(xctlDir, DefaultBinariesDir), nil
}

// IsDebugMode returns a boolean representing if xctl is running is debug mode or not
func IsDebugMode() bool {
	debugMode := os.Getenv("XCTL_DEBUG")

	return debugMode == "true"
}
