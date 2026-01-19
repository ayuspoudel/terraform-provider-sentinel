package clusterSchema

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ClusterDataSourceSchema() schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster_name": schema.StringAttribute{
				Required: true,
			},

			"labels": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},

			"source": schema.StringAttribute{
				Computed: true,
			},

			"registered_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}
