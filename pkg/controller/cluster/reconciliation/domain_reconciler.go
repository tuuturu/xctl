package reconciliation

import (
	"fmt"

	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/controller/common/reconciliation"
)

type domainReconciler struct {
	domainService cloud.DomainService
}

func (d *domainReconciler) Reconcile(ctx reconciliation.Context) (reconciliation.Result, error) {
	action, err := d.determineAction(ctx)
	if err != nil {
		return reconciliation.Result{Requeue: false}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		_, err := d.domainService.CreateSubdomain(ctx.Ctx, ctx.ClusterDeclaration.URL)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("creating domain: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err := d.domainService.DeleteSubdomain(ctx.Ctx, ctx.ClusterDeclaration.URL)
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

func (d *domainReconciler) determineAction(rctx reconciliation.Context) (reconciliation.Action, error) {
	action := reconciliation.DetermineUserIndication(rctx, true)

	var (
		primaryDomainExists bool
		subdomainExists     bool
		err                 error
	)

	primaryDomainExists, err = d.domainService.HasPrimaryDomain(rctx.Ctx, rctx.ClusterDeclaration.URL)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("checking primary domain existence: %w", err)
	}

	if primaryDomainExists {
		subdomainExists, err = d.domainService.HasSubdomain(rctx.Ctx, rctx.ClusterDeclaration.URL)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("checking subdomain existence: %w", err)
		}
	}

	switch action {
	case reconciliation.ActionCreate:
		if !primaryDomainExists {
			return reconciliation.ActionNoop, nil // Should probably return error
		}

		if subdomainExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !subdomainExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

func (d *domainReconciler) String() string {
	return "Domains"
}

func NewDomainReconciler(domainService cloud.DomainService) reconciliation.Reconciler {
	return &domainReconciler{domainService: domainService}
}
