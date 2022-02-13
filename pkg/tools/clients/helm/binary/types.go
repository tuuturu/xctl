package binary

import "github.com/spf13/afero"

const logFeature = "helm/binary"

type externalBinaryHelm struct {
	fs             *afero.Afero
	kubeConfigPath string
	binaryPath     string
}

type commandLogFields struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}
