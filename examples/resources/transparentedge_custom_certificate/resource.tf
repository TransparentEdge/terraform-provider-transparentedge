# Both privatekey and publickey must be in PEM format.

# You can specify an external file (relative to this *.tf file)
resource "transparentedge_custom_certificate" "mysite1" {
  publickey  = file("${path.module}/mysite1.crt")
  privatekey = file("${path.module}/mysite1.key")
}

# Or specify the certificate directly with a heredoc string
resource "transparentedge_custom_certificate" "mysite2" {
  # Full chain recommended, it may include multiple 'BEGIN' and 'END' entries
  publickey = <<EOF
-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----

-----BEGIN CERTIFICATE-----
-----END CERTIFICATE-----
EOF

  # Private Key cannot be password protected
  privatekey = <<EOF
-----BEGIN PRIVATE KEY-----
-----END PRIVATE KEY-----
EOF
}
