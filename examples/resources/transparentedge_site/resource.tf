################################
### Multiple sites in a list ###
################################
variable "sites" {
  type = set(any)
  default = [
    "www.example1.com",
    "www.example2.com",
  ]
}

# Verification string is required for new sites that are not already owned
# in case of doubts please check the documentation or contact with support
data "transparentedge_siteverify" "all" {
  for_each = var.sites
  domain   = each.key
}

output "verification_strings" {
  value = data.transparentedge_siteverify.all
}

resource "transparentedge_site" "all" {
  for_each = var.sites
  domain   = each.key

  # Timeout for site creation, only useful if the site is not verified
  # site creation will be retried until timeout or error
  timeouts = {
    create = "120s"
  }
}

output "all_sites" {
  value = transparentedge_site.all
}


###################
### Single site ###
###################
data "transparentedge_siteverify" "example3" {
  domain = "www.example3.com"
}

output "example3_verifycation" {
  value = data.transparentedge_siteverify.example3
}

resource "transparentedge_site" "www_example3_com" {
  domain = "www.example3.com"
}
