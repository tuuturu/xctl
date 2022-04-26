# Create an environment

First, scaffold an environment configuration file by running the following command:

```shell
xctl scaffold environment > environment.yaml
```

Configure the environment as required, then run the following command:

```shell
xctl apply -f environment.yaml
```

After a few minutes, you'll have a production grade environment ready to be used.