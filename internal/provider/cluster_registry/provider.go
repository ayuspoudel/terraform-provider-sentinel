package clusterProvider

import (
	"context"

	clusterConfig "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/config/cluster_registry"
	clusterDatasource "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/datasource/cluster_registry"
	clusterResource "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/resource/cluster_registry"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type ClusterProvider struct{}

func New() provider.Provider {
	return &ClusterProvider{}
}

/*
Author: @ayuspoudel
This function defines provider-level metadata such as provider name.
*/
func (p *ClusterProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sentinel"
}

/*
Author: @ayuspoudel
This function defines the provider-level configuration schema.
*/
func (p *ClusterProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required: true,
			},
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

/*
Author: @ayuspoudel
This function initializes the provider client and stores it inside ProviderData.
*/
func (p *ClusterProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config clusterConfig.ProviderConfig

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := config.ToClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to configure sentinel provider", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

/*
Author: @ayuspoudel
This function registers all data sources supported by the provider.
*/
func (p *ClusterProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		clusterDatasource.NewClusterDataSource,
	}
}

/*
Author: @ayuspoudel
This function registers all resources supported by the provider.
*/
func (p *ClusterProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		clusterResource.NewClusterResource,
	}
}
