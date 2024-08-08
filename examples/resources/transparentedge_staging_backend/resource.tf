resource "transparentedge_staging_backend" "stagorigin1" {
  name   = "stagorigin1"
  origin = "origin.example.com"
  port   = 443
  ssl    = true

  # health check
  hchost       = "www.origin.example.com"
  hcpath       = "/favicon.ico"
  hcstatuscode = 200
  hcinterval   = 40
}

resource "transparentedge_staging_backend" "stagorigin2" {
  name   = "stagorigin2"
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
  value = transparentedge_staging_backend.stagorigin1
}

output "origin2" {
  value = transparentedge_staging_backend.stagorigin2
}
