# Xctl - A tool for managing a Kubernetes based infrastructure

XCTL is a tool that leverages [CNCF technology](https://www.cncf.io/) to bootstrap a production grade environment.

⚠ Warning: This software is not ready for actual production usage ⚠

```shell
xctl apply -f - << EOF
apiVersion: v1alpha1
kind: Environment

metadata:
  name: demo
  email: demo@example.com

spec:
  domain: example.com
EOF
```

will result in the following resources:

* [Kubernetes](https://kubernetes.io/) for resource orchestration
* [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/) for routing
* [CertManager](https://cert-manager.io/) configured with [Let's encrypt](https://letsencrypt.org/) for TLS
* [Grafana](https://grafana.com/oss/grafana/), [Prometheus](https://grafana.com/oss/prometheus/),
  [Loki](https://grafana.com/oss/loki/) and [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/) for
  monitoring
* [ArgoCD](https://argoproj.github.io/cd/) for continuous deployment (TBA)

<div class="buttons">
    <a class="md-button" href="/getting-started/apply-environment/">Get started</a>
</div>

<style>
div.buttons {
    width: 100%;
    display: flex;

    justify-content: center;
}
</style>

