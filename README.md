# Xctl - A tool for managing a Kubernetes based infrastructure

Inspired by the tool [okctl](https://github.com/oslokommune/okctl) made by Oslo Kommune in Norway.

## Vision

To be the quickest way to provision a complete and familiar infrastructure for your production needs.

## Supported cloud providers:
- [x] Linode

## What does xctl provision?

- [x] Kubernetes as the platform
- [ ] Inbound traffic
    - [ ] Nginx Ingress Controller for routing traffic
    - [ ] Certbot for your TLS needs
    - [ ] ExternalDNS for automatic DNS handling
- [ ] Secrets
    - [ ] Vault as the secret manager
    - [ ] External Secrets for accessing secrets
- [ ] Monitoring
    - [ ] Prometheus for scraping metrics
    - [ ] Promtail for scraping container logs
    - [ ] Grafana for visualizing logs and metrics
- [ ] CI/CD
    - [ ] ArgoCD
