package binary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/sirupsen/logrus"

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
		URL:         fmt.Sprintf("https://releases.hashicorp.com/vault/%s/vault_%s_linux_amd64.zip", version, version),
		UnpackingFn: []binaries.UnpackingFn{binaries.GenerateZipUnpacker("vault")},
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

func (c *client) runVaultCommand(args ...string) (io.Reader, error) {
	cmd := exec.Command(c.vaultPath, args...) //nolint:gosec

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Env = c.envAsArray()
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		}).Debug("executing command")

		err = fmt.Errorf("%s: %w", stderr.String(), err)

		return nil, errorHandler(err, fmt.Errorf("executing command: %w", err))
	}

	return &stdout, nil
}
