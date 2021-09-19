package cmd

import (
	"os"

	"github.com/deifyed/xctl/cmd/handlers"
	"github.com/deifyed/xctl/pkg/apis/xctl"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	applyCmdOpts = handlers.ApplyRunEOpts{ //nolint:gochecknoglobals
		Io: xctl.IOStreams{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
		Filesystem: &afero.Afero{Fs: afero.NewOsFs()},
		Purge:      false,
	}
	applyCmd = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "apply",
		Short: "applies a manifest",
		RunE:  handlers.ApplyRunE(&applyCmdOpts),
	}
)

//nolint:gochecknoinits
func init() {
	flags := applyCmd.Flags()

	flags.StringVarP(&applyCmdOpts.File, "file", "f", "-", "file to apply")

	rootCmd.AddCommand(applyCmd)
}
