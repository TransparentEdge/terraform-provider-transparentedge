terraform {
  required_providers {
    transparentedge = {
      version = ">= 0.1.0"
    }
  }
}

provider "transparentedge" {}

data "transparentedge_sites" "all" {}

output "all_sites" {
  value = data.transparentedge_sites.all
}
