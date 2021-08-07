package cmd

import (
	"fmt"
	"os"

	"github.com/deifyed/xctl/cmd/handlers"
	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type applyOpts struct {
	File string
}

var (
	applyCmdOpts applyOpts         //nolint:gochecknoglobals
	applyCmd     = &cobra.Command{ //nolint:gochecknoglobals
		Use:   "apply",
		Short: "applies a manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := os.Stdout
			fs := &afero.Afero{Fs: afero.NewOsFs()}

			rawContent, err := fs.ReadFile(applyCmdOpts.File)
			if err != nil {
				return fmt.Errorf("reading file: %w", err)
			}

			kind, err := v1alpha1.InferKindFromManifest(rawContent)
			if err != nil {
				return fmt.Errorf("inferring kind: %w", err)
			}

			switch kind {
			case v1alpha1.ClusterKind:
				fmt.Fprintf(out, "Applying cluster manifest %s, please wait\n\n", applyCmdOpts.File)

				return handlers.HandleCluster(out, false, rawContent)
			case v1alpha1.ApplicationKind:
				fmt.Fprintf(out, "Applying application manifest %s, please wait\n\n", applyCmdOpts.File)

				return handlers.HandleApplication(out, false, rawContent)
			default:
				return fmt.Errorf("unknown kind %s", kind)
			}
		},
	}
)

//nolint:gochecknoinits
func init() {
	flags := applyCmd.Flags()

	flags.StringVarP(&applyCmdOpts.File, "file", "f", "-", "file to apply")

	rootCmd.AddCommand(applyCmd)
}
