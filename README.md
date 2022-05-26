# XCTL - A tool for managing a Kubernetes based infrastructure

Inspired by the tool [okctl](https://github.com/oslokommune/okctl) made by Oslo Kommune in Norway.

## Vision

### Mission

To be a simple and familiar way to handle a complete infrastructure for your production needs.

### Hypothesis

By exposing an opinionated set of good enough technologies in an intuitive matter, developers will have more time to
focus on delivering real value. This means:
- One way of orchestrating frontend and backend services, Kubernetes.
- One way of handling jobs, Kubernetes.
- One way of doing queues, TBA.
- One relational database technology, PostgreSQL.
- One noSQL database technology, TBA.

## Supported cloud providers:

- [x] Linode

## Installation

See [here](https://xctl.tuuturu.org/getting-started/preparation/)

## Usage

### Authenticate

First, to authenticate, export the environment variable `LINODE_TOKEN` containing your personal access token.

### Provision an environment

```shell
xctl apply -f - << EOF
apiVersion: v1alpha1
kind: Environment

metadata:
  name: demo
  email: demo@example.com

spec:
  provider: linode
  domain: example.com
  repository: git@github.com:tuuturu/iac.git
EOF
```

After a few minutes you'll have a running Kubernetes cluster with the technologies listed
[here](#what-does-xctl-provision).

### Administrating your environment

To be able to run `kubectl` commands, use `xctl venv -c environment.yaml` to create a subshell with the environment
variable `KUBECONFIG` set.

### Deploying an app

```shell
xctl apply -f - << EOF
apiVersion: v1alpha1
kind: Application

metadata:
  name: hello

spec:
  image: ghcr.io/tuuturu/xctl-hello
  port: 80
  url: hello.tuuturu.org
EOF
```

After a committing and pushing the changes done by `xctl`, ArgoCD should soon spin up your application.

## What does xctl provision?

- [x] Kubernetes as the platform
- [x] Inbound traffic
  - [x] Nginx Ingress Controller for routing traffic
  - [x] Certbot for your TLS needs
- [x] Monitoring
  - [x] Grafana for visualizing logs and metrics
  - [x] Prometheus for scraping metrics
  - [x] Loki for collecting logs and making them queryable
  - [x] Promtail for scraping container logs
- [x] CI/CD
  - [x] ArgoCD for continuous deployment
- [ ] Secrets
  - [ ] Vault as the secret manager
  - [ ] External Secrets for accessing secrets
