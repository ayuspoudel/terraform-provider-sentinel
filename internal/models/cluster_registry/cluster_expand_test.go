package clusterModel

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExpandRegisterRequest(t *testing.T) {
	m := ClusterModel{
		Name:       types.StringValue("sreCluster"),
		Context:    types.StringValue("sreCluster"),
		Kubeconfig: types.StringValue("dummy-kubeconfig-bytes"),
	}

	req, err := ExpandRegisterRequest(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.ClusterName != "sreCluster" {
		t.Fatalf("cluster name mismatch")
	}

	if req.Context != "sreCluster" {
		t.Fatalf("context mismatch")
	}

	if string(req.Kubeconfig) != "dummy-kubeconfig-bytes" {
		t.Fatalf("kubeconfig mismatch")
	}
}

func TestExpandRegisterRequest_MissingKubeconfig(t *testing.T) {
	m := ClusterModel{
		Name: types.StringValue("sreCluster"),
	}

	_, err := ExpandRegisterRequest(m)
	if err == nil {
		t.Fatalf("expected error for missing kubeconfig")
	}
}
