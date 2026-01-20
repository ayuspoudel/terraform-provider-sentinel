package kubeconfigSourceSchema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func KubeconfigSourceSchema() schema.Schema {
	return schema.Schema{
		Description: "A Terraform-managed source of kubeconfig bytes used by Sentinel clusters.",

		Attributes: map[string]schema.Attribute{
			"id":    schema.StringAttribute{Computed: true},
			"bytes": schema.StringAttribute{Computed: true, Sensitive: true},
		},

		Blocks: map[string]schema.Block{
			"kubeconfig": kubeconfigBlock(),
		},
	}
}

func kubeconfigBlock() schema.Block {
	return schema.SingleNestedBlock{
		Blocks: map[string]schema.Block{
			"file": kubeconfigFileBlock(),
			"env":  kubeconfigEnvBlock(),
			"s3":   kubeconfigS3Block(),
		},
	}
}

func kubeconfigFileBlock() schema.Block {
	return schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"path":    schema.StringAttribute{Required: true},
				"context": schema.StringAttribute{Optional: true},
			},
		},
	}
}

func kubeconfigEnvBlock() schema.Block {
	return schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"name":     schema.StringAttribute{Required: true},
				"encoding": schema.StringAttribute{Optional: true},
				"context":  schema.StringAttribute{Optional: true},
			},
		},
	}
}

func kubeconfigS3Block() schema.Block {
	return schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"bucket":  schema.StringAttribute{Required: true},
				"key":     schema.StringAttribute{Required: true},
				"region":  schema.StringAttribute{Optional: true},
				"context": schema.StringAttribute{Optional: true},
			},
		},
	}
}
