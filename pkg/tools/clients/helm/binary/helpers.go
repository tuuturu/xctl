package binary

import (
	"fmt"
	"path"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

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

func generateTempFile(fs *afero.Afero, name string, content []byte) (string, error) {
	tmpDir, err := fs.TempDir("/tmp", "xctl")
	if err != nil {
		return "", fmt.Errorf("creating temp dir for plugin values: %w", err)
	}

	tmpValuesPath := path.Join(tmpDir, name)

	err = fs.WriteFile(tmpValuesPath, content, 0o600)
	if err != nil {
		return "", fmt.Errorf("creating temporary values file: %w", err)
	}

	return tmpValuesPath, nil
}

type generateInstallArgsOpts struct {
	KubeConfigPath string
	Fs             *afero.Afero
	Plugin         v1alpha1.Plugin
}

func generateInstallArgs(opts generateInstallArgsOpts) ([]string, error) {
	argsMap := map[string]string{
		"namespace":  opts.Plugin.Metadata.Namespace,
		"kubeconfig": opts.KubeConfigPath,
		"atomic":     "true",
		"wait":       "true",
	}

	if opts.Plugin.Spec.Helm.Values != "" {
		valuesPath, err := generateTempFile(
			opts.Fs,
			fmt.Sprintf("%s-values.yaml", opts.Plugin.Metadata.Name),
			[]byte(opts.Plugin.Spec.Helm.Values),
		)
		if err != nil {
			return []string{}, fmt.Errorf("generating temporary values file: %w", err)
		}

		argsMap["values"] = valuesPath
	}

	if opts.Plugin.Spec.Helm.Version != "" {
		argsMap["version"] = opts.Plugin.Spec.Helm.Version
	}

	chartWithRepository := fmt.Sprintf("%s/%s",
		opts.Plugin.Spec.Helm.Repository.Name,
		opts.Plugin.Spec.Helm.Chart,
	)

	args := append(mapAsArgs(argsMap), "install", opts.Plugin.Metadata.Name, chartWithRepository)

	return args, nil
}

func mapAsArgs(m map[string]string) []string {
	result := make([]string, len(m))
	index := 0

	for key, value := range m {
		result[index] = fmt.Sprintf("--%s=%s", key, value)

		index++
	}

	return result
}
