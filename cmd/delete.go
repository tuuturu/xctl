package cmd

import (
	"github.com/deifyed/xctl/cmd/handlers"

	"github.com/spf13/cobra"
)

var (
	deleteCmdOpts = handlers.ApplyRunEOpts{Purge: true} //nolint:gochecknoglobals
	deleteCmd     = &cobra.Command{                     //nolint:gochecknoglobals
		Use:   "delete",
		Short: "deletes a resource",
		RunE:  handlers.ApplyRunE(&deleteCmdOpts),
	}
)

//nolint:gochecknoinits
func init() {
	flags := deleteCmd.Flags()

	flags.StringVarP(&deleteCmdOpts.File, "file", "f", "-", "file representing resource to delete")

	rootCmd.AddCommand(deleteCmd)
}
