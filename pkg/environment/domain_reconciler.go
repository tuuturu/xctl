package environment

import (
	"fmt"

	reconciliation2 "github.com/deifyed/xctl/pkg/tools/reconciliation"

	"github.com/pkg/errors"

	"github.com/deifyed/xctl/pkg/cloud"
)

func (d *domainReconciler) Reconcile(ctx reconciliation2.Context) (reconciliation2.Result, error) {
	cluster, err := d.clusterService.GetCluster(ctx.Ctx, ctx.ClusterDeclaration)
	if err != nil {
		if !errors.Is(err, cloud.ErrNotFound) {
			return reconciliation2.Result{}, fmt.Errorf("retrieving cluster: %w", err)
		}
	}

	action, err := d.determineAction(ctx, cluster)
	if err != nil {
		return reconciliation2.Result{}, fmt.Errorf("determining action: %w", err)
	}

	domain := cloud.Domain{Host: fmt.Sprintf("*.%s", ctx.ClusterDeclaration.Spec.RootDomain)}

	switch action {
	case reconciliation2.ActionCreate:
		_, err = d.domainService.CreateSubdomain(ctx.Ctx, domain, cluster.PublicIPv6)
		if err != nil {
			return reconciliation2.Result{}, fmt.Errorf("creating subdomain: %w", err)
		}

		return reconciliation2.Result{Requeue: false}, nil
	case reconciliation2.ActionDelete:
		err = d.domainService.DeleteSubdomain(ctx.Ctx, domain)
		if err != nil {
			return reconciliation2.Result{}, fmt.Errorf("deleting subdomain: %w", err)
		}

		return reconciliation2.Result{Requeue: false}, err
	}

	return reconciliation2.NoopWaitIndecisiveHandler(action)
}

func (d *domainReconciler) determineAction(ctx reconciliation2.Context, cluster cloud.Cluster) (reconciliation2.Action, error) { //nolint:lll
	userIndication := reconciliation2.DetermineUserIndication(ctx, true)
	domain := cloud.Domain{Host: ctx.ClusterDeclaration.Spec.RootDomain}

	hasPrimaryDomain, err := d.domainService.HasPrimaryDomain(ctx.Ctx, domain)
	if err != nil {
		return reconciliation2.ActionNoop, fmt.Errorf("checking primary domain: %w", err)
	}

	if !hasPrimaryDomain {
		return reconciliation2.ActionNoop, fmt.Errorf(
			"the primary domain %s is not available in your account",
			domain.PrimaryDomain(),
		)
	}

	hasSubdomain, err := d.domainService.HasSubdomain(ctx.Ctx, domain)
	if err != nil {
		return "", fmt.Errorf("checking for subdomain: %w", err)
	}

	switch userIndication {
	case reconciliation2.ActionCreate:
		if cluster.PublicIPv6 == "" {
			return reconciliation2.ActionWait, nil
		}

		if hasSubdomain {
			return reconciliation2.ActionNoop, nil
		}

		return reconciliation2.ActionCreate, nil
	case reconciliation2.ActionDelete:
		if !hasSubdomain {
			return reconciliation2.ActionNoop, nil
		}

		return reconciliation2.ActionDelete, nil
	}

	return reconciliation2.ActionNoop, reconciliation2.ErrIndecisive
}

func (d *domainReconciler) String() string {
	return domainReconcilerName
}

func NewDomainReconciler(provider cloud.Provider) reconciliation2.Reconciler {
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
