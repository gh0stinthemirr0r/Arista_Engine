package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

// EAPIClient handles communication with Arista EOS eAPI
type EAPIClient struct {
	http *http.Client
}

// NewEAPIClient creates a new EAPI client
func NewEAPIClient(tlsVerify bool, timeout time.Duration) *EAPIClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !tlsVerify},
	}
	return &EAPIClient{http: &http.Client{Transport: tr, Timeout: timeout}}
}

// RunCmdsParams represents the parameters for runCmds
type RunCmdsParams struct {
	Version       int      `json:"version"`
	Cmds          []string `json:"cmds"`
	Format        string   `json:"format,omitempty"`        // json|text
	AutoComplete  bool     `json:"autoComplete,omitempty"`
	ExpandAliases bool     `json:"expandAliases,omitempty"`
}

// JSONRPCRequest represents a JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  RunCmdsParams `json:"params"`
	ID      string        `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC response
type JSONRPCResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	Result  []any            `json:"result,omitempty"`
	Error   *JSONRPCError    `json:"error,omitempty"`
	ID      string           `json:"id"`
}

// JSONRPCError represents a JSON-RPC error
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// RunCmds executes commands on an EOS device via eAPI
func (c *EAPIClient) RunCmds(ctx context.Context, baseURL, user, pass string, params RunCmdsParams) (*JSONRPCResponse, *http.Response, time.Duration, error) {
	rpc := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "runCmds",
		Params:  params,
		ID:      "1",
	}
	raw, _ := json.Marshal(rpc)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/command-api", bytes.NewReader(raw))
	if err != nil {
		return nil, nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if user != "" {
		req.SetBasicAuth(user, pass)
	}

	start := time.Now()
	resp, err := c.http.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		return nil, resp, elapsed, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	var out JSONRPCResponse
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, resp, elapsed, err
	}
	if out.Error != nil {
		return &out, resp, elapsed, errors.New(out.Error.Message)
	}
	return &out, resp, elapsed, nil
}

// TestConnection tests the connection to an EOS device
func (c *EAPIClient) TestConnection(ctx context.Context, baseURL, user, pass string) (bool, string, time.Duration, error) {
	params := RunCmdsParams{
		Version:      1,
		Cmds:         []string{"show version"},
		Format:       "json",
		AutoComplete: true,
	}

	_, resp, elapsed, err := c.RunCmds(ctx, baseURL, user, pass, params)
	if err != nil {
		return false, err.Error(), elapsed, err
	}

	if resp.StatusCode == 200 {
		return true, "Connection successful", elapsed, nil
	}

	return false, "Connection failed", elapsed, errors.New("non-200 status code")
}

// EnumerateCommands attempts to discover available CLI commands
func (c *EAPIClient) EnumerateCommands(ctx context.Context, baseURL, user, pass string) ([]string, error) {
	// Start with basic command discovery
	discoveryCommands := []string{
		"show ?",
		"show running-config ?",
		"show interfaces ?",
		"show version",
		"show system",
	}

	var allCommands []string

	for _, cmd := range discoveryCommands {
		params := RunCmdsParams{
			Version:      1,
			Cmds:         []string{cmd},
			Format:       "text", // Use text for help output
			AutoComplete: true,
		}

		_, _, _, err := c.RunCmds(ctx, baseURL, user, pass, params)
		if err != nil {
			continue // Skip failed commands
		}

		allCommands = append(allCommands, cmd)
	}

	return allCommands, nil
}
