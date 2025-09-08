package uiapi

import (
	"arista_engine/internal/client"
	"arista_engine/internal/core"
	"arista_engine/internal/store"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ExplorerAPI handles API requests from the UI
type ExplorerAPI struct {
	store      *store.Store
	eapiClient *client.EAPIClient
	cvClient   *client.CloudVisionClient
}

// NewExplorerAPI creates a new ExplorerAPI instance
func NewExplorerAPI(store *store.Store, eapiClient *client.EAPIClient, cvClient *client.CloudVisionClient) *ExplorerAPI {
	return &ExplorerAPI{
		store:      store,
		eapiClient: eapiClient,
		cvClient:   cvClient,
	}
}

// RunAPIRequest executes an API request
func (e *ExplorerAPI) RunAPIRequest(ctx context.Context, request core.ExplorerRequest) (core.ExplorerResponse, error) {
	// Get endpoint configuration
	endpoint, err := e.store.GetEndpoint(request.EndpointID)
	if err != nil {
		return core.ExplorerResponse{}, fmt.Errorf("failed to get endpoint: %w", err)
	}

	// Set timeout
	timeout := 30 * time.Second
	if request.TimeoutMs > 0 {
		timeout = time.Duration(request.TimeoutMs) * time.Millisecond
	}
	
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Execute request based on endpoint type
	var response core.ExplorerResponse
	switch endpoint.Type {
	case core.EndpointEAPI:
		response, err = e.handleEAPIRequest(ctx, endpoint, request)
	case core.EndpointCV:
		response, err = e.handleCloudVisionRequest(ctx, endpoint, request)
	default:
		return core.ExplorerResponse{}, fmt.Errorf("unsupported endpoint type: %s", endpoint.Type)
	}

	if err != nil {
		response.Error = err.Error()
	}

	// Log the request
	record := core.APIQueryRecord{
		ID:         fmt.Sprintf("req_%d", time.Now().UnixNano()),
		EndpointID: request.EndpointID,
		Method:     request.Method,
		Path:       request.Path,
		Body:       request.Body,
		Status:     response.Status,
		Response:   map[string]any{"json": response.JSON, "text": response.Text},
		Timestamp:  time.Now(),
		ElapsedMs:  response.ElapsedMs,
		Error:      response.Error,
	}

	if err := e.store.SaveQueryRecord(record); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to save query record: %v\n", err)
	}

	response.LogID = record.ID
	return response, nil
}

// handleEAPIRequest handles EOS eAPI requests
func (e *ExplorerAPI) handleEAPIRequest(ctx context.Context, endpoint core.Endpoint, request core.ExplorerRequest) (core.ExplorerResponse, error) {
	// Convert request body to RunCmdsParams
	params, err := e.convertToRunCmdsParams(request.Body)
	if err != nil {
		return core.ExplorerResponse{}, fmt.Errorf("failed to convert request body: %w", err)
	}

	// Execute the request
	rpc, httpResp, elapsed, err := e.eapiClient.RunCmds(ctx, endpoint.URL, endpoint.Username, endpoint.Password, params)
	if err != nil && rpc == nil {
		return core.ExplorerResponse{}, err
	}

	// Build response
	response := core.ExplorerResponse{
		Status:     httpResp.StatusCode,
		Headers:    httpResp.Header,
		ElapsedMs:  elapsed.Milliseconds(),
		EndpointID: endpoint.ID,
	}

	// Handle JSON-RPC response
	if rpc != nil {
		if rpc.Error != nil {
			response.Error = rpc.Error.Message
			response.JSON = map[string]any{
				"error": rpc.Error,
			}
		} else {
			response.JSON = map[string]any{
				"result": rpc.Result,
			}
		}
	}

	return response, err
}

// handleCloudVisionRequest handles CloudVision REST requests
func (e *ExplorerAPI) handleCloudVisionRequest(ctx context.Context, endpoint core.Endpoint, request core.ExplorerRequest) (core.ExplorerResponse, error) {
	// Build full URL
	fullURL := endpoint.URL + request.Path

	// Execute the request
	resp, elapsed, err := e.cvClient.DoREST(ctx, request.Method, fullURL, endpoint.Token, request.Body)
	if err != nil {
		return core.ExplorerResponse{}, err
	}
	defer resp.Body.Close()

	// Build response
	response := core.ExplorerResponse{
		Status:     resp.StatusCode,
		Headers:    resp.Header,
		ElapsedMs:  elapsed.Milliseconds(),
		EndpointID: endpoint.ID,
	}

	// Parse response body
	var parsed any
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		// If JSON parsing fails, treat as text
		response.Text = fmt.Sprintf("Failed to parse JSON: %v", err)
	} else {
		response.JSON = parsed
	}

	return response, nil
}

// convertToRunCmdsParams converts a request body to RunCmdsParams
func (e *ExplorerAPI) convertToRunCmdsParams(body map[string]any) (client.RunCmdsParams, error) {
	params := client.RunCmdsParams{
		Version:       1,
		Format:        "json",
		AutoComplete:  true,
		ExpandAliases: true,
	}

	// Extract commands
	if cmds, ok := body["cmds"].([]interface{}); ok {
		for _, cmd := range cmds {
			if cmdStr, ok := cmd.(string); ok {
				params.Cmds = append(params.Cmds, cmdStr)
			}
		}
	} else {
		return params, errors.New("cmds field is required and must be an array of strings")
	}

	// Extract optional fields
	if version, ok := body["version"].(float64); ok {
		params.Version = int(version)
	}
	if format, ok := body["format"].(string); ok {
		params.Format = format
	}
	if autoComplete, ok := body["autoComplete"].(bool); ok {
		params.AutoComplete = autoComplete
	}
	if expandAliases, ok := body["expandAliases"].(bool); ok {
		params.ExpandAliases = expandAliases
	}

	return params, nil
}
