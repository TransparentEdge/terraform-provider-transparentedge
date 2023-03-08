data "transparentedge_staging_backend" "mybackend" {
  name = "mybackendname"
}

output "vclname" {
  # use 'vclname' to associate a backend in VCL Code
  # for example: set req.backend_hint = ${data.transparentedge_staging_backend.mybackend.vclname}.backend();
  value = data.transparentedge_staging_backend.mybackend.vclname
}
