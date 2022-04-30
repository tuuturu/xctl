# Deploy an application

## Configure

First, scaffold an application configuration file by running the following command:

```shell
xctl scaffold application > application.yaml
```

Configure the application as required.

!!! note
    A full list of available configuration parameters can be found [here](/application/manifest)

## Deploy

Generate necessary Kubernetes and ArgoCD configuration by running the following command:

```shell
xctl apply -f application.yaml
```

After you've configured the generated files, commit and push the changes.
