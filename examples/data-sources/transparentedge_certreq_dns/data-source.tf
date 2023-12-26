terraform {
  required_providers {
    transparentedge = {
      source = "TransparentEdge/transparentedge"
      # Available since version 0.5.0
      version = ">=0.5.0"
    }
  }
}

data "transparentedge_certreq_dns" "dns_cr" {
  id = 144
}

output "cred" {
  value = data.transparentedge_certreq_dns.dns_cr
}
