package cmd

import (
	"os"

	"github.com/deifyed/xctl/cmd/handlers"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	applyCmdOpts = handlers.ApplyRunEOpts{ //nolint:gochecknoglobals
		InternalFilesystem: &afero.Afero{Fs: afero.NewMemMapFs()},
		ExternalFilesystem: &afero.Afero{Fs: afero.NewOsFs()},
		Out:                os.Stdout,
		Purge:              false,
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
