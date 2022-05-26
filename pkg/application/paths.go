package application

import (
	"path"

	"github.com/deifyed/xctl/pkg/config"
)

func environmentDir(absoluteRepositoryRootDirectory string, environmentName string) string {
	return path.Join(
		absoluteRepositoryRootDirectory,
		config.DefaultInfrastructureDir,
		environmentName,
	)
}

func applicationsDir(absoluteRepositoryRootDirectory string, appName string) string {
	return path.Join(
		absoluteRepositoryRootDirectory,
		config.DefaultInfrastructureDir,
		config.DefaultApplicationsDir,
		appName,
	)
}
