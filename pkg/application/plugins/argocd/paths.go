package argocd

import "path"

func argoCDApplicationPath(absoluteEnvironmentDirectory string) string {
	return path.Join(absoluteEnvironmentDirectory, "argocd", "applications")
}
