package binary

import (
	"bytes"
	"fmt"
	"os/exec"

	"sigs.k8s.io/yaml"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/spf13/afero"

	"github.com/deifyed/xctl/pkg/clients/kubectl"
)

func (k kubectlBinaryClient) Apply(opts kubectl.ApplyOpts) error {
	log := logging.GetLogger(logFeature, "apply")

	raw, err := yaml.Marshal(opts.Manifest)
	if err != nil {
		return fmt.Errorf("marshalling manifest: %w", err)
	}

	cmd := exec.Command(k.kubectlPath, "apply", "-f", "-") //nolint:gosec

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Env = k.envAsArray()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = bytes.NewReader(raw)

	err = cmd.Run()
	if err != nil {
		log.Debug("executing command", commandLogFields{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
		})

		return fmt.Errorf("executing pod command: %s", err)
	}

	return nil
}

func New(fs *afero.Afero, kubeConfigPath string) (kubectl.Client, error) {
	kubectlPath, err := getKubectlPath(fs)
	if err != nil {
		return nil, fmt.Errorf("acquiring kubectl path: %w", err)
	}

	return &kubectlBinaryClient{
		kubectlPath: kubectlPath,
		env: map[string]string{
			kubeConfigPathKey: kubeConfigPath,
		},
	}, nil
}
