data "transparentedge_certificates" "all" {}

output "all_certificates" {
  value = data.transparentedge_certificates.all
}
