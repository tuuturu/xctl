package binary

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/deifyed/xctl/pkg/clients/vault"

	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/tools/binaries"
	"github.com/spf13/afero"
)

const version = "1.8.4"

var (
	keyRe   = regexp.MustCompile(`Unseal\sKey\s\d:\s(?P<key>.+)`)
	tokenRe = regexp.MustCompile(`Token:\s(?P<token>.+)`)
)

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
	result := vault.InitializationResponse{Keys: make([]string, 5)}
	buffer := bytes.Buffer{}

	_, err := io.Copy(&buffer, reader)
	if err != nil {
		return vault.InitializationResponse{}, fmt.Errorf("preparing buffer: %w", err)
	}

	tokenMatch := tokenRe.FindStringSubmatch(buffer.String())
	result.Token = tokenMatch[tokenRe.SubexpIndex("token")]

	keyMatch := keyRe.FindAllStringSubmatch(buffer.String(), 5)
	for index, match := range keyMatch {
		result.Keys[index] = match[keyRe.SubexpIndex("key")]
	}

	return result, nil
}
