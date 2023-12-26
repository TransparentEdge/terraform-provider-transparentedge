terraform {
  required_providers {
    transparentedge = {
      source = "TransparentEdge/transparentedge"
      # Available since version 0.5.0
      version = ">=0.5.0"
    }
  }
}

# Example for 'AWS (Route53)'
resource "transparentedge_certreq_dns_credential" "aws_01" {
  alias = "aws-creds-01"

  parameters = {
    AWS_SECRET_ACCESS_KEY = "ABC"
    AWS_ACCESS_KEY_ID     = "DEF"
  }
}

# Example for 'Cloudflare'
resource "transparentedge_certreq_dns_credential" "cf_01" {
  alias = "cf-creds-01"

  parameters = {
    CF_Account_ID = "0123"
    CF_Token      = "ABC"
  }
}

# Internal credentials for Transparent Edge
# The following credentials do not require any secrets:

# CNAME Verification credential
# This credential must be attached to a DNS Certificate Request that wants to perform
# the validation using a CNAME. Check the 'transparentedge_certreq_dns_cname_verification' data source
# for more information.
resource "transparentedge_certreq_dns_credential" "cname_verif" {
  alias = "cname-verification"

  parameters = {
    CNAME_VERIFICATION = ""
  }
}

# Transparent Edge DNS Provider credential
# This credential is only valid if your domain DNS has been delegated
# to Transparent Edge.
resource "transparentedge_certreq_dns_credential" "tedge" {
  alias = "tedge-cred"

  parameters = {
    TCDN_DNS_PROVIDER = ""
  }
}
