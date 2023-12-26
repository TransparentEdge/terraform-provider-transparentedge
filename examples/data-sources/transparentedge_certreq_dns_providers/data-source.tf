terraform {
  required_providers {
    transparentedge = {
      source = "TransparentEdge/transparentedge"
      # Available since version 0.5.0
      version = ">=0.5.0"
    }
  }
}

data "transparentedge_certreq_dns_providers" "all" {}

output "providers" {
  value = data.transparentedge_certreq_dns_providers.all
}

# This example only outputs the required parameters for the 'AWS (Route53)' DNS provider
output "aws_dns_provider_parameters" {
  value = [
    for dns_provider in data.transparentedge_certreq_dns_providers.all.providers :
    dns_provider.parameters
    if dns_provider.dns_provider == "AWS (Route53)"
  ]
}
