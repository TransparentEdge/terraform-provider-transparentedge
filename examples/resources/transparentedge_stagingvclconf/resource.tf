# The usage of the resource 'stagingvclconf' is exactly the same as the production resource 'vclconf'
# for extended documentation refer to 'vclconf' take care to replace 'vclconf' by 'stagingvclconf'

resource "transparentedge_stagingvclconf" "staging" {
  vclcode = <<EOF
sub vcl_recv {
    # Staging configuration
    if (req.http.host == "www.example.com") {
        set req.backend_hint = cX_myorigin.backend();
        set req.http.TCDN-i3-transform = "auto_webp";
        set bereq.http.TCDN-Command = "redirect_https, brotli_compress";
    }
}
EOF
}

output "staging_config" {
  value = transparentedge_stagingvclconf.staging
}
