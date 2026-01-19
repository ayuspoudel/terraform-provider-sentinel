package clusterModel

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FlattenClusterResponse(r *ClusterResponse) ClusterModel {
	labels := map[string]attr.Value{}
	for k, v := range r.Labels {
		labels[k] = types.StringValue(v)
	}

	return ClusterModel{
		Name:         types.StringValue(r.ClusterName),
		Labels:       types.MapValueMust(types.StringType, labels),
		Source:       types.StringValue(r.Source),
		RegisteredAt: types.StringValue(r.RegisteredAt),
	}
}
