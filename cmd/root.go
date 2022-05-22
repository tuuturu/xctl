package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/deifyed/xctl/pkg/tools/i18n"

	"github.com/deifyed/xctl/cmd/helpers"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "xctl",
	Short:         i18n.T("cmdRootShortDecsription"),
	Version:       "0.0.alpha",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	helpers.InitializeLogging()

	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		userError := helpers.ErrorTranslator(err)

		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", userError)

		os.Exit(1)
	}
}
