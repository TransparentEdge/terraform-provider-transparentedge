---
layout: ""
page_title: "TransparentEdge Terraform Provider"
description: |-
    The TransparentEdge provider to interact directly with the autoprovisioning
    capabilities of the CDN via API.
---

# TransparentEdge Terraform Provider

A terraform provider for the CDN of [Transparent Edge](https://www.transparentedge.eu/).
This provider is intended to be used by the CDN users with the role "Company Admin", although some services do not require this role.

## Example usage

```terraform
terraform {
  required_providers {
    transparentedge = {
      version = ">= 0.1.0"
    }
  }
}

provider "transparentedge" {}
```

It's recommended to use environment variables:  

```shell
export TCDN_COMPANY_ID=0
export TCDN_CLIENT_ID="xxx"
export TCDN_CLIENT_SECRET="xxx"
```

You can find all the required variables in our [our dashboard](https://dashboard.transparentcdn.com/).
Login, go to your profile -> "Account options" -> "Manage keys".

Make sure that you're using the correct Company ID if you own multiple companies.  

Optional environment variables:  

```shell
export TCDN_HOST_URL="https://api.transparentcdn.com"
export TCDN_VERIFY_SSL=true
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `client_id` (String, Sensitive) Client ID (dashboard -> profile -> account options -> manage keys). May also be provided via TCDN_CLIENT_ID environment variable.
- `client_secret` (String, Sensitive) Client Secret (dashboard -> profile -> account options -> manage keys). May also be provided via TCDN_CLIENT_SECRET environment variable.
- `company_id` (Number) Company ID number (for ex: 300). May also be provided via TCDN_COMPANY_ID environment variable.
- `host_url` (String) URL of Transparent Edge API. default: 'https://api.transparentcdn.com'. May also be provided via TCDN_HOST_URL environment variable.
- `verify_ssl` (Boolean) Ignore SSL certificate for 'host_url'. May also be provided via TCDN_VERIFY_SSL environment variable.