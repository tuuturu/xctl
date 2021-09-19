package binary

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/helm"
	"github.com/spf13/afero"
)

func (e externalBinaryHelm) Install(plugin v1alpha1.Plugin) error {
	tmpDir, err := e.fs.TempDir("/tmp", "xctl")
	if err != nil {
		return fmt.Errorf("creating temp dir for plugin values: %w", err)
	}

	tmpValuesPath := path.Join(tmpDir, fmt.Sprintf("%s-values.yaml", plugin.Metadata.Name))

	err = e.fs.WriteFile(tmpValuesPath, []byte(plugin.Spec.Values), 0o600)
	if err != nil {
		return fmt.Errorf("creating temporary values file: %w", err)
	}

	cmd := exec.Command(e.binaryPath,
		"install",
		plugin.Metadata.Name,
		plugin.Spec.HelmChart,
		fmt.Sprintf("--kubeconfig=%s", e.kubeConfigPath),
		fmt.Sprintf("--values=%s", tmpValuesPath),
		"--atomic",
		"--debug",
	)

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("running Helm install on %s: %w", plugin.Metadata.Name, err)
	}

	return nil
}

func (e externalBinaryHelm) Delete(plugin v1alpha1.Plugin) error {
	cmd := exec.Command(e.binaryPath,
		"uninstall",
		plugin.Metadata.Name,
		fmt.Sprintf("--kubeconfig=%s", e.kubeConfigPath),
	)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("running Helm uninstall on %s: %w", plugin.Metadata.Name, err)
	}

	return nil
}

func (e externalBinaryHelm) Exists(plugin v1alpha1.Plugin) (bool, error) {
	println(fmt.Sprintf("with %s", e.kubeConfigPath))

	cmd := exec.Command(e.binaryPath,
		fmt.Sprintf("--kubeconfig=%s", e.kubeConfigPath),
		"get",
		"manifest",
		plugin.Metadata.Name,
		"--debug",
	)

	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if strings.HasPrefix(stderr.String(), "Error: release: not found") {
			return false, nil
		}

		return false, fmt.Errorf("running Helm get: %s", stderr.String())
	}

	return true, nil
}

func NewExternalBinaryHelm(fs *afero.Afero) helm.Client {
	binaryPath := "/usr/bin/helm"

	return &externalBinaryHelm{
		kubeConfigPath: config.GetAbsoluteKubeconfigPath(),
		binaryPath:     binaryPath,
		fs:             fs,
	}
}
