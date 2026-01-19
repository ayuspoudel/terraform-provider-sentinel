package clusterSchema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ClusterAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cluster_name":  schema.StringAttribute{Required: true},
		"labels":        schema.MapAttribute{Optional: true, ElementType: types.StringType},
		"source":        schema.StringAttribute{Optional: true},
		"registered_at": schema.StringAttribute{Computed: true},
	}
}

func ClusterBlocks() map[string]schema.Block {
	return map[string]schema.Block{
		"kubeconfig": kubeconfigBlock(),
	}
}

func kubeconfigBlock() schema.Block {
	return schema.SingleNestedBlock{
		Blocks: map[string]schema.Block{
			"file": kubeconfigfileBlock(),
			"env":  kubeconfigEnvBlock(),
			"s3":   kubeconfigS3Block(),
		},
	}
}

func kubeconfigfileBlock() schema.Block {
	return schema.SingleNestedBlock{
		Attributes: map[string]schema.Attribute{
			"path":    schema.StringAttribute{Required: true},
			"context": schema.StringAttribute{Optional: true},
		},
	}
}

func kubeconfigEnvBlock() schema.Block {
	return schema.SingleNestedBlock{
		Attributes: map[string]schema.Attribute{
			"name":     schema.StringAttribute{Required: true},
			"encoding": schema.StringAttribute{Optional: true},
			"context":  schema.StringAttribute{Optional: true},
		},
	}
}

func kubeconfigS3Block() schema.Block {
	return schema.SingleNestedBlock{
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{Required: true},
			"key":    schema.StringAttribute{Optional: true},
			"region": schema.StringAttribute{Optional: true},
		},
	}
}
