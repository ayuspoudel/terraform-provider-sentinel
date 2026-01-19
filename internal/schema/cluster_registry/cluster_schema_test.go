package clusterSchema

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestKubeconfigBlock_SubBlocksExist(t *testing.T) {
	block := kubeconfigBlock()

	nested, ok := block.(schema.SingleNestedBlock)
	if !ok {
		t.Fatalf("expected kubeconfig block to be SingleNestedBlock")
	}

	if _, ok := nested.Blocks["file"]; !ok {
		t.Fatalf("expected kubeconfig.file block")
	}

	if _, ok := nested.Blocks["env"]; !ok {
		t.Fatalf("expected kubeconfig.env block")
	}

	if _, ok := nested.Blocks["s3"]; !ok {
		t.Fatalf("expected kubeconfig.s3 block")
	}
}
