package clusterSchema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestClusterSchema_KubeconfigAttribute(t *testing.T) {
	attrs := ClusterAttributes()

	attr, ok := attrs["kubeconfig"]
	if !ok {
		t.Fatalf("expected kubeconfig attribute to exist")
	}

	strAttr, ok := attr.(schema.StringAttribute)
	if !ok {
		t.Fatalf("expected kubeconfig to be StringAttribute")
	}

	if !strAttr.Required {
		t.Fatalf("expected kubeconfig to be required")
	}

	if !strAttr.Sensitive {
		t.Fatalf("expected kubeconfig to be sensitive")
	}
}

func TestClusterSchema_NoBlocks(t *testing.T) {
	blocks := ClusterBlocks()

	if len(blocks) != 0 {
		t.Fatalf("expected no blocks in cluster schema")
	}
}

func TestClusterSchema_ContextAttribute(t *testing.T) {
	attrs := ClusterAttributes()

	if _, ok := attrs["context"]; !ok {
		t.Fatalf("expected context attribute")
	}
}

func TestClusterSchema_NoLabelsAttribute(t *testing.T) {
	attrs := ClusterAttributes()

	if _, ok := attrs["labels"]; ok {
		t.Fatalf("labels must not be configurable")
	}
}
