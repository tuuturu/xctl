package environment

import (
	"fmt"
	"io"
	"strings"

	"github.com/deifyed/xctl/pkg/plugins/certmanager"
	"github.com/deifyed/xctl/pkg/plugins/grafana"
	"github.com/deifyed/xctl/pkg/plugins/loki"
	ingress "github.com/deifyed/xctl/pkg/plugins/nginx-ingress-controller"
	"github.com/deifyed/xctl/pkg/plugins/prometheus"
	"github.com/deifyed/xctl/pkg/plugins/promtail"

	"github.com/deifyed/xctl/pkg/plugins/argocd"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"sigs.k8s.io/yaml"

	_ "embed"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/tools/spinner"
)

// Reconcile knows how to ensure reality for an environment is as declared in an environment manifest
func Reconcile(opts ReconcileOpts) error {
	log := logging.GetLogger("cmd/apply", "environment")

	manifest, err := ExtractManifest(opts.Manifest)
	if err != nil {
		return fmt.Errorf("extracting manifest: %w", err)
	}

	var spinnerOut io.Writer
	if logging.GetLevel() == logging.LevelInfo {
		spinnerOut = opts.Out
	} else {
		spinnerOut = io.Discard
	}

	spin := spinner.NewSpinner(spinnerOut)
	spin.FinalMSG = "✅"

	schedulerOpts := reconciliation.SchedulerOpts{
		Filesystem:                      opts.Filesystem,
		Out:                             opts.Out,
		PurgeFlag:                       opts.Purge,
		EnvironmentManifest:             manifest,
		ReconciliationLoopDelayFunction: reconciliation.DefaultDelayFunction,
		QueueStepFunc: func(identifier string) {
			log.Debug(fmt.Sprintf("reconciling %s", identifier))

			spin.Suffix = fmt.Sprintf(" Reconciling %s", identifier)
		},
	}

	scheduler := reconciliation.NewScheduler(schedulerOpts,
		NewClusterReconciler(opts.Provider),
		ingress.NewReconciler(opts.Provider),
		NewDomainReconciler(opts.Provider),
		certmanager.NewReconciler(opts.Provider),
		prometheus.NewReconciler(opts.Provider),
		grafana.NewReconciler(opts.Provider),
		loki.NewReconciler(opts.Provider),
		promtail.NewReconciler(opts.Provider),
		argocd.NewReconciler(opts.Provider),
	)

	spin.Start()

	_, err = scheduler.Run(opts.Context)
	if err != nil {
		spin.FinalMSG = "❌\n\n"
		spin.Stop()

		return fmt.Errorf("scheduling: %w", err)
	}

	spin.Suffix = "Reconciliation finished"

	spin.Stop()

	fmt.Fprintf(opts.Out, "\n\nReconciliation complete\n")

	return nil
}

//go:embed environment-template.yaml
var environmentTemplate string //nolint:gochecknoglobals

// Scaffold returns a stream containing an environment template
func Scaffold() io.Reader {
	return strings.NewReader(environmentTemplate)
}

// ExtractManifest knows how to produce an environment manifest from a reader source
func ExtractManifest(source io.Reader) (v1alpha1.Environment, error) {
	manifest := v1alpha1.NewDefaultEnvironment()

	rawManifest, err := io.ReadAll(source)
	if err != nil {
		return v1alpha1.Environment{}, fmt.Errorf("reading: %w", err)
	}

	err = yaml.Unmarshal(rawManifest, &manifest)
	if err != nil {
		return v1alpha1.Environment{}, fmt.Errorf("unmarshalling: %w", err)
	}

	return manifest, nil
}
