# Environment manifest

Below is a list of all configuration attributes for an environment manifest.

## `apiVersion`<span class="required">*</span>

<span class="required">required</span>

| Type   | Example  | Description                        |
|--------|----------|------------------------------------|
| string | v1alpha1 | Defines the version of the schema. |

## `kind`<span class="required">*</span>

<span class="required">required</span>

| Type   | Example     | Description                     |
|--------|-------------|---------------------------------|
| string | Application | Defines the kind of the schema. |

## `metadata`

### `name`<span class="required">*</span>

<span class="required">required</span>

| Type   | Regex                  | Example     | Description                          |
|--------|------------------------|-------------|--------------------------------------|
| string | `[a-z]+(-[a-z]+){0,2}` | hello-world | Defines the name of the application. |

## `spec`

### `image`<span class="required">*</span>

<span class="required">required</span>

| Type   | Example                    | Description                                                           |
|--------|----------------------------|-----------------------------------------------------------------------|
| string | ghcr.io/tuuturu/xctl-hello | Defines the URI of the container image to be used in the application. |

### `port`

| Type   | Example | Default | Description                                                                                                |
|--------|---------|---------|------------------------------------------------------------------------------------------------------------|
| number | 80      | None    | Defines the port your application listens to traffic on. Only required if the application expects traffic. |

### `url`

| Type   | Example           | Default | Description                                                                                                                      |
|--------|-------------------|---------|----------------------------------------------------------------------------------------------------------------------------------|
| string | hello.tuuturu.org | None    | Defines the URL your application will be available on. Only required if the application should be available outside the cluster. |

<style>
span.required {
    color: red;
}
</style>
