package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// CloudVisionClient handles communication with Arista CloudVision
type CloudVisionClient struct {
	http *http.Client
}

// NewCloudVisionClient creates a new CloudVision client
func NewCloudVisionClient(tlsVerify bool, timeout time.Duration) *CloudVisionClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !tlsVerify},
	}
	return &CloudVisionClient{http: &http.Client{Transport: tr, Timeout: timeout}}
}

// DoREST performs a REST API call to CloudVision
func (c *CloudVisionClient) DoREST(ctx context.Context, method, url, bearer string, body any) (*http.Response, time.Duration, error) {
	var rdr io.Reader
	if body != nil {
		raw, _ := json.Marshal(body)
		rdr = bytes.NewReader(raw)
	}
	
	req, err := http.NewRequestWithContext(ctx, method, url, rdr)
	if err != nil {
		return nil, 0, err
	}
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	start := time.Now()
	resp, err := c.http.Do(req)
	return resp, time.Since(start), err
}

// TestConnection tests the connection to CloudVision
func (c *CloudVisionClient) TestConnection(ctx context.Context, baseURL, token string) (bool, string, time.Duration, error) {
	// Test with a simple API call
	testURL := baseURL + "/api/resources/inventory/v1/Devices?limit=1"
	
	resp, elapsed, err := c.DoREST(ctx, "GET", testURL, token, nil)
	if err != nil {
		return false, err.Error(), elapsed, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return true, "Connection successful", elapsed, nil
	}

	return false, "Connection failed", elapsed, nil
}

// GetModels attempts to discover available CloudVision models
func (c *CloudVisionClient) GetModels(ctx context.Context, baseURL, token string) ([]string, error) {
	// Known CloudVision models from the documentation
	models := []string{
		"action.v1",
		"alert.v1",
		"bugexposure.v1",
		"changecontrol.v1",
		"configlet.v1",
		"configstatus.v1",
		"connectivitymonitor.v1",
		"dashboard.v1",
		"endpointlocation.v1",
		"event.v1",
		"identityprovider.v1",
		"imagestatus.v1",
		"inventory.v1",
		"lifecycle.v1",
		"redirector.v1",
		"serviceaccount.v1",
		"softwaremanagement.v1",
		"studio.v1",
		"studio_topology.v1",
		"tag.v2",
		"workspace.v1",
	}

	// Test which models are available by attempting to access them
	var availableModels []string
	
	for _, model := range models {
		testURL := baseURL + "/api/resources/" + model
		resp, _, err := c.DoREST(ctx, "GET", testURL, token, nil)
		if err != nil {
			continue
		}
		resp.Body.Close()
		
		if resp.StatusCode == 200 || resp.StatusCode == 404 { // 404 means model exists but no data
			availableModels = append(availableModels, model)
		}
	}

	return availableModels, nil
}

// GetDevices retrieves device information from CloudVision
func (c *CloudVisionClient) GetDevices(ctx context.Context, baseURL, token string) (any, error) {
	url := baseURL + "/api/resources/inventory/v1/Devices"
	resp, _, err := c.DoREST(ctx, "GET", url, token, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetEvents retrieves events from CloudVision
func (c *CloudVisionClient) GetEvents(ctx context.Context, baseURL, token string) (any, error) {
	url := baseURL + "/api/resources/event/v1/Events"
	resp, _, err := c.DoREST(ctx, "GET", url, token, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
