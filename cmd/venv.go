package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/deifyed/xctl/pkg/cloud/linode"

	"github.com/deifyed/xctl/pkg/tools/venv"
	"github.com/deifyed/xctl/pkg/tools/venv/shells/zsh"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type venvOpts struct {
	clusterDeclarationPath string
}

var (
	venvCmdOpts venvOpts          //nolint:gochecknoglobals
	venvCmd     = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "venv",
		Short: "activates a virtual environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			fs := &afero.Afero{Fs: afero.NewOsFs()}

			provider := linode.NewLinodeProvider()

			err := provider.Authenticate()
			if err != nil {
				return fmt.Errorf("authenticating with cloud provider: %w", err)
			}

			shellPath, err := venv.GetCurrentShell(fs)
			if err != nil {
				return fmt.Errorf("getting current shell path: %w", err)
			}

			var shell venv.Shell

			switch {
			case strings.HasSuffix(shellPath, "zsh"):
				shell = zsh.NewZshShell(fs, shellPath)
			default:
				return fmt.Errorf("handling shell %s. Not a known shell", shellPath)
			}

			env := venv.MergeVariables(os.Environ(), []string{
				fmt.Sprintf("KUBECONFIG=%s", "something"),
			})

			shellCmd, err := shell.Command(env)
			if err != nil {
				return fmt.Errorf("running virtual environment shell: %w", err)
			}

			fmt.Printf("Launching virtual environment")

			err = shellCmd.Run()
			if err != nil {
				return fmt.Errorf("running shell: %w", err)
			}

			err = shell.Teardown()
			if err != nil {
				return fmt.Errorf("tearing down shell: %w", err)
			}

			fmt.Printf("Successfully exited virtual environment")

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
