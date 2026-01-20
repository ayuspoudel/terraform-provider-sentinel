package clusterDatasource

import (
	"context"

	clusterClient "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/client/cluster_registry"
	clusterModel "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/models/cluster_registry"
	clusterSchema "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/schema/cluster_registry"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type ClusterDataSource struct {
	client *clusterClient.Client
}

func NewClusterDataSource() datasource.DataSource {
	return &ClusterDataSource{}
}

func (d *ClusterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "sentinel_cluster"
}

func (d *ClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = clusterSchema.ClusterDataSourceSchema()
}

func (d *ClusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data := req.ProviderData.(map[string]any)
	d.client = data["cluster"].(*clusterClient.Client)
}

func (d *ClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state clusterModel.ClusterDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cluster, err := d.client.GetCluster(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get cluster", err.Error())
		return
	}

	if cluster == nil {
		resp.Diagnostics.AddError("cluster not found", "sentinel returned no cluster")
		return
	}

	newState := clusterModel.FlattenClusterResponseDataSource(cluster)
	resp.State.Set(ctx, newState)
}
