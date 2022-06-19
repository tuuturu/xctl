package kustomize

const defaultKustomizationFilename = "kustomization.yaml"

type file struct {
	Resources []string `json:"resources"`
}
