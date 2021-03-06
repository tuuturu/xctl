package binary

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/deifyed/xctl/pkg/tools/logging"
)

type runCommandOpts struct {
	Log       logging.Logger
	Namespace string
	Stdin     io.Reader
	Args      []string
}

func (k kubectlBinaryClient) runCommand(opts runCommandOpts) (io.Reader, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("acquiring user home directory: %w", err)
	}

	opts.Args = append([]string{"--cache-dir", path.Join(homeDir, ".kube", "cache")}, opts.Args...)

	if opts.Namespace != "" {
		opts.Args = append([]string{"--namespace", opts.Namespace}, opts.Args...)
	}

	cmd := exec.Command(k.kubectlPath, opts.Args...)

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Env = k.envAsArray()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if opts.Stdin != nil {
		cmd.Stdin = opts.Stdin
	}

	err = cmd.Run()
	if err != nil {
		opts.Log.Debug("executing command", commandLogFields{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
		})

		err = fmt.Errorf("%s: %w", stderr.String(), err)

		return nil, errorHandler(err, fmt.Errorf("executing pod command: %s", err))
	}

	return &stdout, nil
}
