# The verbatim VCL code is set here in a heredoc string: https://developer.hashicorp.com/terraform/language/expressions/strings#heredoc-strings
# if you already have a VCL configuration it's recommended to import it directly from the dashboard.
variable "code" {
  type    = string
  default = <<EOF
sub vcl_recv {
    if (req.http.host == "www.example.com") {
        set req.backend_hint = cX_myorigin.backend();
        set req.http.TCDN-i3-transform = "auto_webp";
        set bereq.http.TCDN-Command = "redirect_https, brotli_compress";
    }
}
EOF
}

resource "transparentedge_vclconf" "prod" {
  vclcode = var.code
}

output "prod_config" {
  value = transparentedge_vclconf.prod
}
