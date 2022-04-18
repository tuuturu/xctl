package grafana

import (
	"errors"
	"fmt"

	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/deifyed/xctl/pkg/tools/clients/helm"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/google/uuid"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/cloud"
)

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	log := logging.GetLogger(logFeature, "reconciliation")

	clients, err := prepareClients(rctx.Filesystem, rctx.EnvironmentManifest)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("preparing clients: %w", err)
	}

	action, err := r.determineAction(rctx, clients.helm, clients.kubectl)
	if err != nil {
		return reconciliation.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		log.Debug("installing")

		err = r.install(clients, rctx.EnvironmentManifest)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("installing: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		log.Debug("deleting")

		err = r.uninstall(clients)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("uninstalling: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.NoopWaitIndecisiveHandler(action)
}

func (r reconciler) install(clients clientContainer, cluster v1alpha1.Environment) error {
	username := uuid.New().String()
	password := uuid.New().String()

	err := clients.secrets.Put(secretName(), map[string]string{
		"adminUsername": username,
		"adminPassword": password,
	}) //nolint:godox    // TODO: Injecting into template is not ok
	if err != nil {
		return fmt.Errorf("creating secrets: %w", err)
	}

	grafanaPlugin, err := NewPlugin(NewPluginOpts{
		Host:          fmt.Sprintf("grafana.%s", cluster.Spec.Domain),
		AdminUsername: username,
		AdminPassword: password,
	})
	if err != nil {
		return fmt.Errorf("creating plugin: %w", err)
	}

	err = clients.helm.Install(grafanaPlugin)
	if err != nil {
		return fmt.Errorf("running helm install: %w", err)
	}

	return nil
}

func (r reconciler) uninstall(clients clientContainer) error {
	grafanaPlugin, err := NewPlugin(NewPluginOpts{})
	if err != nil {
		return fmt.Errorf("creating plugin: %w", err)
	}

	// err = secretsClient.Delete("grafana")

	err = clients.helm.Delete(grafanaPlugin)
	if err != nil {
		return fmt.Errorf("running helm delete: %w", err)
	}

	return nil
}

//nolint:funlen
func (r reconciler) determineAction(rctx reconciliation.Context, helm helm.Client, kubectlClient kubectl.Client) (reconciliation.Action, error) { //nolint:lll
	indication := reconciliation.DetermineUserIndication(rctx, rctx.EnvironmentManifest.Spec.Plugins.Grafana)

	var (
		clusterExists   = true
		componentExists = true
		vaultReady      = true
	)

	_, err := r.cloudProvider.GetCluster(rctx.Ctx, rctx.EnvironmentManifest)
	if err != nil {
		if !errors.Is(err, cloud.ErrNotFound) {
			return "", fmt.Errorf("acquiring cluster: %w", err)
		}

		clusterExists = false
	}

	plugin, err := NewPlugin(NewPluginOpts{})
	if err != nil {
		return "", fmt.Errorf("preparing plugin: %w", err)
	}

	if clusterExists {
		componentExists, err = helm.Exists(plugin)
		if err != nil {
			return "", fmt.Errorf("checking component existence: %w", err)
		}

		vaultReady, err = kubectlClient.PodReady(kubectl.Pod{
			Name:      "vault-0",
			Namespace: "kube-system",
		})
		if err != nil {
			if !errors.Is(err, kubectl.ErrNotFound) {
				return "", fmt.Errorf("checking vault ready status: %w", err)
			}

			vaultReady = false
		}
	}

	switch indication {
	case reconciliation.ActionCreate:
		if !clusterExists {
			return reconciliation.ActionWait, nil
		}

		if componentExists {
			return reconciliation.ActionNoop, nil
		}

		if !vaultReady {
			return reconciliation.ActionWait, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !clusterExists {
			return reconciliation.ActionNoop, nil
		}

		if !componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

func (r reconciler) String() string {
	return "Grafana"
}

func NewReconciler(cloudProvider cloud.Provider) reconciliation.Reconciler {
	return &reconciler{
		cloudProvider: cloudProvider,
	}
}
