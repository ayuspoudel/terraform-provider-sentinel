package clusterModel

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExpandRegisterRequest_File(t *testing.T) {
	model := ClusterModel{
		Name:   types.StringValue("sreCluster"),
		Source: types.StringValue("terraform"),
		Labels: types.MapValueMust(types.StringType, map[string]attr.Value{
			"env": types.StringValue("prod"),
		}),
		Kubeconfig: KubeconfigModel{
			File: &KubeconfigFileModel{
				Path:    types.StringValue("/tmp/kubeconfig"),
				Context: types.StringValue("default"),
			},
		},
	}

	req, err := ExpandRegisterRequest(model)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.ClusterName != "sreCluster" {
		t.Fatalf("cluster name mismatch")
	}

	if req.Kubeconfig.File == nil {
		t.Fatalf("expected file kubeconfig")
	}
}

func TestExpandRegisterRequest_MultipleSources(t *testing.T) {
	model := ClusterModel{
		Name: types.StringValue("sreCluster"),
		Kubeconfig: KubeconfigModel{
			File: &KubeconfigFileModel{Path: types.StringValue("/tmp/k")},
			Env:  &KubeconfigEnvModel{Name: types.StringValue("KUBECONFIG")},
		},
	}

	_, err := ExpandRegisterRequest(model)
	if err == nil {
		t.Fatalf("expected error for multiple kubeconfig sources")
	}
}
