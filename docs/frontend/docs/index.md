# Xctl - A tool for managing a Kubernetes based infrastructure

<div class="buttons">
    <a class="md-button" href="/getting-started/apply-environment/">Get started</a>
</div>

XCTL is a tool that leverages [CNCF technology](https://www.cncf.io/) to bootstrap a production grade environment.

!!! danger "Warning"
    This software is not ready for actual production usage

```bash
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

<div class="technology">
  <img class="icon" src="/img/kubernetes-icon.svg">
  <span class="content">
    <a href="https://kubernetes.io/">Kubernetes</a> for resource orchestration
  </span>
</div>

<div class="technology">
  <img class="icon" src="/img/nginx-icon.svg">
  <span class="content">
    <a href="https://kubernetes.github.io/ingress-nginx/">NGINX Ingress Controller</a> for traffic routing
  </span>
</div>

<div class="technology">
  <img class="icon" src="/img/certmanager-icon.png">
  <span class="content">
    <a href="https://cert-manager.io/">CertManager</a>
    configured with <a href="https://letsencrypt.org/">Let's encrypt</a> for TLS
  </span>
</div>

<div class="technology">
  <img class="icon" src="/img/grafana-icon.png">
  <span class="content">
    <a href="https://grafana.com/oss/grafana/">Grafana</a>,
    <a href="https://grafana.com/oss/prometheus/">Prometheus</a>,
    <a href="https://grafana.com/oss/loki/">Loki</a> and
    <a href="https://grafana.com/docs/loki/latest/clients/promtail/">Promtail</a> for monitoring
  </span>
</div>

<div class="technology">
  <img class="icon" src="/img/argocd-icon.png">
  <span class="content">
    <a href="https://argoproj.github.io/cd/">ArgoCD</a>
    for continuous deployment (TBA)
  </span>
</div>

<style>
div.buttons {
    width: 100%;
    display: flex;

    justify-content: center;

    margin-bottom: 2em;
}

div.technology {
  display: flex;
  align-items: center;

  margin-top: 1em;
  padding-left: 1em;
}

img.icon {
  max-width: 48px;
  max-height: 48px;
  min-height: 48px;
}

span.content {
  margin-left: 1em;
}
</style>

