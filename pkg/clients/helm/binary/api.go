package binary

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path"
	"strings"

	"github.com/deifyed/xctl/pkg/tools/logging"
	"github.com/sirupsen/logrus"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/clients/helm"
	"github.com/spf13/afero"
)

func (e externalBinaryHelm) Install(plugin v1alpha1.Plugin) error {
	log := logging.CreateEntry(logrus.StandardLogger(), logFeature, "install")

	tmpDir, err := e.fs.TempDir("/tmp", "xctl")
	if err != nil {
		return fmt.Errorf("creating temp dir for plugin values: %w", err)
	}

	tmpValuesPath := path.Join(tmpDir, fmt.Sprintf("%s-values.yaml", plugin.Metadata.Name))

	err = e.fs.WriteFile(tmpValuesPath, []byte(plugin.Spec.Helm.Values), 0o600)
	if err != nil {
		return fmt.Errorf("creating temporary values file: %w", err)
	}

	cmd := exec.Command(e.binaryPath,
		fmt.Sprintf("--namespace=%s", plugin.Metadata.Namespace),
		fmt.Sprintf("--kubeconfig=%s", e.kubeConfigPath),
		"install",
		"--atomic",
		plugin.Metadata.Name,
		plugin.Spec.Helm.Chart,
		fmt.Sprintf("--values=%s", tmpValuesPath),
	)

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err = cmd.Run()
	if err != nil {
		log.WithFields(logrus.Fields{
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		}).Debug("executing command")

		return fmt.Errorf("running Helm install on %s: %w", plugin.Metadata.Name, err)
	}

	return nil
}

func (e externalBinaryHelm) Delete(plugin v1alpha1.Plugin) error {
	log := logging.CreateEntry(logrus.StandardLogger(), logFeature, "delete")

	cmd := exec.Command(e.binaryPath,
		fmt.Sprintf("--namespace=%s", plugin.Metadata.Namespace),
		fmt.Sprintf("--kubeconfig=%s", e.kubeConfigPath),
		"uninstall",
		plugin.Metadata.Name,
	)

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		log.WithFields(logrus.Fields{
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		}).Debug("executing command")

		return fmt.Errorf("running Helm uninstall on %s: %w", plugin.Metadata.Name, err)
	}

	return nil
}

func (e externalBinaryHelm) Exists(plugin v1alpha1.Plugin) (bool, error) {
	log := logging.CreateEntry(logrus.StandardLogger(), logFeature, "delete")

	cmd := exec.Command(e.binaryPath,
		fmt.Sprintf("--namespace=%s", plugin.Metadata.Namespace),
		fmt.Sprintf("--kubeconfig=%s", e.kubeConfigPath),
		"get",
		"manifest",
		plugin.Metadata.Name,
	)

	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr
	cmd.Stdout = io.Discard

	err := cmd.Run()
	if err != nil {
		if strings.HasPrefix(stderr.String(), "Error: release: not found") {
			return false, nil
		}

		log.WithFields(logrus.Fields{
			"stderr": stderr.String(),
		}).Debug("executing command")

		return false, fmt.Errorf("running Helm get: %s", stderr.String())
	}

	return true, nil
}

func New(fs *afero.Afero, kubeConfigPath string) (helm.Client, error) {
	binaryPath, err := getHelmPath(fs)
	if err != nil {
		return nil, fmt.Errorf("acquiring Helm path: %w", err)
	}

	return &externalBinaryHelm{
		kubeConfigPath: kubeConfigPath,
		binaryPath:     binaryPath,
		fs:             fs,
	}, nil
}
