package grafana

import (
	"errors"
	"fmt"
	"strings"

	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/deifyed/xctl/pkg/tools/clients/helm"

	"github.com/google/uuid"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/cloud"
)

func (r reconciler) Reconcile(rctx reconciliation.Context) (reconciliation.Result, error) {
	log := logging.GetLogger("plugin", r.String())

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

		err = r.install(clients)
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

func (r reconciler) install(clients clientContainer) error {
	err := clients.secrets.Put(secretName(), map[string]string{
		adminUsernameKey: uuid.New().String(),
		adminPasswordKey: uuid.New().String(),
	})
	if err != nil {
		return fmt.Errorf("creating secrets: %w", err)
	}

	grafanaPlugin := NewPlugin()

	err = clients.helm.Install(grafanaPlugin)
	if err != nil {
		return fmt.Errorf("running helm install: %w", err)
	}

	return nil
}

func (r reconciler) uninstall(clients clientContainer) error {
	grafanaPlugin := NewPlugin()

	err := clients.secrets.Delete(secretName())
	if err != nil {
		return fmt.Errorf("deleting secret: %w", err)
	}

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
		ready           = true
	)

	_, err := r.cloudProvider.GetCluster(rctx.Ctx, rctx.EnvironmentManifest)
	if err != nil {
		if !errors.Is(err, cloud.ErrNotFound) {
			return "", fmt.Errorf("acquiring cluster: %w", err)
		}

		clusterExists = false
	}

	plugin := NewPlugin()

	if clusterExists {
		componentExists, err = helm.Exists(plugin)
		if err != nil {
			return "", fmt.Errorf("checking component existence: %w", err)
		}

		ready, err = kubectlClient.IsReady(kubectl.Selector{
			Namespace: pluginNamespace,
			Kind:      kubectl.DeploymentKind,
			Name:      pluginName,
		})
		if err != nil {
			if !errors.Is(err, kubectl.ErrNotFound) {
				return "", fmt.Errorf("checking ready status: %w", err)
			}

			ready = false
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

		if componentExists && !ready {
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
	return strings.Title(pluginName)
}

func NewReconciler(cloudProvider cloud.Provider) reconciliation.Reconciler {
	return &reconciler{
		cloudProvider: cloudProvider,
	}
}
