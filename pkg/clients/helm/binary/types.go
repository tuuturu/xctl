package binary

import "github.com/spf13/afero"

const logFeature = "helm/binary"

type externalBinaryHelm struct {
	fs             *afero.Afero
	kubeConfigPath string
	binaryPath     string
}
