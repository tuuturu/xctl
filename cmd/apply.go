package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/deifyed/xctl/pkg/tools/spinner"

	"github.com/deifyed/xctl/pkg/cloud/linode"
	"github.com/deifyed/xctl/pkg/config"
	clusterrec "github.com/deifyed/xctl/pkg/controller/cluster/reconciliation"
	"github.com/deifyed/xctl/pkg/controller/common/reconciliation"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
				return handleCluster(false, rawContent)
			case v1alpha1.ApplicationKind:
				return handleApplication(false, rawContent)
			default:
				return fmt.Errorf("unknown kind %s", kind)
			}
		},
	}
)

func handleCluster(purge bool, content []byte) error {
	var manifest v1alpha1.Cluster

	err := yaml.Unmarshal(content, &manifest)
	if err != nil {
		return fmt.Errorf("parsing cluster manifest: %w", err)
	}

	provider := linode.NewLinodeProvider()

	err = provider.Authenticate()
	if err != nil {
		return fmt.Errorf("authenticating with cloud provider: %w", err)
	}

	spin := spinner.NewSpinner(os.Stdout)

	opts := reconciliation.SchedulerOpts{
		Out:                             os.Stdout,
		PurgeFlag:                       purge,
		ClusterDeclaration:              manifest,
		ReconciliationLoopDelayFunction: func() { time.Sleep(config.DefaultReconciliationLoopDelayDuration) },
		QueueStepFunc: func(identifier string) {
			spin.Suffix = fmt.Sprintf(" reconciling %s", identifier)
		},
	}

	scheduler := reconciliation.NewScheduler(opts,
		clusterrec.NewClusterReconciler(provider),
	)

	spin.Start()
	defer spin.Stop()

	_, err = scheduler.Run(context.Background())

	return err
}

func handleApplication(_ bool, content []byte) error {
	var manifest v1alpha1.Application

	err := yaml.Unmarshal(content, &manifest)
	if err != nil {
		return fmt.Errorf("parsing application manifest: %w", err)
	}

	println(fmt.Sprintf("finished handling %s", manifest.Metadata.Name))

	return nil
}

//nolint:gochecknoinits
func init() {
	flags := applyCmd.Flags()

	flags.StringVarP(&applyCmdOpts.File, "file", "f", "-", "file to apply")

	rootCmd.AddCommand(applyCmd)
}
