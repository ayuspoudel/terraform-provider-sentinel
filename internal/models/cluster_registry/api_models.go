package clusterModel

/*
Author: @ayuspoudel

These structs define the API payload expected by Sentinel Cluster Registry
for registering clusters with credentials. These are NOT terraform models
and should never use terraform-plugin-framework types.
*/

type RegisterWithCredentialsPayload struct {
	ClusterName string
	Kubeconfig  []byte
	Context     string
	Source      string
}

type ClusterResponse struct {
	Name          string            `json:"name"`
	CredentialRef string            `json:"credential_ref,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	CreatedAt     string            `json:"created_at"`
}
