package clusterClient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	clusterModel "github.com/ayuspoudel/sentinel-sre/terraform-provider/internal/models/cluster_registry"
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{baseURL: baseURL, token: token, http: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
}

func (c *Client) RegisterWithCredentials(ctx context.Context, req *clusterModel.RegisterWithCredentialsRequest) (*clusterModel.ClusterResponse, error) {
	apiEndpoint := c.baseURL + "/v1/clusters/register-with-credentials"
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	c.addHeaders(httpReq)

	response, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusConflict {
		return nil, fmt.Errorf("cluster already exists")
	}
	if response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("failed to register cluster, status code: %d", response.StatusCode)
	}
	var cluster clusterModel.ClusterResponse
	err = json.NewDecoder(response.Body).Decode(&cluster)
	if err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (c *Client) GetCluster(ctx context.Context, name string) (*clusterModel.ClusterResponse, error) {
	url := fmt.Sprintf("%s/v1/clusters/%s", c.baseURL, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	c.addHeaders(req)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("failed to get cluster: %s", resp.Status)
	}

	var out clusterModel.ClusterResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) DeleteCluster(ctx context.Context, name string) error {
	url := fmt.Sprintf("%s/v1/clusters/%s", c.baseURL, name)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	c.addHeaders(req)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("failed to delete cluster: %s", resp.Status)
	}

	return nil
}
