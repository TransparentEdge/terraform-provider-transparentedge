data "transparentedge_backends" "all" {}

output "all_backends" {
  value = data.transparentedge_backends.all
}
