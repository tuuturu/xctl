package cmd

import (
	"github.com/deifyed/xctl/cmd/handlers"
	"github.com/spf13/cobra"
)

var (
	applyCmdOpts = handlers.ApplyRunEOpts{Purge: false} //nolint:gochecknoglobals
	applyCmd     = &cobra.Command{                      //nolint:gochecknoglobals
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
