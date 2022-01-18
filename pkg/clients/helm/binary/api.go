package binary

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/clients/helm"
	"github.com/deifyed/xctl/pkg/tools/logging"
	"github.com/spf13/afero"
)

func (e externalBinaryHelm) Install(plugin v1alpha1.Plugin) error {
	log := logging.GetLogger(logFeature, "install")

	err := e.addRepository(plugin.Spec.Helm.Repository)
	if err != nil {
		return fmt.Errorf("adding repository: %w", err)
	}

	err = e.updateRepositories()
	if err != nil {
		return fmt.Errorf("updating repositories: %w", err)
	}

	args, err := generateInstallArgs(generateInstallArgsOpts{
		KubeConfigPath: e.kubeConfigPath,
		Fs:             e.fs,
		Plugin:         plugin,
	})
	if err != nil {
		return fmt.Errorf("preparing installation arguments: %w", err)
	}

	cmd := exec.Command(e.binaryPath, args...)

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err = cmd.Run()
	if err != nil {
		log.Debug("executing command", commandLogFields{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
		})

		switch {
		case isUnreachable(err):
			return helm.ErrUnreachable
		case isConnectionTimedOut(err):
			return helm.ErrTimeout
		default:
			return fmt.Errorf("running Helm install on %s: %w", plugin.Metadata.Name, err)
		}
	}

	return nil
}

func (e externalBinaryHelm) Delete(plugin v1alpha1.Plugin) error {
	log := logging.GetLogger(logFeature, "delete")

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
		log.Debug("executing command", commandLogFields{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
		})

		return fmt.Errorf("running Helm uninstall on %s: %w", plugin.Metadata.Name, err)
	}

	return nil
}

func (e externalBinaryHelm) Exists(plugin v1alpha1.Plugin) (bool, error) {
	log := logging.GetLogger(logFeature, "exists")

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

		log.Debug("executing command", commandLogFields{
			Stderr: stderr.String(),
		})

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
