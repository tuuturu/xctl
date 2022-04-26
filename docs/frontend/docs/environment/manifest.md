# Environment manifest

Below is a list of all configuration attributes for an environment manifest.

## `apiVersion` <span class="required">required</span>

| Type   | Example  | Description                        |
|--------|----------|------------------------------------|
| string | v1alpha1 | Defines the version of the schema. |

## `kind` <span class="required">required</span>

| Type   | Example     | Description                     |
|--------|-------------|---------------------------------|
| string | Environment | Defines the kind of the schema. |

## `metadata`

### `name` <span class="required">required</span>

| Type   | Regex                  | Example            | Description                          |
|--------|------------------------|--------------------|--------------------------------------|
| string | `[a-z]+(-[a-z]+){0,2}` | tuuturu-production | Defines the name of the environment. |

### `email` <span class="required">required</span>

| Type   | Example          | Description                                                                                                                  |
|--------|------------------|------------------------------------------------------------------------------------------------------------------------------|
| string | xctl@example.com | Defines the email associated with the environment. Currently, only used for registering SSL certificates with Let's Encrypt. |

## `spec`  <span class="required">required</span>

### `domain` <span class="required">required</span>

| Type   | Example     | Description                                       |
|--------|-------------|---------------------------------------------------|
| string | tuuturu.org | Defines the associated domain for the environment |

### `plugins`

#### `nginxIngressController`

| Type    | Default | Description                          |
|---------|---------|--------------------------------------|
| boolean | true    | Handles traffic into the environment |

#### `certBot`

| Type    | Default | Description                          |
|---------|---------|--------------------------------------|
| boolean | true    | Generates SSL certificates on demand |

#### `grafana`

| Type    | Default | Description                      |
|---------|---------|----------------------------------|
| boolean | true    | Visualises metrics, logs, traces |

#### `prometheus`

| Type    | Default | Description                     |
|---------|---------|---------------------------------|
| boolean | true    | Stores and exposes metrics data |

#### `loki`

| Type    | Default | Description                 |
|---------|---------|-----------------------------|
| boolean | true    | Stores and exposes log data |

#### `promtail`

| Type    | Default | Description                          |
|---------|---------|--------------------------------------|
| boolean | true    | Scrapes logs and pushes them to Loki |

<style>
span.required {
    color: red;
}
</style>
