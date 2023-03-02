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

resource "transparentedge_backend" "origin2" {
  name   = "origin2"
  origin = "my-origin2.com"
  port   = 80
  ssl    = false

  # healthcheck
  hchost       = "www.my-origin2.com"
  hcpath       = "/favicon.ico"
  hcstatuscode = 403
}

output "origin1" {
  value = transparentedge_backend.origin1
}

output "origin2" {
  value = transparentedge_backend.origin2
}
