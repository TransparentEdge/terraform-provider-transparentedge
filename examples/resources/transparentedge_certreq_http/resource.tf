terraform {
  required_providers {
    transparentedge = {
      source = "TransparentEdge/transparentedge"
      # Available since version 0.5.0
      version = ">=0.5.0"
    }
  }
}

resource "transparentedge_certreq_http" "cr_example_com" {
  domains    = ["www.example.com", "static.example.com"]
  standalone = true
}

data "transparentedge_certreq_http" "cr_example_com" {
  id = transparentedge_certreq_http.cr_example_com.id
}

output "cr_example_com" {
  value = data.transparentedge_certreq_http.cr_example_com
}
