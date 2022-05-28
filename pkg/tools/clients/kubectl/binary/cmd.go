package binary

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/deifyed/xctl/pkg/tools/logging"
)

func (k kubectlBinaryClient) runCommand(log logging.Logger, args ...string) (io.Reader, error) {
	cmd := exec.Command(k.kubectlPath, args...)

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Env = k.envAsArray()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Debug("executing command", commandLogFields{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
		})

		err = fmt.Errorf("%s: %w", stderr.String(), err)

		return nil, errorHandler(err, fmt.Errorf("executing pod command: %s", err))
	}

	return &stdout, nil
}
