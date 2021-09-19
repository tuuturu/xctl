package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/deifyed/xctl/cmd/preruns"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "xctl",
	Short: "xctl provisions a known and complete production environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute(errOut io.Writer) {
	if err := rootCmd.Execute(); err != nil {
		userError := preruns.ErrorTranslator(err)

		_, _ = fmt.Fprintln(errOut, userError)

		os.Exit(1)
	}
}
