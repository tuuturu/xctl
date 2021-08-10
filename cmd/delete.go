package cmd

import (
	"fmt"
	"os"

	"github.com/deifyed/xctl/cmd/handlers"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type deleteOpts struct {
	File string
}

var (
	deleteCmdOpts deleteOpts        //nolint:gochecknoglobals
	deleteCmd     = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "delete",
		Short: "deletes a resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := os.Stdout
			fs := &afero.Afero{Fs: afero.NewOsFs()}

			rawContent, err := fs.ReadFile(deleteCmdOpts.File)
			if err != nil {
				return fmt.Errorf("reading file: %w", err)
			}

			kind, err := v1alpha1.InferKindFromManifest(rawContent)
			if err != nil {
				return fmt.Errorf("inferring kind: %w", err)
			}

			switch kind {
			case v1alpha1.ClusterKind:
				fmt.Fprintf(out, "Deleting resources associated with cluster manifest, please wait\n\n")

				return handlers.HandleCluster(out, true, rawContent)
			case v1alpha1.ApplicationKind:
				fmt.Fprintf(out, "Deleting resources associated with application manifest %s, please wait\n\n", deleteCmdOpts.File)

				return handlers.HandleApplication(out, true, rawContent)
			default:
				return fmt.Errorf("unknown kind %s", kind)
			}
		},
	}
)

//nolint:gochecknoinits
func init() {
	flags := deleteCmd.Flags()

	flags.StringVarP(&deleteCmdOpts.File, "file", "f", "-", "file representing resource to delete")

	rootCmd.AddCommand(deleteCmd)
}
