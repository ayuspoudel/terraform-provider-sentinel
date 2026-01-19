package clusterModel

import "testing"

func TestFlattenClusterResponse(t *testing.T) {
	resp := &ClusterResponse{
		ClusterName:  "sreCluster",
		Source:       "terraform",
		RegisteredAt: "2026-01-18T10:00:00Z",
		Labels: map[string]string{
			"env": "prod",
		},
	}

	model := FlattenClusterResponse(resp)

	if model.Name.ValueString() != "sreCluster" {
		t.Fatalf("name mismatch")
	}

	if model.Source.ValueString() != "terraform" {
		t.Fatalf("source mismatch")
	}
}
