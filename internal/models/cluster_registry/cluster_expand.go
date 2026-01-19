package clusterModel

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

/*
Author: @ayuspoudel

ExpandRegisterRequest converts terraform-facing ClusterModel into an
API-facing RegisterWithCredentialsRequest. This is the boundary where
we move from terraform types to plain Go structs that will be serialized
and sent to Sentinel Cluster Registry.

At this layer we:
- validate exactly one kubeconfig source is provided
- extract values from terraform types
- build the request expected by the HTTP client
*/

func ExpandRegisterRequest(m ClusterModel) (*RegisterWithCredentialsRequest, error) {
	req := &RegisterWithCredentialsRequest{ClusterName: m.Name.ValueString(), Source: m.Source.ValueString()}
	if !m.Labels.IsNull() {
		labels := map[string]string{}
		for k, v := range m.Labels.Elements() {
			labels[k] = v.(types.String).ValueString()
		}
		req.Labels = labels
	}
	// Kubeconfig source extraction

	var sources int
	var kc KubeconfigCredentials
	if m.Kubeconfig.File != nil {
		sources++
		if !m.Kubeconfig.File.Path.IsNull() {
			kc.File = &KubeconfigFileCredentials{
				Path:    m.Kubeconfig.File.Path.ValueString(),
				Context: m.Kubeconfig.File.Context.ValueString(),
			}
		}
	}
	if m.Kubeconfig.Env != nil {
		sources++
		if !m.Kubeconfig.Env.Name.IsNull() {
			kc.Env = &KubeconfigEnvCredentials{
				Name:     m.Kubeconfig.Env.Name.ValueString(),
				Encoding: m.Kubeconfig.Env.Encoding.ValueString(),
				Context:  m.Kubeconfig.Env.Context.ValueString(),
			}
		}
	}
	if m.Kubeconfig.S3 != nil {
		sources++
		if !m.Kubeconfig.S3.Bucket.IsNull() {
			kc.S3 = &KubeconfigS3Credentials{
				Bucket: m.Kubeconfig.S3.Bucket.ValueString(),
				Key:    m.Kubeconfig.S3.Key.ValueString(),
				Region: m.Kubeconfig.S3.Region.ValueString(),
			}
		}
	}

	if sources == 0 {
		return nil, fmt.Errorf("at least one kubeconfig source must be provided")
	}
	if sources > 1 {
		return nil, fmt.Errorf("only one kubeconfig source can be provided")
	}

	req.Kubeconfig = kc
	return req, nil

}
