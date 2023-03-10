resource "transparentedge_staging_backend" "stagorigin1" {
  name   = "stagorigin1"
  origin = "origin.example.com"
  port   = 443
  ssl    = true

  # healthcheck
  hchost       = "www.origin.example.com"
  hcpath       = "/favicon.ico"
  hcstatuscode = 200
}

resource "transparentedge_staging_backend" "stagorigin2" {
  name   = "stagorigin2"
  origin = "origin2.example.com"
  port   = 80
  ssl    = false

  # healthcheck
  hchost       = "www.origin2.example.com"
  hcpath       = "/favicon.ico"
  hcstatuscode = 403
}

output "origin1" {
  value = transparentedge_staging_backend.stagorigin1
}

output "origin2" {
  value = transparentedge_staging_backend.stagorigin2
}
