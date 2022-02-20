package binary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/deifyed/xctl/pkg/tools/clients/vault"

	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/tools/binaries"
	"github.com/spf13/afero"
)

const version = "1.9.3"

func getVaultPath(fs *afero.Afero) (string, error) {
	binariesDir, err := config.GetAbsoluteBinariesDir()
	if err != nil {
		return "", fmt.Errorf("acquiring binaries directory: %w", err)
	}

	path, err := binaries.Download(binaries.DownloadOpts{
		Name:        "vault",
		Version:     version,
		Fs:          fs,
		BinariesDir: binariesDir,
		URL:         fmt.Sprintf("https://github.com/hashicorp/vault/archive/refs/tags/v%s.tar.gz", version),
		UnpackingFn: []binaries.UnpackingFn{binaries.GzipUnpacker, binaries.GenerateTarUnpacker("vault")},
	})
	if err != nil {
		return "", fmt.Errorf("downloading and checking checksum: %w", err)
	}

	return path, nil
}

func parseInitializationResponse(reader io.Reader) (vault.InitializationResponse, error) {
	response := vault.InitializationResponse{}
	buffer := bytes.Buffer{}

	_, err := io.Copy(&buffer, reader)
	if err != nil {
		return vault.InitializationResponse{}, fmt.Errorf("preparing buffer: %w", err)
	}

	err = json.Unmarshal(buffer.Bytes(), &response)
	if err != nil {
		return vault.InitializationResponse{}, fmt.Errorf("unmarshalling response: %w", err)
	}

	return response, nil
}

func (c *client) envAsArray() []string {
	env := make([]string, len(c.env))
	index := 0

	for key, value := range c.env {
		env[index] = fmt.Sprintf("%s=%s", key, value)

		index++
	}

	return env
}

func attributesAsArray(m map[string]string) []string {
	result := make([]string, len(m))
	index := 0

	for key, value := range m {
		result[index] = fmt.Sprintf("%s=%s", key, value)

		index++
	}

	return result
}
