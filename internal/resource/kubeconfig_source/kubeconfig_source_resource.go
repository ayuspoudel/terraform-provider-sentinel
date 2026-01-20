package kubeconfigSourceResource

import (
	"context"

	kubeconfigSourceModel "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/models/kubeconfig_source"
	kubeconfigSourceSchema "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/schema/kubeconfig_source"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*
Author: @ayuspoudel

This resource represents a Terraform-managed kubeconfig source.
It resolves kubeconfig bytes locally and never talks to Sentinel.
*/
type KubeconfigSourceResource struct{}

func New() resource.Resource { return &KubeconfigSourceResource{} }

/*
Metadata sets the Terraform resource name.
*/
func (r *KubeconfigSourceResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "sentinel_kubeconfig_source"
}

/*
Schema defines the kubeconfig source structure.
*/
func (r *KubeconfigSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = kubeconfigSourceSchema.KubeconfigSourceSchema()
}

/*
Create resolves kubeconfig bytes and stores them in state.
*/
func (r *KubeconfigSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan kubeconfigSourceModel.KubeconfigSourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, bytes, _, err := kubeconfigSourceModel.Expand(plan)
	if err != nil {
		resp.Diagnostics.AddError("invalid kubeconfig source", err.Error())
		return
	}

	plan.ID = types.StringValue(id)
	plan.Bytes = types.StringValue(string(bytes))

	resp.State.Set(ctx, plan)
}

/*
Read re-resolves the kubeconfig source to detect drift.
*/
func (r *KubeconfigSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state kubeconfigSourceModel.KubeconfigSourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, bytes, _, err := kubeconfigSourceModel.Expand(state)
	if err != nil {
		resp.Diagnostics.AddError("failed to read kubeconfig source", err.Error())
		return
	}

	state.ID = types.StringValue(id)
	state.Bytes = types.StringValue(string(bytes))

	resp.State.Set(ctx, state)
}

/*
Update behaves the same as Create since the resource is idempotent.
*/
func (r *KubeconfigSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan kubeconfigSourceModel.KubeconfigSourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, bytes, _, err := kubeconfigSourceModel.Expand(plan)
	if err != nil {
		resp.Diagnostics.AddError("invalid kubeconfig source", err.Error())
		return
	}

	plan.ID = types.StringValue(id)
	plan.Bytes = types.StringValue(string(bytes))

	resp.State.Set(ctx, plan)
}

/*
Delete removes the resource from Terraform state only.
*/
func (r *KubeconfigSourceResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}
