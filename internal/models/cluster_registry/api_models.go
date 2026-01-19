package clusterModel

type RegisterWithCredentialsRequest struct {
	ClusterName string            `json:"cluster_name"`
	Kubeconfig  string            `json:"kubeconfig"`
	Labels      map[string]string `json:"labels,omitempty"`
	Source      string            `json:"source,omitempty"`
}

type ClusterResponse struct {
	ClusterName   string            `json:"cluster_name"`
	CredentialRef string            `json:"credential_ref,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	RegisteredAt  string            `json:"registered_at"`
	Source        string            `json:"source"`
}
