#################
### EXAMPLE 1 ###
#################
# The verbatim VCL code is set here in a heredoc string: https://developer.hashicorp.com/terraform/language/expressions/strings#heredoc-strings
# if you already have a VCL configuration it's recommended to import it directly from the dashboard.
variable "code" {
  type    = string
  default = <<EOF
sub vcl_recv {
    if (req.http.host == "www.example.com") {
        set req.backend_hint = cX_myorigin.backend();
        set req.http.TCDN-i3-transform = "auto_webp";
        set req.http.TCDN-Command = "brotli_compress";
        call redirect_https;
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

#################
### EXAMPLE 2 ###
#################
# You can also set the heredoc string directly on the resource or reference an external file:
resource "transparentedge_vclconf" "external_file" {
  vclcode = file("${path.module}/config.vcl")
}

#################
### EXAMPLE 3 ###
#################
# If you keep the configuration inside of a tf file you'll be able to reference other variables.
# This is a good approach since the backends are tied to the configuration
resource "transparentedge_backend" "myorig" {
  name   = "origin1"
  origin = "origin.example.com"
  port   = 443
  ssl    = true

  # healthcheck
  hchost       = "www.origin.example.com"
  hcpath       = "/favicon.ico"
  hcstatuscode = 200
}

resource "transparentedge_vclconf" "backend_dependency" {
  vclcode = <<EOF
sub vcl_recv {
    if (req.http.host == "www.example.com") {
        set req.backend_hint = ${resource.transparentedge_backend.myorig.vclname}.backend();
    }
}
EOF
}

# Another option would be to use template files:
# https://developer.hashicorp.com/terraform/language/functions/templatefile
