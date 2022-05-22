package cmd

import (
	"os"

	"github.com/deifyed/xctl/cmd/hooks"
	"github.com/deifyed/xctl/pkg/apis/xctl"
	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/environment"
	"github.com/spf13/afero"

	"github.com/deifyed/xctl/cmd/helpers"
	"github.com/deifyed/xctl/pkg/tools/i18n"
	"github.com/spf13/cobra"
)

var (
	loginCmdOpts = loginCmdOptsContainer{
		io: xctl.IOStreams{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
		fs: &afero.Afero{Fs: afero.NewOsFs()},
	}
	loginCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "login",
		Short: i18n.T("cmdLoginShortDescription"),
		Args:  cobra.ExactArgs(0),
		PreRunE: hooks.EnvironmentManifestInitializer(hooks.EnvironmentManifestInitializerOpts{
			Io:                  loginCmdOpts.io,
			Fs:                  loginCmdOpts.fs,
			EnvironmentManifest: &loginCmdOpts.EnvironmentContext,
			SourcePath:          &loginCmdOpts.EnvironmentContextPath,
		}),
		RunE: environment.Authenticate(&loginCmdOpts.EnvironmentContext),
	}
)

//nolint:gochecknoinits
func init() {
	helpers.AddEnvironmentContextFlag(loginCmd.Flags(), &loginCmdOpts.EnvironmentContextPath)

	rootCmd.AddCommand(loginCmd)
}

type loginCmdOptsContainer struct {
	EnvironmentContextPath string
	EnvironmentContext     v1alpha1.Environment
	fs                     *afero.Afero
	io                     xctl.IOStreams
}
