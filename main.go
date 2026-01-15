package main

import (
	"context"
	"log"

	policyProvider "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/provider/policy_registry"
	tfprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	ctx := context.Background()

	err := providerserver.Serve(
		ctx,
		func() tfprovider.Provider {
			return policyProvider.New()
		},
		providerserver.ServeOpts{Address: "registry.terraform.io/ayuspoudel/sentinel"},
	)

	if err != nil {
		log.Fatal(err)
	}
}
