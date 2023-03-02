terraform {
  required_providers {
    transparentedge = {
      version = ">= 0.1.0"
    }
  }
}

provider "transparentedge" {}

data "transparentedge_backends" "all" {}

output "all_backends" {
  value = data.transparentedge_backends.all
}
