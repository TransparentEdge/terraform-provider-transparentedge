terraform {
  required_providers {
    transparentedge = {
      source = "TransparentEdge/transparentedge"
      # Available since version 0.5.0
      version = ">=0.5.0"
    }
  }
}

data "transparentedge_certreq_http" "http_cr" {
  # ID of the HTTP Certificate Request
  id = 1056
}

output "http_cr" {
  value = data.transparentedge_certreq_http.http_cr
}
