data "transparentedge_stagingbackends" "all" {}

output "all_staging_backends" {
  value = data.transparentedge_stagingbackends.all
}
