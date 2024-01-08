# TransparentEdge Terraform Provider

A terraform provider for the CDN of [Transparent Edge](https://www.transparentedge.eu/).  
This provider is intended to be used by the CDN users with the role _Company Admin_, although some services do not require this role.

## Links

- [Overview](https://registry.terraform.io/providers/TransparentEdge/transparentedge/latest)
- [Official documentation](https://registry.terraform.io/providers/TransparentEdge/transparentedge/latest/docs)
- [API documentation](https://api.transparentcdn.com/docs/)

## Example usage

```terraform
terraform {
  required_providers {
    transparentedge = {
      source = "TransparentEdge/transparentedge"
      version = ">=0.3.3"
    }
  }
}

provider "transparentedge" {
  # Provider configuration overrides environment variables
  # it's recommended to use environment variables for company_id, client_id and client_secret
  company_id    = 300
  client_id     = "XXX"
  client_secret = "XXX"
  insecure      = false                            # this is the default value
  api_url       = "https://api.transparentcdn.com" # this is the default value
}
```

It's recommended to use environment variables:

```shell
export TCDN_COMPANY_ID=0
export TCDN_CLIENT_ID="xxx"
export TCDN_CLIENT_SECRET="xxx"
```

You can find all the required variables in the [dashboard](https://dashboard.transparentcdn.com/): `"Profile" -> "Account options" -> "Manage keys".`

Make sure that you're using the correct Company ID if you own multiple companies.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This builds and installs the provider binary in the `$GOPATH/bin` or `$GOBIN` directory if you have it set.

To generate or update documentation, run `go generate`.

To override the provider locally, it's usually enough to create the file `~/.terraformrc` (UNIX) or `%APPDATA%\terraform.rc` (Windows). Replace `${GOPATH}` or `${GOBIN}` with the correct one for your machine:

```
# ~/.terraformrc
provider_installation {
  dev_overrides {
      # Use $GOPATH/bin or $GOBIN
      "TransparentEdge/transparentedge" = "${GOBIN}"
      "registry.terraform.io/hashicorp/transparentedge" = "${GOBIN}"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```
