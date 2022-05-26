package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/deifyed/xctl/pkg/tools/secrets/keyring"

	"github.com/deifyed/xctl/cmd/helpers"

	"github.com/deifyed/xctl/pkg/tools/i18n"

	"github.com/deifyed/xctl/cmd/hooks"
	"github.com/deifyed/xctl/pkg/config"

	"github.com/deifyed/xctl/pkg/apis/xctl"

	"github.com/deifyed/xctl/pkg/cloud/linode"

	"github.com/deifyed/xctl/pkg/cloud"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

	"github.com/deifyed/xctl/pkg/tools/venv"
	"github.com/deifyed/xctl/pkg/tools/venv/shells/zsh"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type venvOpts struct {
	io                         xctl.IOStreams
	fs                         *afero.Afero
	environmentDeclarationPath string
	environmentManifest        v1alpha1.Environment
}

var (
	venvCmdOpts = venvOpts{ //nolint:gochecknoglobals
		io: xctl.IOStreams{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
		fs: &afero.Afero{Fs: afero.NewOsFs()},
	}
	venvCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "venv",
		Short: i18n.T("cmdVenvShortDescription"),
		PreRunE: hooks.EnvironmentManifestInitializer(hooks.EnvironmentManifestInitializerOpts{
			Io:                  venvCmdOpts.io,
			Fs:                  venvCmdOpts.fs,
			EnvironmentManifest: &venvCmdOpts.environmentManifest,
			SourcePath:          &venvCmdOpts.environmentDeclarationPath,
		}),
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := linode.NewLinodeProvider()
			keyringClient := keyring.Client{EnvironmentName: venvCmdOpts.environmentManifest.Metadata.Name}

			err := provider.Authenticate(keyringClient)
			if err != nil {
				return fmt.Errorf("authenticating with cloud provider: %w", err)
			}

			kubeConfigPath, err := ensureKubeConfig(ensureKubeConfigOpts{
				fs:                  venvCmdOpts.fs,
				ctx:                 cmd.Context(),
				provider:            provider,
				environmentManifest: venvCmdOpts.environmentManifest,
			})
			if err != nil {
				return fmt.Errorf("setting up kubeconfig: %w", err)
			}

			env := venv.MergeVariables(os.Environ(), []string{
				fmt.Sprintf("KUBECONFIG=%s", kubeConfigPath),
				fmt.Sprintf("XCTL_CONTEXT=%s", venvCmdOpts.environmentDeclarationPath),
			})

			workDir, err := venvCmdOpts.fs.TempDir("/tmp", "xctl")
			if err != nil {
				return fmt.Errorf("creating temporary directory: %w", err)
			}

			shellCmd, err := acquireShellCommand(venvCmdOpts.fs, workDir, env)
			if err != nil {
				return fmt.Errorf("acquiring shell: %w", err)
			}

			fmt.Fprintf(venvCmdOpts.io.Out, "Launching virtual environment\n")

			err = shellCmd.Run()
			if err != nil {
				return fmt.Errorf("running shell: %w", err)
			}

			fmt.Fprintf(venvCmdOpts.io.Out, "Successfully exited virtual environment\n")

			return nil
		},
	}
)

//nolint:gochecknoinits
func init() {
	helpers.AddEnvironmentContextFlag(venvCmd.Flags(), &venvCmdOpts.environmentDeclarationPath)

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

type ensureKubeConfigOpts struct {
	fs  *afero.Afero
	ctx context.Context

	provider            cloud.ClusterService
	environmentManifest v1alpha1.Environment
}

// ensureKubeConfig fetches a kubeconfig from provider and stores it b64 decoded in workdir as config.yaml
func ensureKubeConfig(opts ensureKubeConfigOpts) (string, error) {
	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(opts.environmentManifest.Metadata.Name)
	if err != nil {
		return "", fmt.Errorf("acquiring KubeConfig path: %w", err)
	}

	kubeConfigPath = path.Join(path.Dir(kubeConfigPath), "kubeconfig-venv.yaml")

	if _, err := os.Stat(kubeConfigPath); err == nil {
		return kubeConfigPath, nil
	}

	kubeConfig, err := opts.provider.GetKubeConfig(opts.ctx, opts.environmentManifest)
	if err != nil {
		return "", fmt.Errorf("acquiring kubeconfig: %w", err)
	}

	decodedConfig, err := base64.StdEncoding.DecodeString(string(kubeConfig))
	if err != nil {
		return "", fmt.Errorf("decoding kubeconfig: %w", err)
	}

	err = opts.fs.WriteFile(kubeConfigPath, decodedConfig, 0o600)
	if err != nil {
		return "", fmt.Errorf("writing kubeconfig: %w", err)
	}

	return kubeConfigPath, nil
}
