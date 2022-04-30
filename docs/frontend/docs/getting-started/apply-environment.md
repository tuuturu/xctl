# Create an environment

## Configure

First, scaffold an environment configuration file by running the following command:

```shell
xctl scaffold environment > environment.yaml
```

Configure the environment as required.

!!! note
    A full list of available configuration parameters can be found [here](/environment/manifest)

## Create

To create an environment based on the configuration file, run the following command:

```shell
xctl apply -f environment.yaml
```

After a few minutes, you'll have a production grade environment ready to be used.