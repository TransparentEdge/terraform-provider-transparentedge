data "transparentedge_sites" "all" {}

output "all_sites" {
  value = data.transparentedge_sites.all
}
