---
layout: ""
page_title: "Full configuration example"
description: |-
    A demonstration using a sample configuration that leverages on all the resources from this provider.
---

## Full configuration example

First, we create the file `provider.tf`:  

```terraform
terraform {
  required_providers {
    transparentedge = {
      version = ">= 0.2.0"
    }
  }
}

provider "transparentedge" {
  company_id = 300
}
```

In this example we are provisioning two sites: `www.example1.com` and `www.example.com`.  
On another file: `sites.tf`:  

```terraform
variable "sites" {
  type = set(any)
  default = [
    "www.example1.com",
    "www.example2.com",
  ]
}

# Verification string is required for new sites that are not already owned
# please check the documentation or contact with support
data "transparentedge_siteverify" "all" {
  for_each = var.sites
  domain   = each.key
}

# Show all the verification strings
output "verification_strings" {
  value = data.transparentedge_siteverify.all
}

# Provision the sites
resource "transparentedge_site" "all" {
  for_each = var.sites
  domain   = each.key

  # Timeout for site creation, only useful if the site is not verified
  # site creation will be retried until timeout or error
  timeouts = {
    create = "120s"
  }
}
```

For each site we output the verification string, when you add new sites to the CDN a verification process is required.   
If you don't known the steps please login into our [dashboard](https://dashboard.transparentcdn.com/) and try to add a new site, everything is documented there.  

Then we create three new files a terraform file named `certificates.tf` and two files that contain the public and private key in PEM format: `mysite1.crt` and `mysite1.key`:  

```terraform
resource "transparentedge_custom_certificate" "mysite1" {
  publickey  = file("${path.module}/mysite1.crt")
  privatekey = file("${path.module}/mysite1.key")
}
```

A file named `backends.tf`:  

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
```

And finally, a file named `vclconfig.tf` with the example code:  

```terraform
resource "transparentedge_vclconf" "myconfig" {
  vclcode = <<EOF
sub vcl_recv {
    if (req.http.host ~ "www.example(1|2).com") {
        set req.backend_hint = ${resource.transparentedge_backend.origin1.vclname}.backend();
        set req.http.TCDN-i3-transform = "auto_webp";
        set bereq.http.TCDN-Command = "redirect_https, brotli_compress";
    }
}
```

