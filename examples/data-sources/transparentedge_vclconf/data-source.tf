data "transparentedge_vclconf" "vclconfig" {}

output "active_vcl_config" {
  value = data.transparentedge_vclconf.vclconfig
}
