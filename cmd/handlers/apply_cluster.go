package handlers

import (
	"context"
	"fmt"
	"io"
	"time"

	ingress "github.com/deifyed/xctl/pkg/plugins/nginx-ingress-controller"

	"github.com/deifyed/xctl/pkg/plugins/vault"

	"github.com/spf13/afero"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/cloud/linode"
	"github.com/deifyed/xctl/pkg/config"
	clusterrec "github.com/deifyed/xctl/pkg/controller/cluster/reconciliation"
	"github.com/deifyed/xctl/pkg/controller/common/reconciliation"
	"github.com/deifyed/xctl/pkg/tools/spinner"

	"sigs.k8s.io/yaml"
)

func handleCluster(fs *afero.Afero, out io.Writer, purge bool, clusterManifestSource io.Reader) error {
	var manifest v1alpha1.Cluster

	content, err := io.ReadAll(clusterManifestSource)
	if err != nil {
		return fmt.Errorf("reading cluster manifest: %w", err)
	}

	err = yaml.Unmarshal(content, &manifest)
	if err != nil {
		return fmt.Errorf("parsing cluster manifest: %w", err)
	}

	provider := linode.NewLinodeProvider()

	err = provider.Authenticate()
	if err != nil {
		return fmt.Errorf("authenticating with cloud provider: %w", err)
	}

	spin := spinner.NewSpinner(out)
	spin.FinalMSG = "✅"

	opts := reconciliation.SchedulerOpts{
		Filesystem:                      fs,
		Out:                             out,
		PurgeFlag:                       purge,
		ClusterDeclaration:              manifest,
		ReconciliationLoopDelayFunction: func() { time.Sleep(config.DefaultReconciliationLoopDelayDuration) },
		QueueStepFunc: func(identifier string) {
			spin.Suffix = fmt.Sprintf(" Reconciling %s", identifier)
		},
	}

	scheduler := reconciliation.NewScheduler(opts,
		clusterrec.NewClusterReconciler(provider),
		vault.NewVaultReconciler(provider),
		ingress.NewNginxIngressControllerReconciler(provider),
	)

	spin.Start()

	_, err = scheduler.Run(context.Background())

	spin.Stop()

	fmt.Fprintf(out, "\n\nCluster reconciliation complete")

	return err
}
