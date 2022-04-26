# Deploy an application

First, scaffold an application configuration file by running the following command:

```shell
xctl scaffold application > application.yaml
```

Configure the application as required, then run the following command:

```shell
xctl apply -f application.yaml
```

After you've configured the generated files, commit and push the changes.
