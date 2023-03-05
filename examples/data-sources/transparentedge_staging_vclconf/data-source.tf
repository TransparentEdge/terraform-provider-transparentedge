data "transparentedge_staging_vclconf" "vclconfig" {}

output "staging_vcl_config" {
  value = data.transparentedge_staging_vclconf.vclconfig
}
