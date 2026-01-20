package clusterClient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
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
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
}

/*
Author: @ayuspoudel
This struct & its methods are provider facing. Provider will provide
attributes to these methods. It abstracts the core sentinel details
like API endpoints, HTTP methods from provider.
*/
func (c *Client) RegisterWithCredentials(ctx context.Context, req *clusterModel.RegisterWithCredentialsPayload) (*clusterModel.ClusterResponse, error) {
	apiEndpoint := c.baseURL + "/v1/clusters/register-with-credentials"

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	_ = writer.WriteField("name", req.ClusterName)

	if req.Context != "" {
		_ = writer.WriteField("context", req.Context)
	}

	part, err := writer.CreateFormFile("kubeconfig", "kubeconfig")
	if err != nil {
		return nil, err
	}

	_, err = part.Write(req.Kubeconfig)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiEndpoint, &body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
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

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	c.addHeaders(httpReq)

	resp, err := c.http.Do(httpReq)
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
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *Client) DeleteCluster(ctx context.Context, name string) error {
	url := fmt.Sprintf("%s/v1/clusters/%s", c.baseURL, name)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	c.addHeaders(httpReq)

	resp, err := c.http.Do(httpReq)
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
