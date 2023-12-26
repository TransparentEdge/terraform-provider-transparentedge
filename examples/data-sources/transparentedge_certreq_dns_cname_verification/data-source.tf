terraform {
  required_providers {
    transparentedge = {
      source = "TransparentEdge/transparentedge"
      # Available since version 0.5.0
      version = ">=0.5.0"
    }
  }
}

data "transparentedge_certreq_dns_cname_verification" "dns_cname" {}

output "cname" {
  value = data.transparentedge_certreq_dns_cname_verification.dns_cname.cname
}
