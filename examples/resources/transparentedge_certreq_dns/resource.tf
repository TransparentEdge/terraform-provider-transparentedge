terraform {
  required_providers {
    transparentedge = {
      source = "TransparentEdge/transparentedge"
      # Available since version 0.5.0
      version = ">=0.5.0"
    }
  }
}

# This resource requires a credential
# Example for 'AWS (Route53)'
resource "transparentedge_certreq_dns_credential" "aws_01" {
  alias = "aws-creds-01"

  parameters = {
    AWS_SECRET_ACCESS_KEY = "ABC"
    AWS_ACCESS_KEY_ID     = "DEF"
  }
}

resource "transparentedge_certreq_dns" "dns_cr_example_com" {
  credential = transparentedge_certreq_dns_credential.aws_01.id

  # DNS Certificate Request can generated wildcard certificates
  domains = ["example.com", "*.example.com"]
}

data "transparentedge_certreq_dns" "dns_cr_example_com" {
  id = transparentedge_certreq_dns.dns_cr_example_com.id
}

output "example_com" {
  value = data.transparentedge_certreq_dns.dns_cr_example_com
}
