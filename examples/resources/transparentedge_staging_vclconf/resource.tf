# The usage of the resource 'staging_vclconf' is exactly the same as the production resource 'vclconf'
# for extended documentation refer to 'vclconf' taking care to replace 'vclconf' by 'staging_vclconf'

resource "transparentedge_staging_vclconf" "staging" {
  vclcode = <<EOF
sub vcl_recv {
    # Staging configuration
    if (req.http.host == "www.example.com") {
        set req.backend_hint = cX_stmyorigin.backend();
        set req.http.TCDN-i3-transform = "auto_webp";
        set req.http.TCDN-Command = "redirect_https, brotli_compress";
    }
}
EOF
}

output "staging_config" {
  value = transparentedge_staging_vclconf.staging
}
