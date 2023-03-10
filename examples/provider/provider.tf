terraform {
  required_providers {
    transparentedge = {
      source  = "TransparentEdge/transparentedge"
      version = ">=0.3.3"
    }
  }
}

provider "transparentedge" {
  # Provider configuration overrides environment variables
  # it's recommended to use environment variables for company_id, client_id and client_secret
  company_id    = 300
  client_id     = "XXX"
  client_secret = "XXX"
  insecure      = false                            # this is the default value
  api_url       = "https://api.transparentcdn.com" # this is the default value
  auth          = true                             # false if you only use data-sources that do not require authentication such as 'transparentedge_ip_ranges'
}
