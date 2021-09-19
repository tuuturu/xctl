package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/deifyed/xctl/pkg/apis/xctl"

	"github.com/deifyed/xctl/pkg/cloud/linode"

	"github.com/deifyed/xctl/cmd/preruns"
	"github.com/deifyed/xctl/pkg/cloud"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

	"github.com/deifyed/xctl/pkg/tools/venv"
	"github.com/deifyed/xctl/pkg/tools/venv/shells/zsh"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type venvOpts struct {
	io                     xctl.IOStreams
	fs                     *afero.Afero
	clusterDeclarationPath string
	clusterManifest        v1alpha1.Cluster
}

var (
	venvCmdOpts = venvOpts{
		io: xctl.IOStreams{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
		fs: &afero.Afero{Fs: afero.NewOsFs()},
	}
	venvCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "venv",
		Short: "activates a virtual environment enabling manipulation of the production environment",
		PreRunE: preruns.ClusterManifestIniter(preruns.ClusterManifestIniterOpts{
			Io:              venvCmdOpts.io,
			Fs:              venvCmdOpts.fs,
			ClusterManifest: &venvCmdOpts.clusterManifest,
			SourcePath:      &venvCmdOpts.clusterDeclarationPath,
		}),
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

			fmt.Fprintf(venvCmdOpts.io.Out, "Launching virtual environment\n")

			err = shellCmd.Run()
			if err != nil {
				return fmt.Errorf("running shell: %w", err)
			}

			err = teardown(venvCmdOpts.fs, workDir)
			if err != nil {
				return fmt.Errorf("tearing down virtual environment: %w", err)
			}

			fmt.Fprintf(venvCmdOpts.io.Out, "Successfully exited virtual environment\n")

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

// acquireShellCommand figures out what the preferred shell is and returns an exec.Cmd version of it
func acquireShellCommand(fs *afero.Afero, workDir string, env []string) (*exec.Cmd, error) {
	shellPath, err := venv.GetCurrentShell(fs)
	if err != nil {
		return nil, fmt.Errorf("getting current shell path: %w", err)
	}

	var shell venv.Shell

	switch {
	case strings.HasSuffix(shellPath, "zsh"):
		shell = zsh.NewZshShell(fs, workDir, shellPath)
	default:
		return nil, fmt.Errorf("handling shell %s. Not a known shell", shellPath)
	}

	shellCommand, err := shell.Command(venvCmdOpts.io, env)
	if err != nil {
		return nil, fmt.Errorf("acquiring shell command: %w", err)
	}

	return shellCommand, nil
}

// setupKubeconfig fetches a kubeconfig from provider and stores it b64 decoded in workdir as config.yaml
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

// teardown deletes the folder workdir and all of it's content
func teardown(fs *afero.Afero, workDir string) error {
	return fs.RemoveAll(workDir)
}
