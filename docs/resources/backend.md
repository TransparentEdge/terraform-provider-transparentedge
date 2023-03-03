---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "transparentedge_backend Resource - transparentedge"
subcategory: ""
description: |-
  
---

# transparentedge_backend (Resource)



## Example Usage

```terraform
resource "transparentedge_backend" "origin1" {
  name   = "origin1"
  origin = "my-origin.com"
  port   = 443
  ssl    = true

  # healthcheck
  hchost       = "www.my-origin.com"
  hcpath       = "/favicon.ico"
  hcstatuscode = 200
}

resource "transparentedge_backend" "origin2" {
  name   = "origin2"
  origin = "my-origin2.com"
  port   = 80
  ssl    = false

  # healthcheck
  hchost       = "www.my-origin2.com"
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

- `hchost` (String) Host header that the healthcheck probe will send to the origin, for example: www.my-origin.com
- `hcpath` (String) Path that the healthcheck probe will used, for example: /favicon.ico
- `hcstatuscode` (Number) Status code expected when the probe receives the HTTP healthcheck response, for example: 200
- `name` (String) Name of the backend
- `origin` (String) Origin is the IP or DNS address to the origin backend, for example: 'my-origin.com'
- `port` (Number) Port where the origin is listening to HTTP requests, for example: 80 or 443
- `ssl` (Boolean) If the origin should be contacted using TLS encription.

### Read-Only

- `company` (Number) Company ID that owns this backend
- `id` (Number) ID of the backend

## Import

Import is supported using the following syntax:

```shell
# Import a backend by its name
terraform import 'transparentedge_backend.origin1' 'origin1'
```