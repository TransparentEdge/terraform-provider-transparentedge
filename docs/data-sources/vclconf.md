---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "transparentedge_vclconf Data Source - transparentedge"
subcategory: ""
description: |-
  VCL Configuration listing.
---

# transparentedge_vclconf (Data Source)

VCL Configuration listing.

## Example Usage

```terraform
data "transparentedge_vclconf" "vclconfig" {}

output "active_vcl_config" {
  value = data.transparentedge_vclconf.vclconfig
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `company` (Number) Company ID that owns this VCL config.
- `id` (Number) ID of the VCL Config.
- `productiondate` (String) Date when the configuration was fully applied in the CDN.
- `uploaddate` (String) Date when the configuration was uploaded.
- `user` (String) User that created the configuration.
- `vclcode` (String) Verbatim of the VCL code.
