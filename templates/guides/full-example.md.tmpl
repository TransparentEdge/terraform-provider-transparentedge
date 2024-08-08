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
      source = "TransparentEdge/transparentedge"
      version = ">=0.6.0"
    }
  }
}

provider "transparentedge" {
  company_id = 300
}
```

In this example we are provisioning two sites: `www.example1.com` and `www.example2.com`, create the file: `sites.tf`:  

```terraform
variable "sites" {
  type = set(any)
  default = [
    "www.example1.com",
    "www.example2.com",
  ]
}

# Verification string is required for new sites that are not already owned
# in case of doubts please check the documentation or contact with support
data "transparentedge_siteverify" "all" {
  for_each = var.sites
  domain   = each.key
}

# Show the verification strings for each domain
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

For each site we output the verification string. When you add new sites to the CDN a verification process is required to ensure that you own the domain.   
If you don't know the steps please login into our [dashboard](https://dashboard.transparentcdn.com/) and try to add a new site, everything is documented in the process.  

Then we create three new files: a terraform file named `certificates.tf` and two files that contain the public and private key in `PEM` format: `mysite1.crt` and `mysite1.key`:  

```terraform
resource "transparentedge_custom_certificate" "mysite1" {
  publickey  = file("${path.module}/mysite1.crt")
  privatekey = file("${path.module}/mysite1.key")
}
```

Now with the backend definition, we create a file named `backends.tf`:  

```terraform
resource "transparentedge_backend" "origin1" {
  name   = "origin1"
  origin = "origin.example.com"
  port   = 443
  ssl    = true

  # health check
  hchost       = "www.origin.example.com"
  hcpath       = "/favicon.ico"
  hcstatuscode = 200
  hcinterval   = 40
  hcdisabled   = false
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
        set req.http.TCDN-Command = "brotli_compress";
        call redirect_https;
    }
}
EOF
}
```

