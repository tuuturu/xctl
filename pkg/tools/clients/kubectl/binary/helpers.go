package binary

import (
	"fmt"

	"github.com/deifyed/xctl/pkg/tools/binaries"

	"github.com/deifyed/xctl/pkg/config"

	"github.com/spf13/afero"
)

const (
	version = "1.23.0"
	hash    = "3f0398d4c8a5ff633e09abd0764ed3b9091fafbe3044970108794b02731c72d6"
)

func getKubectlPath(fs *afero.Afero) (string, error) {
	binariesDir, err := config.GetAbsoluteBinariesDir()
	if err != nil {
		return "", fmt.Errorf("acquiring binaries directory: %w", err)
	}

	path, err := binaries.Download(binaries.DownloadOpts{
		Name:        "kubectl",
		Version:     version,
		Fs:          fs,
		BinariesDir: binariesDir,
		URL:         fmt.Sprintf("https://dl.k8s.io/release/v%s/bin/linux/amd64/kubectl", version),
		Hash:        hash,
	})
	if err != nil {
		return "", fmt.Errorf("downloading and checking checksum: %w", err)
	}

	return path, nil
}

func (k kubectlBinaryClient) envAsArray() []string {
	env := make([]string, len(k.env))
	index := 0

	for key, value := range k.env {
		env[index] = fmt.Sprintf("%s=%s", key, value)

		index++
	}

	return env
}
