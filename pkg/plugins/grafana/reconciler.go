package grafana

import (
	"errors"
	"fmt"

	reconciliation2 "github.com/deifyed/xctl/pkg/tools/reconciliation"

	helm2 "github.com/deifyed/xctl/pkg/tools/clients/helm"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/google/uuid"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/cloud"
)

func (r reconciler) Reconcile(rctx reconciliation2.Context) (reconciliation2.Result, error) {
	log := logging.GetLogger(logFeature, "reconciliation")

	clients, err := prepareClients(rctx.Filesystem, rctx.ClusterDeclaration)
	if err != nil {
		return reconciliation2.Result{}, fmt.Errorf("preparing clients: %w", err)
	}

	stopFn, err := openVaultConnection(clients.kubectl)
	if err != nil {
		return reconciliation2.Result{}, fmt.Errorf("opening vault connection: %w", err)
	}

	defer func() {
		_ = stopFn()
	}()

	action, err := r.determineAction(rctx, clients.helm)
	if err != nil {
		return reconciliation2.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation2.ActionCreate:
		log.Debug("installing")

		err = r.install(clients, rctx.ClusterDeclaration)
		if err != nil {
			if errors.Is(err, helm2.ErrUnreachable) {
				return reconciliation2.Result{Requeue: true}, nil
			}

			return reconciliation2.Result{}, fmt.Errorf("installing: %w", err)
		}

		return reconciliation2.Result{Requeue: false}, nil
	case reconciliation2.ActionDelete:
		log.Debug("deleting")

		err = r.uninstall(clients)
		if err != nil {
			if errors.Is(err, helm2.ErrUnreachable) {
				return reconciliation2.Result{Requeue: true}, nil
			}

			return reconciliation2.Result{}, fmt.Errorf("uninstalling: %w", err)
		}

		return reconciliation2.Result{Requeue: false}, nil
	}

	return reconciliation2.NoopWaitIndecisiveHandler(action)
}

func (r reconciler) install(clients clientContainer, cluster v1alpha1.Cluster) error {
	username := uuid.New().String()
	password := uuid.New().String()

	err := clients.secrets.Put("grafana", map[string]string{
		"adminUsername": username,
		"adminPassword": password,
	}) //nolint:godox    // TODO: Injecting into template is not ok
	if err != nil {
		return fmt.Errorf("creating secrets: %w", err)
	}

	grafanaPlugin, err := NewPlugin(NewPluginOpts{
		Host:          fmt.Sprintf("grafana.%s", cluster.Spec.RootDomain),
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

func (r reconciler) determineAction(rctx reconciliation2.Context, helm helm2.Client) (reconciliation2.Action, error) { //nolint:lll
	indication := reconciliation2.DetermineUserIndication(rctx, rctx.ClusterDeclaration.Spec.Plugins.Grafana)

	var (
		clusterExists   = true
		componentExists = true
	)

	_, err := r.cloudProvider.GetCluster(rctx.Ctx, rctx.ClusterDeclaration)
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
	}

	switch indication {
	case reconciliation2.ActionCreate:
		if !clusterExists {
			return reconciliation2.ActionWait, nil
		}

		if componentExists {
			return reconciliation2.ActionNoop, nil
		}

		return reconciliation2.ActionCreate, nil
	case reconciliation2.ActionDelete:
		if !clusterExists {
			return reconciliation2.ActionNoop, nil
		}

		if !componentExists {
			return reconciliation2.ActionNoop, nil
		}

		return reconciliation2.ActionDelete, nil
	}

	return reconciliation2.ActionNoop, reconciliation2.ErrIndecisive
}

func (r reconciler) String() string {
	return "Grafana"
}

func NewReconciler(cloudProvider cloud.Provider) reconciliation2.Reconciler {
	return &reconciler{
		cloudProvider: cloudProvider,
	}
}