package binary

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/deifyed/xctl/pkg/tools/clients/helm"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

	"github.com/deifyed/xctl/pkg/tools/logging"
)

func (e externalBinaryHelm) addRepository(repository v1alpha1.PluginSpecHelmRepository) error {
	log := logging.GetLogger(logFeature, "addRepository")

	cmd := exec.Command(e.binaryPath, "repo", "add", repository.Name, repository.URL)

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	log.Debugf("adding repository %s as %s", repository.URL, repository.Name)

	err := cmd.Run()
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
			return fmt.Errorf("running command: %w", err)
		}
	}

	return nil
}

func (e externalBinaryHelm) updateRepositories() error {
	log := logging.GetLogger(logFeature, "updateRepositories")

	cmd := exec.Command(e.binaryPath, "repo", "update")

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	log.Debug("updating repositories")

	err := cmd.Run()
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
			return fmt.Errorf("running command: %w", err)
		}
	}

	return nil
}
