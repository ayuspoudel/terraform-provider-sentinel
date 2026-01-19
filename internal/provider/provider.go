package provider

import (
	"context"

	clusterDatasource "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/datasource/cluster_registry"
	policyDatasource "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/datasource/policy_registry"

	clusterResource "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/resource/cluster_registry"
	policyResource "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/resource/policy_registry"

	clusterConfig "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/config/cluster_registry"
	policyConfig "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/config/policy_registry"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type SentinelProvider struct{}

func New() provider.Provider {
	return &SentinelProvider{}
}

func (p *SentinelProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sentinel"
}

func (p *SentinelProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
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

func (p *SentinelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var policyCfg policyConfig.ProviderConfig
	diags := req.Config.Get(ctx, &policyCfg)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyClient, err := policyCfg.ToClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(" failed to configure policy client", err.Error())
		return
	}

	clusterCfg := clusterConfig.ProviderConfig(policyCfg)
	clusterClient, err := clusterCfg.ToClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("failed to configure cluster client", err.Error())
		return
	}

	resp.DataSourceData = map[string]any{
		"policy":  policyClient,
		"cluster": clusterClient,
	}

	resp.ResourceData = resp.DataSourceData
}

func (p *SentinelProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		policyDatasource.NewPolicyStatusDataSource,
		clusterDatasource.NewClusterDataSource,
	}
}

func (p *SentinelProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		policyResource.NewPolicyResource,
		clusterResource.NewClusterResource,
	}
}
