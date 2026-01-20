package kubeconfigSourceModel

import "github.com/hashicorp/terraform-plugin-framework/types"

type KubeconfigSourceModel struct {
	ID         types.String   `tfsdk:"id"`
	Bytes      types.String   `tfsdk:"bytes"`
	Kubeconfig KubeconfigSpec `tfsdk:"kubeconfig"`
}

type KubeconfigSpec struct {
	File []KubeconfigFile `tfsdk:"file"`
	Env  []KubeconfigEnv  `tfsdk:"env"`
	S3   []KubeconfigS3   `tfsdk:"s3"`
}

type KubeconfigFile struct {
	Path    types.String `tfsdk:"path"`
	Context types.String `tfsdk:"context"`
}

type KubeconfigEnv struct {
	Name     types.String `tfsdk:"name"`
	Encoding types.String `tfsdk:"encoding"`
	Context  types.String `tfsdk:"context"`
}

type KubeconfigS3 struct {
	Bucket  types.String `tfsdk:"bucket"`
	Key     types.String `tfsdk:"key"`
	Region  types.String `tfsdk:"region"`
	Context types.String `tfsdk:"context"`
}
