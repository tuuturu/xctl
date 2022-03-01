package cmd

import (
	"os"

	"github.com/deifyed/xctl/pkg/tools/i18n"

	"github.com/deifyed/xctl/cmd/handlers"
	"github.com/deifyed/xctl/pkg/apis/xctl"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	applyCmdOpts = handlers.ApplyRunEOpts{ //nolint:gochecknoglobals
		Io: xctl.IOStreams{
			Out: os.Stdout,
			Err: os.Stderr,
			In:  os.Stdin,
		},
		Filesystem: &afero.Afero{Fs: afero.NewOsFs()},
		Purge:      false,
	}
	applyCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "apply",
		Short: i18n.T("cmdApplyShortDescription"),
		RunE:  handlers.ApplyRunE(&applyCmdOpts),
	}
)

//nolint:gochecknoinits
func init() {
	flags := applyCmd.Flags()

	flags.StringVarP(&applyCmdOpts.File, "file", "f", "-", "file to apply")

	rootCmd.AddCommand(applyCmd)
}
