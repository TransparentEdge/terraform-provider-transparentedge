package main

import (
	"context"
	"flag"
	"log"

	transparentedge "github.com/TransparentEdge/terraform-provider-transparentedge/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name transparentedge --rendered-provider-name TransparentEdge

var (
	// variables are set by goreleaser
	version = "dev"
	commit  = "none"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/TransparentEdge/transparentedge",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), transparentedge.New(version, commit), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
