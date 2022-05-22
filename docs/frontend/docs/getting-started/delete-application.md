# Delete an application

To delete an application, run the following command:

```shell
xctl delete --file application.yaml --context environment.yaml
```

Commit and push to finalize the changes.

!!! tip
    To avoid having to specify the environment context for contextual commands, use `xctl venv`
