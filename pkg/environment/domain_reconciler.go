package environment

import (
	"fmt"

	"github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/pkg/errors"

	"github.com/deifyed/xctl/pkg/cloud"
)

// Reconcile knows how to ensure reality for a domain is as declared in an environment manifest
func (d *domainReconciler) Reconcile(ctx reconciliation.Context) (reconciliation.Result, error) {
	cluster, err := d.clusterService.GetCluster(ctx.Ctx, ctx.EnvironmentManifest)
	if err != nil {
		if !errors.Is(err, cloud.ErrNotFound) {
			return reconciliation.Result{}, fmt.Errorf("retrieving cluster: %w", err)
		}
	}

	action, err := d.determineAction(ctx, cluster)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining action: %w", err)
	}

	domain := cloud.Domain{Host: fmt.Sprintf("*.%s", ctx.EnvironmentManifest.Spec.Domain)}

	switch action {
	case reconciliation.ActionCreate:
		_, err = d.domainService.CreateSubdomain(ctx.Ctx, domain, cluster.PublicIPv6)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("creating subdomain: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = d.domainService.DeleteSubdomain(ctx.Ctx, domain)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting subdomain: %w", err)
		}

		return reconciliation.Result{Requeue: false}, err
	}

	return reconciliation.NoopWaitIndecisiveHandler(action)
}

func (d *domainReconciler) determineAction(ctx reconciliation.Context, cluster cloud.Cluster) (reconciliation.Action, error) { //nolint:lll
	userIndication := reconciliation.DetermineUserIndication(ctx, true)
	domain := cloud.Domain{Host: ctx.EnvironmentManifest.Spec.Domain}

	hasPrimaryDomain, err := d.domainService.HasPrimaryDomain(ctx.Ctx, domain)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("checking primary domain: %w", err)
	}

	if !hasPrimaryDomain {
		return reconciliation.ActionNoop, fmt.Errorf(
			"the primary domain %s is not available in your account",
			domain.PrimaryDomain(),
		)
	}

	hasSubdomain, err := d.domainService.HasSubdomain(ctx.Ctx, domain)
	if err != nil {
		return "", fmt.Errorf("checking for subdomain: %w", err)
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		if cluster.PublicIPv6 == "" {
			return reconciliation.ActionWait, nil
		}

		if hasSubdomain {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !hasSubdomain {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns a string representing the reconciler
func (d *domainReconciler) String() string {
	return domainReconcilerName
}

// NewDomainReconciler returns an initialized domain reconciler
func NewDomainReconciler(provider cloud.Provider) reconciliation.Reconciler {
	return &domainReconciler{
		domainService:  provider,
		clusterService: provider,
	}
}

const domainReconcilerName = "Cluster Domain"

type domainReconciler struct {
	clusterService cloud.ClusterService
	domainService  cloud.DomainService
}
