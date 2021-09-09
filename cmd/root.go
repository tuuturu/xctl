package cmd

import (
	"fmt"
	"os"

	"github.com/deifyed/xctl/cmd/helpers"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "xctl",
	Short: "xctl provisions things",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		userError := helpers.ErrorTranslator(err)

		_, _ = fmt.Fprintln(os.Stderr, userError)

		os.Exit(1)
	}
}
