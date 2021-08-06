package reconciliation

import (
	"fmt"

	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/controller/common/reconciliation"
)

type clusterReconciler struct {
	clusterService cloud.ClusterService
}

func (c *clusterReconciler) Reconcile(ctx reconciliation.Context) (reconciliation.Result, error) {
	action := reconciliation.DetermineUserIndication(ctx, true)

	switch action {
	case reconciliation.ActionCreate:
		err := c.clusterService.CreateCluster(ctx.Ctx, ctx.ClusterDeclaration)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("creating cluster: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err := c.clusterService.DeleteCluster(ctx.Ctx, ctx.ClusterDeclaration)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting cluster: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, reconciliation.ErrIndecisive
}

func (c *clusterReconciler) String() string {
	return "Kubernetes Cluster"
}

func NewClusterReconciler(clusterService cloud.ClusterService) reconciliation.Reconciler {
	return &clusterReconciler{clusterService: clusterService}
}
