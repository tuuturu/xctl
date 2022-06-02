
## Accessing ArgoCD

First, acquire credentials for Grafana by running the `get credentials` command:

```shell
xctl --context env.yaml get credentials argocd
```

After you've gotten the credentials, open a tunnel to the Grafana instance by running the following command:

```shell
xctl --context env.yaml forward argocd
```

## Preexisting ArgoCD applications

### Applications

This ArgoCD application tracks the `/infrastructure/<environment name>/argocd/applications/` directory. This directory
contains ArgoCD applications referencing actual applications. This setup has the following effects:

- We can deploy a new application without administrating the environment directly by adding an ArgoCD application
    manifest to this directory
- This renders the ArgoCD setup stateless, meaning we can upgrade ArgoCD by toggling the integration off and on again in
    the `environment.yaml` manifest.

### Namespaces

The namespaces ArgoCD application tracks the `/infrastructure/<environment name>/argocd/namespaces/` directory. This
directory contains Kubernetes manifests defining all namespaces in the cluster. This lets us disassociate a namespace
from an application, allowing namespaces to outlive an application in the case where multiple applications reside in the
same namespace, but one is deleted.