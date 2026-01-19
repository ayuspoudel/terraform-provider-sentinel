package clusterModel

/*
Author: @ayuspoudel

These structs define the API payload expected by Sentinel Cluster Registry
for registering clusters with credentials. These are NOT terraform models
and should never use terraform-plugin-framework types.
*/

type RegisterWithCredentialsRequest struct {
	ClusterName string                `json:"cluster_name"`
	Source      string                `json:"source,omitempty"`
	Labels      map[string]string     `json:"labels,omitempty"`
	Kubeconfig  KubeconfigCredentials `json:"kubeconfig"`
}

type KubeconfigCredentials struct {
	File *KubeconfigFileCredentials `json:"file,omitempty"`
	Env  *KubeconfigEnvCredentials  `json:"env,omitempty"`
	S3   *KubeconfigS3Credentials   `json:"s3,omitempty"`
}

type KubeconfigFileCredentials struct {
	Path    string `json:"path"`
	Context string `json:"context,omitempty"`
}

type KubeconfigEnvCredentials struct {
	Name     string `json:"name"`
	Encoding string `json:"encoding,omitempty"`
	Context  string `json:"context,omitempty"`
}

type KubeconfigS3Credentials struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key,omitempty"`
	Region string `json:"region,omitempty"`
}

type ClusterResponse struct {
	ClusterName   string            `json:"cluster_name"`
	CredentialRef string            `json:"credential_ref,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	RegisteredAt  string            `json:"registered_at"`
	Source        string            `json:"source"`
}
