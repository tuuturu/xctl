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
	deleteCmdOpts = handlers.ApplyRunEOpts{ //nolint:gochecknoglobals
		Io: xctl.IOStreams{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
		Filesystem: &afero.Afero{Fs: afero.NewOsFs()},
		Purge:      true,
	}
	deleteCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "delete",
		Short: i18n.T("cmdDeleteShortDescription"),
		RunE:  handlers.ApplyRunE(&deleteCmdOpts),
	}
)

//nolint:gochecknoinits
func init() {
	flags := deleteCmd.Flags()

	flags.StringVarP(&deleteCmdOpts.File, "file", "f", "-", "file representing resource to delete")

	rootCmd.AddCommand(deleteCmd)
}
