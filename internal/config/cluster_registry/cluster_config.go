package clusterConfig

import (
	"context"
	"fmt"

	clusterClient "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/client/cluster_registry"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*
Author: @ayuspoudel

ProviderConfig defines configuration accepted by the Sentinel Cluster Registry
provider block.
*/
type ProviderConfig struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Token    types.String `tfsdk:"token"`
}

/*
Author: @ayuspoudel

ToClient validates provider config and constructs a Cluster Registry client.
*/
func (p *ProviderConfig) ToClient(ctx context.Context) (*clusterClient.Client, error) {
	if p.Endpoint.IsNull() || p.Endpoint.ValueString() == "" {
		return nil, fmt.Errorf("sentinel provider: endpoint must be set")
	}

	var token string
	if !p.Token.IsNull() {
		token = p.Token.ValueString()
	}

	return clusterClient.NewClient(p.Endpoint.ValueString(), token), nil
}
