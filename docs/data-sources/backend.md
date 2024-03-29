---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "transparentedge_backend Data Source - transparentedge"
subcategory: ""
description: |-
  Read a backend.
---

# transparentedge_backend (Data Source)

Read a backend.

## Example Usage

```terraform
terraform {
  required_providers {
    transparentedge = {
      source = "TransparentEdge/transparentedge"
      # Available since version 0.3.0
      version = ">=0.3.0"
    }
  }
}

data "transparentedge_backend" "mybackend" {
  name = "mybackendname"
}

output "vclname" {
  # use 'vclname' to associate a backend in VCL Code
  # for example: set req.backend_hint = ${data.transparentedge_backend.mybackend.vclname}.backend();
  value = data.transparentedge_backend.mybackend.vclname
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the backend.

### Read-Only

- `company` (Number) Company ID that owns this backend.
- `hchost` (String) Host header that the healthcheck probe will send to the origin, for example: www.my-origin.com.
- `hcpath` (String) Host header that the healthcheck probe will send to the origin, for example: `www.my-origin.com`.
- `hcstatuscode` (Number) Status code expected when the probe receives the HTTP healthcheck response, for example: `200`.
- `id` (Number) ID of the backend.
- `origin` (String) IP or DNS name pointing to the origin backend, for example: `my-origin.com`.
- `port` (Number) Port where the origin is listening to HTTP requests, for example: `80` or `443`.
- `ssl` (Boolean) Use TLS encription when contacting with the origin backend.
- `vclname` (String) Final unique name of the backend to be referenced in VCL Code: `c{company_id}_{name}`.
