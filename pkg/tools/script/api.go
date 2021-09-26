package script

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"

	"github.com/deifyed/xctl/pkg/config"
	"github.com/spf13/afero"
)

func NewScriptRunner(fs *afero.Afero, env map[string]string) *Runner {
	return &Runner{
		fs:  fs,
		env: env,
	}
}

func (receiver Runner) Execute(script []byte) (int, error) {
	workDir, err := receiver.fs.TempDir("/tmp", config.ApplicationName)
	if err != nil {
		return -1, fmt.Errorf("creating temporary directory: %w", err)
	}

	scriptPath := path.Join(workDir, "script.sh")

	err = receiver.fs.WriteFile(scriptPath, script, 0o700)
	if err != nil {
		return -1, fmt.Errorf("creating script file: %w", err)
	}

	cmd := exec.Command("sh", scriptPath)

	stderr := bytes.Buffer{}

	cmd.Env = mapAsSlice(receiver.env)
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return -1, fmt.Errorf("running script: %w", fmt.Errorf("%s: %w", stderr.String(), err))
	}

	return 0, nil
}
