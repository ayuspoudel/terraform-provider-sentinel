package clusterModel

import "github.com/hashicorp/terraform-plugin-framework/types"

type ClusterModel struct {
	Name         types.String    `tfsdk:"cluster_name"`
	Labels       types.Map       `tfsdk:"labels"`
	Source       types.String    `tfsdk:"source"`
	RegisteredAt types.String    `tfsdk:"registered_at"`
	Kubeconfig   KubeconfigModel `tfsdk:"kubeconfig"`
}

type KubeconfigModel struct {
	File *KubeconfigFileModel `tfsdk:"file"`
	Env  *KubeconfigEnvModel  `tfsdk:"env"`
	S3   *KubeconfigS3Model   `tfsdk:"s3"`
}

type KubeconfigFileModel struct {
	Path    types.String `tfsdk:"path"`
	Context types.String `tfsdk:"context"`
}

type KubeconfigEnvModel struct {
	Name     types.String `tfsdk:"name"`
	Encoding types.String `tfsdk:"encoding"`
	Context  types.String `tfsdk:"context"`
}

type KubeconfigS3Model struct {
	Bucket types.String `tfsdk:"bucket"`
	Key    types.String `tfsdk:"key"`
	Region types.String `tfsdk:"region"`
}
