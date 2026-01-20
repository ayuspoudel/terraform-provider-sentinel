package clusterModel

import "github.com/hashicorp/terraform-plugin-framework/types"

type ClusterDataSourceModel struct {
	Name         types.String `tfsdk:"cluster_name"`
	Labels       types.Map    `tfsdk:"labels"`
	Source       types.String `tfsdk:"source"`
	RegisteredAt types.String `tfsdk:"registered_at"`
}

func FlattenClusterResponseDataSource(r *ClusterResponse) ClusterDataSourceModel {
	var labels types.Map

	if len(r.Labels) > 0 {
		labels = types.MapValueMust(types.StringType, stringMapToAttr(r.Labels))
	} else {
		labels = types.MapNull(types.StringType)
	}

	return ClusterDataSourceModel{
		Name:         types.StringValue(r.Name),
		Labels:       labels,
		RegisteredAt: types.StringValue(r.CreatedAt),
	}
}
