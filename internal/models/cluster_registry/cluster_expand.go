package clusterModel

import (
	"fmt"
)

/*
Author: @ayuspoudel

ExpandRegisterRequest converts terraform-facing ClusterModel into a payload
usable by the Sentinel Cluster Registry RegisterWithCredentials API.

Terraform manages:
- cluster_name
- context
- kubeconfig bytes

Labels are derived server-side (labels.context).
*/
func ExpandRegisterRequest(m ClusterModel) (*RegisterWithCredentialsPayload, error) {
	if m.Name.IsNull() || m.Name.ValueString() == "" {
		return nil, fmt.Errorf("cluster_name must be provided")
	}

	if m.Kubeconfig.IsNull() || m.Kubeconfig.ValueString() == "" {
		return nil, fmt.Errorf("kubeconfig must be provided")
	}

	out := &RegisterWithCredentialsPayload{
		ClusterName: m.Name.ValueString(),
		Kubeconfig:  []byte(m.Kubeconfig.ValueString()),
		Source:      "terraform",
	}

	if !m.Context.IsNull() && m.Context.ValueString() != "" {
		out.Context = m.Context.ValueString()
	}

	return out, nil
}
