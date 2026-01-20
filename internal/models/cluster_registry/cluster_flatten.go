package clusterModel

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FlattenClusterResponse(resp *ClusterResponse, plan ClusterModel) ClusterModel {
	var context types.String

	if ctx, ok := resp.Labels["context"]; ok {
		context = types.StringValue(ctx)
	} else {
		context = types.StringNull()
	}

	return ClusterModel{
		Name:         plan.Name,
		Context:      context,
		Source:       types.StringValue("terraform"),
		RegisteredAt: types.StringValue(resp.CreatedAt),
		Kubeconfig:   types.StringNull(),
	}
}

func stringMapToAttr(m map[string]string) map[string]attr.Value {
	out := make(map[string]attr.Value, len(m))
	for k, v := range m {
		out[k] = types.StringValue(v)
	}
	return out
}
