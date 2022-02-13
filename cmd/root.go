package cmd

import (
	"fmt"
	"os"

	"github.com/deifyed/xctl/pkg/tools/i18n"

	"github.com/deifyed/xctl/cmd/helpers"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "xctl",
	Short: i18n.T("cmdRootShortDecsription"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() {
	helpers.InitializeLogging()

	if err := rootCmd.Execute(); err != nil {
		userError := helpers.ErrorTranslator(err)

		_, _ = fmt.Fprintln(os.Stderr, userError)

		os.Exit(1)
	}
}
