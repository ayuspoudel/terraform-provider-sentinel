package clusterResource

import (
	"context"

	clusterClient "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/client/cluster_registry"
	clusterModel "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/models/cluster_registry"
	clusterSchema "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/schema/cluster_registry"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ClusterResource struct {
	client *clusterClient.Client
}

func NewClusterResource() resource.Resource {
	return &ClusterResource{}
}

/*
Author: @ayuspoudel
Provider has the data for client stored in ProviderData inside req struct.
This needs to be passed to resource as req.ProviderData.(*client.ClusterClient)
so r.client is initialized with the client needed.
At this point our ClusterResource will have all the methods defined for creation
deletion, update and get of clusters from sentinel which was the details abstracted
by the client.
*/
func (r *ClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data := req.ProviderData.(map[string]any)
	r.client = data["cluster"].(*clusterClient.Client)
}

/*
Author: @ayuspoudel
All this does is set metadata for resource name.
This is similar to aws_instance or aws_iam_profile and so on.
*/
func (r *ClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "sentinel_cluster"
}

/*
Author: @ayuspoudel
At this function we are passing the definitions we have created in
schema/cluster_schema.go for sentinel cluster resource. This will let
resource know it has this following schema valid, and users should
pass the inputs. It also has optional or defined fields already fulfiled.
*/
func (r *ClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: clusterSchema.ClusterAttributes(),
		Blocks:     clusterSchema.ClusterBlocks(),
	}
}

/*
Author: @ayuspoudel
This function creates the actual resource cluster we have defined.
It uses request given by terraform from user input, which we store in our
variable of type ClusterModel which is defined based on terraform types. Then
we see if there are any errors else respond to terraform with the error.
Then we use client.RegisterCluster(...) to register the cluster in sentinel
and if apply is successful we set the state file with the contents of our plan.
*/
func (r *ClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan clusterModel.ClusterModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq, err := clusterModel.ExpandRegisterRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError("invalid sentinel_cluster configuration", err.Error())
		return
	}
	cluster, err := r.client.RegisterWithCredentials(ctx, apiReq)

	if err != nil {
		resp.Diagnostics.AddError("failed to create sentinel_cluster", err.Error())
		return
	}
	state := clusterModel.FlattenClusterResponse(cluster, plan)
	state.Kubeconfig = plan.Kubeconfig

	resp.State.Set(ctx, state)
}

/*
Author: @ayuspoudel
This is used by terraform to get the current details of the cluster via the GET
endpoint provided by sentinel and abstracted by our client.
*/
func (r *ClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state clusterModel.ClusterModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cluster, err := r.client.GetCluster(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to get current cluster", err.Error())
		return
	}

	if cluster == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	newState := clusterModel.FlattenClusterResponse(cluster, state)
	newState.Kubeconfig = state.Kubeconfig
	resp.State.Set(ctx, newState)

}

/*
Author: @ayuspoudel
Sentinel Cluster Registry API is idempotent so internally both update and create
function contents remain the same.
*/
func (r *ClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan clusterModel.ClusterModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq, err := clusterModel.ExpandRegisterRequest(plan)
	if err != nil {
		resp.Diagnostics.AddError("invalid sentinel_cluster configuration", err.Error())
		return
	}

	cluster, err := r.client.RegisterWithCredentials(ctx, apiReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to update sentinel_cluster", err.Error())
		return
	}

	state := clusterModel.FlattenClusterResponse(cluster, plan)
	state.Kubeconfig = plan.Kubeconfig
	resp.State.Set(ctx, state)
}

/*
Author: @ayuspoudel
This function is used by terraform to delete any resources. For this we simply
get the current state from req resource.DeleteRequest's State.Get() method
and use client.DeleteCluster(...) method to delete the cluster from sentinel.
*/
func (r *ClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state clusterModel.ClusterModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCluster(ctx, state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("failed to delete sentinel_cluster", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
func (r *ClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	clusterName := req.ID
	if clusterName == "" {
		resp.Diagnostics.AddError("invalid import ID", "import ID cannot be empty")
		return
	}
	cluster, err := r.client.GetCluster(ctx, clusterName)
	if err != nil {
		resp.Diagnostics.AddError("failed to get cluster during import", err.Error())
		return
	}
	if cluster == nil {
		resp.Diagnostics.AddError("cluster not found", "no cluster found with the given name")
		return
	}
	plan := clusterModel.ClusterModel{Name: types.StringValue(clusterName), Kubeconfig: types.StringNull()}

	state := clusterModel.FlattenClusterResponse(cluster, plan)
	resp.State.Set(ctx, state)
}
