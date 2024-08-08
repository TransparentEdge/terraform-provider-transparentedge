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
}

resource "transparentedge_backend" "origin2" {
  name   = "origin2"
  origin = "origin2.example.com"
  port   = 80
  ssl    = false

  # health check (disabled) - hchost, hcpath, hcstatuscode are still required
  hchost       = "www.origin2.example.com"
  hcpath       = "/favicon.ico"
  hcstatuscode = 403
  hcdisabled   = true # the health check probe is disabled
}

output "origin1" {
  value = transparentedge_backend.origin1
}

output "origin2" {
  value = transparentedge_backend.origin2
}
