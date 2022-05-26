package manifests

const kustomizationFilename = "kustomization.yaml"

type patchTarget struct {
	Kind string `json:"kind"`
}

type patch struct {
	Path   string      `json:"path"`
	Target patchTarget `json:"target"`
}

type kustomize struct {
	Resources []string `json:"resources"`
	Patches   []patch  `json:"patches,omitempty"`
}
