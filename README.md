# Sentinel Terraform Provider
[![CI, Validation, and Conditional Release](https://github.com/ayuspoudel/terraform-provider-sentinel/actions/workflows/release.yml/badge.svg)](https://github.com/ayuspoudel/terraform-provider-sentinel/actions/workflows/release.yml)

The **Sentinel Terraform Provider** allows you to manage **Sentinel SRE policies** using Terraform.
It integrates directly with the Sentinel Control Plane API to apply, update, delete, and observe
policy definitions as infrastructure-as-code.

This provider is designed to be declarative, idempotent, and safe to use in CI/CD pipelines.




## Requirements

- Terraform >= 1.3
- A running Sentinel Control Plane
- Network access to the Sentinel API endpoint



## Installation

```hcl
terraform {
  required_providers {
    sentinel = {
      source  = "ayuspoudel/sentinel"
      version = "~> 0.1.0"
    }
  }
}
````



## Provider Configuration

```hcl
provider "sentinel" {
  endpoint = "http://localhost:8080"
  token    = "optional-api-token"
}
```

### Arguments

| Name     | Type   | Required | Description                            |
| -- |  | -- | -- |
| endpoint | string | yes      | Base URL of the Sentinel Control Plane |
| token    | string | no       | Bearer token used for authentication   |



## Resources

### `sentinel_policy`

Manages a Sentinel policy definition.
Policies are **idempotent** — updates and creates use the same API semantics.

#### Example

```hcl
resource "sentinel_policy" "example" {
  name = "checkout"

  metadata {
    owner       = "team-a"
    environment = "prod"
  }

  target {
    cluster   = "sreCluster"
    namespace = "default"
    service   = "checkout"
  }

  signals {
    traffic {
      query   = "sum(rate(http_requests_total[1m]))"
      min_rps = 10
    }

    errors {
      query = "sum(rate(http_requests_errors_total[1m]))"
    }

    slo {
      objective = 99.9
      window    = "720h"
    }
  }

  policy {
    budget {
      fast_burn {
        window    = "5m"
        threshold = 2.0
      }

      slow_burn {
        window    = "1h"
        threshold = 1.0
      }
    }
  }
}
```



## Data Sources

### `sentinel_policy_status`

Reads the **evaluated status** of an existing Sentinel policy.

This data source is useful for:

* Observability dashboards
* CI/CD gating
* Validating policy readiness

#### Example

```hcl
data "sentinel_policy_status" "checkout" {
  name = "checkout"
}
```

#### Attributes

| Name                 | Description                              |
| -- | - |
| cluster_exists       | Whether the target cluster is registered |
| cluster_reachable    | Whether the cluster API is reachable     |
| namespace_exists     | Whether the namespace exists             |
| agent_installed      | Whether Sentinel agent is installed      |
| agent_healthy        | Whether Sentinel agent is healthy        |
| prometheus_reachable | Whether Prometheus is reachable          |
| queries_valid        | Whether signal queries are valid         |
| last_validated_at    | Last validation timestamp                |
| last_error           | Last validation error, if any            |



## Import

Existing policies can be imported into Terraform state.

```bash
terraform import sentinel_policy.example checkout
```



## Development

### Running Tests

```bash
go test ./internal/...
```

### Building the Provider

```bash
go build -o terraform-provider-sentinel
```
 


## Versioning

This provider follows semantic versioning.

* `0.x` releases may introduce breaking changes
* `1.x` releases will maintain backward compatibility



## License

MIT


