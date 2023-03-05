data "transparentedge_staging_backends" "all" {}

output "all_staging_backends" {
  value = data.transparentedge_staging_backends.all
}
