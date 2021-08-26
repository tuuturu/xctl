package binary

import "github.com/spf13/afero"

type externalBinaryHelm struct {
	fs             *afero.Afero
	kubeConfigPath string
	binaryPath     string
}
