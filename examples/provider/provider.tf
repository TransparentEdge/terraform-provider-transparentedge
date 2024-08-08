terraform {
  required_providers {
    transparentedge = {
      source  = "TransparentEdge/transparentedge"
      version = ">=0.6.0"
    }
  }
}

# Configure the provider
provider "transparentedge" {
  company_id = 300
}

# Create a backend
resource "transparentedge_backend" "origin1" {
  name   = "origin1"
  origin = "origin.example.com"
  port   = 443
  ssl    = true

  # health check
  hchost       = "www.origin.example.com"
  hcpath       = "/favicon.ico"
  hcstatuscode = 200
}
