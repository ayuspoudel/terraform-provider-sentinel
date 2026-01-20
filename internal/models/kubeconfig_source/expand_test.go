package kubeconfigSourceModel

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestExpand_EnvPlain(t *testing.T) {
	const envName = "TEST_KUBECONFIG_PLAIN"
	defer os.Unsetenv(envName)

	os.Setenv(envName, "plain-kubeconfig")

	m := KubeconfigSourceModel{
		Kubeconfig: KubeconfigSpec{
			Env: []KubeconfigEnv{
				{
					Name:     types.StringValue(envName),
					Encoding: types.StringNull(),
				},
			},
		},
	}

	id, bytes, _, err := Expand(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(id) == 0 {
		t.Fatalf("expected id to be set")
	}

	if string(bytes) != "plain-kubeconfig" {
		t.Fatalf("unexpected bytes: %q", string(bytes))
	}
}

func TestExpand_EnvBase64(t *testing.T) {
	const envName = "TEST_KUBECONFIG_B64"
	defer os.Unsetenv(envName)

	raw := "base64-kubeconfig"
	encoded := base64.StdEncoding.EncodeToString([]byte(raw))
	os.Setenv(envName, encoded)

	m := KubeconfigSourceModel{
		Kubeconfig: KubeconfigSpec{
			Env: []KubeconfigEnv{
				{
					Name:     types.StringValue(envName),
					Encoding: types.StringValue("base64"),
				},
			},
		},
	}

	_, bytes, _, err := Expand(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(bytes) != raw {
		t.Fatalf("decoded bytes mismatch: %q", string(bytes))
	}
}

func TestExpand_EnvInvalidEncoding(t *testing.T) {
	const envName = "TEST_KUBECONFIG_BAD_ENCODING"
	defer os.Unsetenv(envName)

	os.Setenv(envName, "data")

	m := KubeconfigSourceModel{
		Kubeconfig: KubeconfigSpec{
			Env: []KubeconfigEnv{
				{
					Name:     types.StringValue(envName),
					Encoding: types.StringValue("gzip"),
				},
			},
		},
	}

	_, _, _, err := Expand(m)
	if err == nil {
		t.Fatalf("expected error for invalid encoding")
	}
}

func TestExpand_EnvEmptyValue(t *testing.T) {
	const envName = "TEST_KUBECONFIG_EMPTY"
	defer os.Unsetenv(envName)

	os.Setenv(envName, "")

	m := KubeconfigSourceModel{
		Kubeconfig: KubeconfigSpec{
			Env: []KubeconfigEnv{
				{
					Name: types.StringValue(envName),
				},
			},
		},
	}

	_, _, _, err := Expand(m)
	if err == nil {
		t.Fatalf("expected error for empty env value")
	}
}
