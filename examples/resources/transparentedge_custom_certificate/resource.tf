# Both privatekey and publickey must be in PEM format.

resource "transparentedge_custom_certificate" "mysite" {

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

output "mysite" {
  value = transparentedge_custom_certificate.mysite
}
