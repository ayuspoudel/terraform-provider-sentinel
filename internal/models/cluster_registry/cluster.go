package clusterModel

import "github.com/hashicorp/terraform-plugin-framework/types"

type ClusterModel struct {
	Name         types.String `tfsdk:"cluster_name"`
	Context      types.String `tfsdk:"context"`
	Source       types.String `tfsdk:"source"`
	RegisteredAt types.String `tfsdk:"registered_at"`

	Kubeconfig types.String `tfsdk:"kubeconfig"`
}
