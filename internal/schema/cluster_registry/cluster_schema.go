package clusterSchema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func ClusterAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cluster_name":  schema.StringAttribute{Required: true},
		"context":       schema.StringAttribute{Optional: true},
		"source":        schema.StringAttribute{Computed: true},
		"registered_at": schema.StringAttribute{Computed: true},
		/*
		 Author: @ayuspoudel

		 kubeconfig contains resolved kubeconfig BYTES.
		 Terraform is responsible for resolving file/env/s3.
		 Sentinel treats this as opaque data.
		*/
		"kubeconfig": schema.StringAttribute{Required: true, Sensitive: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}},
	}
}

func ClusterBlocks() map[string]schema.Block {
	return map[string]schema.Block{}
}
