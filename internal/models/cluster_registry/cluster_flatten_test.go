package clusterModel

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenClusterResponse(t *testing.T) {
	plan := ClusterModel{
		Name:    types.StringValue("sreCluster"),
		Context: types.StringValue("sreCluster"),
	}

	resp := &ClusterResponse{
		Name:      "sreCluster",
		CreatedAt: "2026-01-18T10:00:00Z",
		Labels: map[string]string{
			"context": "sreCluster",
		},
	}

	model := FlattenClusterResponse(resp, plan)

	if model.Name.ValueString() != "sreCluster" {
		t.Fatalf("name mismatch")
	}

	if model.Context.ValueString() != "sreCluster" {
		t.Fatalf("context mismatch")
	}

	if model.RegisteredAt.ValueString() != "2026-01-18T10:00:00Z" {
		t.Fatalf("registered_at mismatch")
	}
}
