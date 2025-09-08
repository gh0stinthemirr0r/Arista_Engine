package core

import "time"

// EndpointType represents the type of Arista endpoint
type EndpointType string

const (
	EndpointEAPI     EndpointType = "eapi"
	EndpointCV       EndpointType = "cloudvision"
	EndpointEOSREST  EndpointType = "eos_rest"
	EndpointTelemetry EndpointType = "telemetry"
)

// Endpoint represents an Arista device or CloudVision controller
type Endpoint struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Type      EndpointType `json:"type"`
	URL       string       `json:"url"`
	Username  string       `json:"username,omitempty"`
	Password  string       `json:"password,omitempty"` // store encrypted
	Token     string       `json:"token,omitempty"`    // store encrypted
	Created   time.Time    `json:"created"`
	Tags      []string     `json:"tags"`
	TLSVerify bool         `json:"tlsVerify"`
	Status    string       `json:"status,omitempty"` // Connected, Warning, Failed
}

// APIDefinition represents a discovered API endpoint
type APIDefinition struct {
	ID          string   `json:"id"`
	Service     string   `json:"service"`     // eapi, cloudvision, telemetry
	Method      string   `json:"method"`      // GET/POST/PUT/DELETE
	Path        string   `json:"path"`
	Description string   `json:"description"`
	Params      []string `json:"params"`
	Category    string   `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// ExplorerRequest represents an API request from the UI
type ExplorerRequest struct {
	EndpointID string                 `json:"endpointId"`
	Method     string                 `json:"method"` // GET/POST/PUT/DELETE (CV REST) or "runCmds" (eAPI)
	Path       string                 `json:"path"`   // e.g. "/command-api" or "/api/resources/..."
	Body       map[string]any         `json:"body,omitempty"`
	TimeoutMs  int                    `json:"timeoutMs,omitempty"`
}

// ExplorerResponse represents the response from an API call
type ExplorerResponse struct {
	Status     int                    `json:"status"`
	Headers    map[string][]string    `json:"headers"`
	JSON       any                    `json:"json,omitempty"`
	Text       string                 `json:"text,omitempty"`
	ElapsedMs  int64                  `json:"elapsedMs"`
	EndpointID string                 `json:"endpointId"`
	LogID      string                 `json:"logId"`
	Error      string                 `json:"error,omitempty"`
}

// APIQueryRecord represents a logged API query
type APIQueryRecord struct {
	ID         string                 `json:"id"`
	EndpointID string                 `json:"endpointId"`
	Method     string                 `json:"method"`
	Path       string                 `json:"path"`
	Body       map[string]any         `json:"body,omitempty"`
	Status     int                    `json:"status"`
	Response   map[string]any         `json:"response"`
	Timestamp  time.Time              `json:"timestamp"`
	ElapsedMs  int64                  `json:"elapsedMs"`
	Error      string                 `json:"error,omitempty"`
}

// ConnectionTestResult represents the result of testing an endpoint connection
type ConnectionTestResult struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode,omitempty"`
	ElapsedMs  int64  `json:"elapsedMs"`
	Details    any    `json:"details,omitempty"`
}

// DeviceInventory represents a formal device inventory record
type DeviceInventory struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	DeviceType   string    `json:"deviceType" db:"device_type"`
	URL          string    `json:"url" db:"url"`
	Username     string    `json:"username" db:"username"`
	Password     string    `json:"password" db:"password"` // Note: In production, this should be encrypted
	Type         string    `json:"type" db:"type"`         // eapi, cloudvision, eos_rest, telemetry
	Status       string    `json:"status" db:"status"`     // connected, disconnected, testing, failed
	AddedAt      time.Time `json:"addedAt" db:"added_at"`
	LastTested   time.Time `json:"lastTested" db:"last_tested"`
	TestCount    int       `json:"testCount" db:"test_count"`
	SuccessCount int       `json:"successCount" db:"success_count"`
	Notes        string    `json:"notes" db:"notes"`
}

// APICatalog represents the complete enumerated API surface
type APICatalog struct {
	EAPI        map[string]APIDefinition `json:"eapi"`
	CloudVision map[string]APIDefinition `json:"cloudvision"`
	EOSREST     map[string]APIDefinition `json:"eos_rest"`
	Telemetry   map[string]APIDefinition `json:"telemetry"`
	LastUpdated time.Time                `json:"lastUpdated"`
}

// CommandTemplate represents a pre-built command template
type CommandTemplate struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Service     string                 `json:"service"` // eapi, cloudvision
	Method      string                 `json:"method"`
	Path        string                 `json:"path"`
	Body        map[string]any         `json:"body,omitempty"`
	Category    string                 `json:"category"`
	Tags        []string               `json:"tags"`
}

// PolicyRule represents a safety policy rule
type PolicyRule struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Resource    string            `json:"resource"`
	Action      string            `json:"action"`
	Conditions  map[string]string `json:"conditions"`
	Effect      string            `json:"effect"` // allow, deny
	Enabled     bool              `json:"enabled"`
}
