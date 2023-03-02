package main

import (
	"context"
	"flag"
	"log"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name transparentedge --rendered-provider-name TransparentEdge

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/TransparentEdge/transparentedge",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), transparentedge.New, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
