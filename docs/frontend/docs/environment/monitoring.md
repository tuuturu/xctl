

## Accessing Grafana

First, acquire credentials for Grafana by running the `get credentials` command:

```shell
xctl --context env.yaml get credentials grafana
```

After you've gotten the credentials, open a tunnel to the Grafana instance by running the following command:

```shell
xctl --context env.yaml forward grafana
```

## Viewing logs

To view logs, first click on the explore button in the menu on the left side of the screen. Make sure you select "Loki"
as the datasource.

Follow the steps in the query builder and click "Show logs".

Check out the [LogQL](https://grafana.com/docs/loki/latest/logql/log_queries/) documentation for more information on
how to use the log query interface.

## Viewing metrics

To view metrics, first click on the explore button in the menu on the left side of the screen. Make sure you select
"Prometheus" as the datasource.

To find certain metrics you can click the dropdown button to the left of the query field to select relevant metrics.

Check out the [LogQL](https://grafana.com/docs/loki/latest/logql/metric_queries/) documentation for more information on
how to use the metric query interface.
