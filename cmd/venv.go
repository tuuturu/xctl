package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/deifyed/xctl/cmd/helpers"
	"github.com/deifyed/xctl/pkg/cloud"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

	"github.com/deifyed/xctl/pkg/cloud/linode"

	"github.com/deifyed/xctl/pkg/tools/venv"
	"github.com/deifyed/xctl/pkg/tools/venv/shells/zsh"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type venvOpts struct {
	fs                     *afero.Afero
	clusterDeclarationPath string
	clusterManifest        v1alpha1.Cluster
}

var (
	venvCmdOpts = venvOpts{ //nolint:gochecknoglobals
		fs: &afero.Afero{Fs: afero.NewOsFs()},
	}
	venvCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "venv",
		Short: "activates a virtual environment",
		PreRunE: helpers.ClusterManifestIniter(
			venvCmdOpts.fs,
			&venvCmdOpts.clusterDeclarationPath,
			&venvCmdOpts.clusterManifest,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			workDir, err := venvCmdOpts.fs.TempDir("/tmp", "xctl")
			if err != nil {
				return fmt.Errorf("creating temporary directory: %w", err)
			}

			provider := linode.NewLinodeProvider()

			err = provider.Authenticate()
			if err != nil {
				return fmt.Errorf("authenticating with cloud provider: %w", err)
			}

			err = setupKubeconfig(ctx, provider, venvCmdOpts.fs, venvCmdOpts.clusterManifest.Metadata.Name, workDir)
			if err != nil {
				return fmt.Errorf("setting up kubeconfig: %w", err)
			}

			env := venv.MergeVariables(os.Environ(), []string{
				fmt.Sprintf("KUBECONFIG=%s", path.Join(workDir, "config.yaml")),
			})

			shellCmd, err := acquireShellCommand(venvCmdOpts.fs, workDir, env)
			if err != nil {
				return fmt.Errorf("acquiring shell: %w", err)
			}

			fmt.Printf("Launching virtual environment\n")

			err = shellCmd.Run()
			if err != nil {
				return fmt.Errorf("running shell: %w", err)
			}

			err = teardown(venvCmdOpts.fs, workDir)
			if err != nil {
				return fmt.Errorf("tearing down virtual environment: %w", err)
			}

			fmt.Printf("Successfully exited virtual environment\n")

			return nil
		},
	}
)

//nolint:gochecknoinits
func init() {
	flags := venvCmd.Flags()

	flags.StringVarP(
		&venvCmdOpts.clusterDeclarationPath,
		"cluster-declaration",
		"c",
		"-",
		"cluster declaration representing context of the virtual environment",
	)

	rootCmd.AddCommand(venvCmd)
}

func acquireShellCommand(fs *afero.Afero, workDir string, env []string) (*exec.Cmd, error) {
	shellPath, err := venv.GetCurrentShell(fs)
	if err != nil {
		return nil, fmt.Errorf("getting current shell path: %w", err)
	}

	var shell venv.Shell

	switch {
	case strings.HasSuffix(shellPath, "zsh"):
		shell = zsh.NewZshShell(venvCmdOpts.fs, workDir, shellPath)
	default:
		return nil, fmt.Errorf("handling shell %s. Not a known shell", shellPath)
	}

	shellCommand, err := shell.Command(env)
	if err != nil {
		return nil, fmt.Errorf("acquiring shell command: %w", err)
	}

	return shellCommand, nil
}

func setupKubeconfig(ctx context.Context, provider cloud.Provider, fs *afero.Afero, clusterName, workDir string) error {
	kubeConfig, err := provider.GetKubeConfig(ctx, clusterName)
	if err != nil {
		return fmt.Errorf("acquiring kubeconfig: %w", err)
	}

	decodedConfig, err := base64.StdEncoding.DecodeString(string(kubeConfig))
	if err != nil {
		return fmt.Errorf("decoding kubeconfig: %w", err)
	}

	err = fs.WriteFile(path.Join(workDir, "config.yaml"), decodedConfig, 0o644)
	if err != nil {
		return fmt.Errorf("writing kubeconfig: %w", err)
	}

	return nil
}

func teardown(fs *afero.Afero, workDir string) error {
	return fs.RemoveAll(workDir)
}
