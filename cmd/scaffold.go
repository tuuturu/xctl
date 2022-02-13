package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/deifyed/xctl/pkg/tools/i18n"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/tools/scaffolding"
	"github.com/spf13/cobra"
)

var scaffoldCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "scaffold",
	Short: i18n.T("cmdScaffoldShortDescription"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resource := strings.ToLower(args[0])

		switch resource {
		case strings.ToLower(v1alpha1.ClusterKind):
			_, err := io.Copy(cmd.OutOrStdout(), scaffolding.Cluster())
			if err != nil {
				return fmt.Errorf("scaffolding cluster template: %w", err)
			}
		default:
			return fmt.Errorf("unable to recognize resource type \"%s\". Valid resource types are: %+v",
				resource,
				[]string{strings.ToLower(v1alpha1.ClusterKind)},
			)
		}

		return nil
	},
}

//nolint:gochecknoinits
func init() {
	rootCmd.AddCommand(scaffoldCmd)
}
