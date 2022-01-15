package binary

import (
	"fmt"
	"regexp"

	"github.com/deifyed/xctl/pkg/tools/binaries"

	"github.com/deifyed/xctl/pkg/config"
	"github.com/spf13/afero"
)

const (
	version = "3.7.0"
	hash    = "096e30f54c3ccdabe30a8093f8e128dba76bb67af697b85db6ed0453a2701bf9"
)

func getHelmPath(fs *afero.Afero) (string, error) {
	binariesDir, err := config.GetAbsoluteBinariesDir()
	if err != nil {
		return "", fmt.Errorf("acquiring binaries directory: %w", err)
	}

	helmPath, err := binaries.Download(binaries.DownloadOpts{
		Name:        "helm",
		Version:     version,
		Fs:          fs,
		BinariesDir: binariesDir,
		UnpackingFn: []binaries.UnpackingFn{binaries.GzipUnpacker, binaries.GenerateTarUnpacker("helm")},
		URL:         fmt.Sprintf("https://get.helm.sh/helm-v%s-linux-amd64.tar.gz", version),
		Hash:        hash,
	})
	if err != nil {
		return "", fmt.Errorf("downloading and checking checksum: %w", err)
	}

	return helmPath, nil
}

var reUnreachableErr = regexp.MustCompile(`.*EOF: Kubernetes cluster unreachable.*`)

func isUnreachable(err error) bool {
	return reUnreachableErr.MatchString(err.Error())
}

var reTimedOutErr = regexp.MustCompile(`.*connection timed out.*`)

func isConnectionTimedOut(err error) bool {
	return reTimedOutErr.MatchString(err.Error())
}
