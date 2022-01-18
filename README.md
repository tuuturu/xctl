# Xctl - A tool for managing a Kubernetes based infrastructure

Inspired by the tool [okctl](https://github.com/oslokommune/okctl) made by Oslo Kommune in Norway.

## Vision

To be the quickest way to provision a complete and familiar infrastructure for your production needs.

## Supported cloud providers:

- [x] Linode

## Usage

### Authenticate

First, to authenticate, export the environment variable `LINODE_TOKEN` containing your personal access token.

### Provision a cluster

1. Scaffold a cluster manifest template by running `xctl scaffold cluster > cluster.yaml`.
2. Configure the scaffolded `cluster.yaml` template.
3. Apply your manifest with `xctl apply -f cluster.yaml`.

After a few minutes you'll have a running cluster with the integrations listed [here](#what-does-xctl-provision)

### Administrating your cluster

To be able to run `kubectl` commands, use `xctl venv -c cluster.yaml` to create a subshell with the environment variable
`KUBECONFIG` set.

### Deploying an app (To be implemented)

1. Scaffold an application manifest template by running `xctl scaffold application > app.yaml`.
2. Configure the scaffolded `app.yaml` template.
3. Apply your manifest with `xctl apply -f app.yaml`

After a committing and pushing the changes done by `xctl`, ArgoCD should soon spin up your application.

## What does xctl provision?

- [x] Kubernetes as the platform
- [x] Inbound traffic
  - [x] Nginx Ingress Controller for routing traffic
  - [x] Certbot for your TLS needs
- [ ] Secrets
  - [x] Vault as the secret manager
  - [ ] External Secrets for accessing secrets
- [ ] Monitoring
    - [ ] Prometheus for scraping metrics
    - [ ] Promtail for scraping container logs
    - [ ] Grafana for visualizing logs and metrics
- [ ] CI/CD
    - [ ] ArgoCD for continuous deployment
