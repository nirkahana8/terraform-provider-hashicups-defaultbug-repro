package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"terraform-provider-usfu/internal/provider"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.New("dev"), providerserver.ServeOpts{
		Address: "registry.terraform.io/hashicorp/usfu",
	})
	if err != nil {
		log.Fatal(err)
	}
}
