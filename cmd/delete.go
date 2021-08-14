package cmd

import (
	"fmt"
	"io"
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
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			out := os.Stdout
			fs := &afero.Afero{Fs: afero.NewOsFs()}

			var manifestSource io.Reader

			if applyCmdOpts.File == "-" {
				manifestSource = os.Stdin
			} else {
				manifestSource, err = fs.Open(applyCmdOpts.File)
				if err != nil {
					return fmt.Errorf("opening manifest file: %w", err)
				}
			}

			kind, err := v1alpha1.InferKindFromManifest(manifestSource)
			if err != nil {
				return fmt.Errorf("inferring kind: %w", err)
			}

			switch kind {
			case v1alpha1.ClusterKind:
				fmt.Fprintf(out, "Deleting resources associated with cluster manifest, please wait\n\n")

				return handlers.HandleCluster(out, true, manifestSource)
			case v1alpha1.ApplicationKind:
				fmt.Fprintf(out, "Deleting resources associated with application manifest %s, please wait\n\n", deleteCmdOpts.File)

				return handlers.HandleApplication(out, true, manifestSource)
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
