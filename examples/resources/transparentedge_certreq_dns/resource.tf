terraform {
  required_providers {
    transparentedge = {
      source = "TransparentEdge/transparentedge"
      # Available since version 0.5.0
      version = ">=0.5.0"
    }
  }
}

resource "transparentedge_certreq_dns" "dns_cr_example_com" {
  credential = 30
  domains    = ["example.com", "*.example.com"]
}

data "transparentedge_certreq_dns" "dns_cr_example_com" {
  id = transparentedge_certreq_dns.dns_cr_example_com.id
}

output "example_com" {
  value = data.transparentedge_certreq_dns.dns_cr_example_com
}
