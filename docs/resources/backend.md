---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "transparentedge_backend Resource - transparentedge"
subcategory: ""
description: |-
  Provides a Backend resource. This allows backends to be created, updated and deleted.
---

# transparentedge_backend (Resource)

Provides a Backend resource. This allows backends to be created, updated and deleted.

## Example Usage

```terraform
resource "transparentedge_backend" "origin1" {
  name   = "origin1"
  origin = "origin.example.com"
  port   = 443
  ssl    = true

  # healthcheck
  hchost       = "www.origin.example.com"
  hcpath       = "/favicon.ico"
  hcstatuscode = 200
}

resource "transparentedge_backend" "origin2" {
  name   = "origin2"
  origin = "origin2.example.com"
  port   = 80
  ssl    = false

  # healthcheck
  hchost       = "www.origin2.example.com"
  hcpath       = "/favicon.ico"
  hcstatuscode = 403
}

output "origin1" {
  value = transparentedge_backend.origin1
}

output "origin2" {
  value = transparentedge_backend.origin2
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `hchost` (String) Host header that the healthcheck probe will send to the origin, for example: `www.my-origin.com`.
- `hcpath` (String) Path that the healthcheck probe will use, for example: `/favicon.ico`.
- `hcstatuscode` (Number) Status code expected when the probe receives the HTTP healthcheck response, for example: `200`.
- `name` (String) Name of the backend.
- `origin` (String) IP or DNS name pointing to the origin backend, for example: `my-origin.com`.
- `port` (Number) Port where the origin is listening to HTTP requests, for example: `80` or `443`.
- `ssl` (Boolean) Use TLS encription when contacting with the origin backend.

### Read-Only

- `company` (Number) Company ID that owns this backend.
- `id` (Number) ID of the backend.
- `vclname` (String) Final unique name of the backend to be referenced in VCL Code: `c{company_id}_{name}`.

## Import

Import is supported using the following syntax:

```shell
# Import a backend by name
terraform import 'transparentedge_backend.origin1' 'origin1'
```
