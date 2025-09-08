package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"arista_engine/internal/client"
	"arista_engine/internal/core"
	"arista_engine/internal/enum"
	"arista_engine/internal/netvisor"
	"arista_engine/internal/store"
	"arista_engine/internal/uiapi"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// App struct
type App struct {
	ctx        context.Context
	logger     *zap.Logger
	store      *store.Store
	eapiClient *client.EAPIClient
	cvClient   *client.CloudVisionClient
	apiParser  *enum.APIParser
	uiAPI      *uiapi.ExplorerAPI
	netvisorDB *netvisor.NetVisorDB
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Create logging directory
	logDir := "Logging"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create logging directory: %v", err))
	}

	// Initialize logger
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"stdout",
		filepath.Join(logDir, fmt.Sprintf("arista_engine_%s.log", time.Now().Format("20060102_150405"))),
	}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = "stacktrace"

	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}

	// Initialize store with proper path
	dbPath := "data.db"
	logger.Info("Initializing database", zap.String("path", dbPath))

	store, err := store.NewStore(dbPath)
	if err != nil {
		logger.Fatal("Failed to initialize store", zap.Error(err))
	}
	logger.Info("Database initialized successfully")

	// Initialize API clients
	eapiClient := client.NewEAPIClient(true, 30*time.Second)
	cvClient := client.NewCloudVisionClient(true, 30*time.Second)

	// Initialize API parser
	apiParser := enum.NewAPIParser()

	// Initialize UI API
	uiAPI := uiapi.NewExplorerAPI(store, eapiClient, cvClient)

	// Initialize NetVisor database
	netvisorDB, err := netvisor.NewNetVisorDB("netvisor_api_v711.db")
	if err != nil {
		logger.Error("Failed to initialize NetVisor database", zap.Error(err))
		// Continue without NetVisor database - it's optional
		netvisorDB = nil
	} else {
		logger.Info("NetVisor database initialized successfully")
	}

	return &App{
		logger:     logger,
		store:      store,
		eapiClient: eapiClient,
		cvClient:   cvClient,
		apiParser:  apiParser,
		uiAPI:      uiAPI,
		netvisorDB: netvisorDB,
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Create required directories
	dirs := []string{"Logging", "Reports", "Exports"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			a.logger.Error("Failed to create directory",
				zap.String("dir", dir),
				zap.Error(err),
			)
			runtime.LogError(ctx, fmt.Sprintf("Failed to create directory %s: %v", dir, err))
		} else {
			a.logger.Info("Created directory", zap.String("dir", dir))
		}
	}

	// Load API catalog if it exists, otherwise parse the enumerated API
	catalogPath := "api_catalog.json"
	if _, err := os.Stat(catalogPath); os.IsNotExist(err) {
		a.logger.Info("API catalog not found, parsing enumerated API...")
		if err := a.apiParser.ParseEnumeratedAPI("Enumerated_API.md"); err != nil {
			a.logger.Error("Failed to parse enumerated API", zap.Error(err))
			runtime.LogError(ctx, fmt.Sprintf("Failed to parse enumerated API: %v", err))
		} else {
			// Save the parsed catalog
			if err := a.apiParser.SaveCatalog(catalogPath); err != nil {
				a.logger.Error("Failed to save API catalog", zap.Error(err))
			} else {
				a.logger.Info("API catalog saved successfully")
			}
		}
	} else {
		// Load existing catalog
		if err := a.apiParser.LoadCatalog(catalogPath); err != nil {
			a.logger.Error("Failed to load API catalog", zap.Error(err))
		} else {
			a.logger.Info("API catalog loaded successfully")
		}
	}

	a.logger.Info("Arista Engine started successfully")
	runtime.LogInfo(ctx, "Arista Engine started successfully")
}

// GetEndpoints returns all configured endpoints
func (a *App) GetEndpoints() ([]core.Endpoint, error) {
	return a.store.GetEndpoints()
}

// AddEndpoint adds a new endpoint
func (a *App) AddEndpoint(endpoint core.Endpoint) error {
	endpoint.ID = fmt.Sprintf("ep_%d", time.Now().UnixNano())
	endpoint.Created = time.Now()

	if err := a.store.SaveEndpoint(endpoint); err != nil {
		a.logger.Error("Failed to save endpoint", zap.Error(err))
		return err
	}

	// Also add to device inventory
	device := core.DeviceInventory{
		ID:           endpoint.ID,
		Name:         endpoint.Name,
		DeviceType:   string(endpoint.Type), // This will be the API type (eapi, cloudvision, etc.)
		URL:          endpoint.URL,
		Username:     endpoint.Username,
		Password:     endpoint.Password,
		Type:         string(endpoint.Type),
		Status:       "disconnected",
		AddedAt:      time.Now(),
		LastTested:   time.Time{},
		TestCount:    0,
		SuccessCount: 0,
		Notes:        fmt.Sprintf("Added via Endpoint Manager - %s", endpoint.Type),
	}

	if err := a.store.AddDeviceToInventory(device); err != nil {
		a.logger.Error("Failed to add device to inventory", zap.Error(err))
		// Don't fail the endpoint creation if inventory fails
	}

	a.logger.Info("Endpoint added", zap.String("id", endpoint.ID), zap.String("name", endpoint.Name))
	return nil
}

// UpdateEndpoint updates an existing endpoint
func (a *App) UpdateEndpoint(endpoint core.Endpoint) error {
	if err := a.store.SaveEndpoint(endpoint); err != nil {
		a.logger.Error("Failed to update endpoint", zap.Error(err))
		return err
	}

	a.logger.Info("Endpoint updated", zap.String("id", endpoint.ID), zap.String("name", endpoint.Name))
	return nil
}

// DeleteEndpoint deletes an endpoint
func (a *App) DeleteEndpoint(endpointID string) error {
	if err := a.store.DeleteEndpoint(endpointID); err != nil {
		a.logger.Error("Failed to delete endpoint", zap.Error(err))
		return err
	}

	a.logger.Info("Endpoint deleted", zap.String("id", endpointID))
	return nil
}

// TestConnection tests connection to an endpoint
func (a *App) TestConnection(endpointID string) (core.ConnectionTestResult, error) {
	endpoint, err := a.store.GetEndpoint(endpointID)
	if err != nil {
		return core.ConnectionTestResult{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var success bool
	var message string
	var elapsed time.Duration

	switch endpoint.Type {
	case core.EndpointEAPI:
		success, message, elapsed, err = a.eapiClient.TestConnection(ctx, endpoint.URL, endpoint.Username, endpoint.Password)
	case core.EndpointCV:
		success, message, elapsed, err = a.cvClient.TestConnection(ctx, endpoint.URL, endpoint.Token)
	default:
		return core.ConnectionTestResult{}, fmt.Errorf("unsupported endpoint type: %s", endpoint.Type)
	}

	result := core.ConnectionTestResult{
		Success:   success,
		Message:   message,
		ElapsedMs: elapsed.Milliseconds(),
	}

	if err != nil {
		result.Message = err.Error()
	}

	// Update endpoint status
	endpoint.Status = "Connected"
	if !success {
		endpoint.Status = "Failed"
	}
	a.store.SaveEndpoint(endpoint)

	return result, nil
}

// GetAPICatalog returns the complete API catalog
func (a *App) GetAPICatalog() (*core.APICatalog, error) {
	return a.apiParser.GetCatalog(), nil
}

// ServeAPICatalog serves the API catalog as JSON for frontend consumption
func (a *App) ServeAPICatalog() (map[string]interface{}, error) {
	catalog := a.apiParser.GetCatalog()
	if catalog == nil {
		return nil, fmt.Errorf("API catalog not loaded")
	}

	// Convert to map for JSON serialization
	result := map[string]interface{}{
		"eapi":        catalog.EAPI,
		"cloudvision": catalog.CloudVision,
		"eos_rest":    catalog.EOSREST,
		"telemetry":   catalog.Telemetry,
		"lastUpdated": catalog.LastUpdated,
	}

	return result, nil
}

// GetEndpointsByService returns endpoints for a specific service
func (a *App) GetEndpointsByService(service string) (map[string]core.APIDefinition, error) {
	return a.apiParser.GetEndpointsByService(service), nil
}

// SearchEndpoints searches for API endpoints
func (a *App) SearchEndpoints(query string) ([]core.APIDefinition, error) {
	return a.apiParser.SearchEndpoints(query), nil
}

// RunAPIRequest executes an API request
func (a *App) RunAPIRequest(request core.ExplorerRequest) (core.ExplorerResponse, error) {
	return a.uiAPI.RunAPIRequest(context.Background(), request)
}

// GetQueryLog returns the query log
func (a *App) GetQueryLog() ([]core.APIQueryRecord, error) {
	return a.store.GetQueryLog()
}

// ExportResults exports results in various formats
func (a *App) ExportResults(format string, data []core.APIQueryRecord) (string, error) {
	// This will be implemented in the export system
	return "", fmt.Errorf("export not yet implemented")
}

// GetDeviceInventory retrieves the device inventory
func (a *App) GetDeviceInventory() ([]core.DeviceInventory, error) {
	return a.store.GetDeviceInventory()
}

// ServeDeviceInventory serves the device inventory as JSON for frontend consumption
func (a *App) ServeDeviceInventory() ([]core.DeviceInventory, error) {
	return a.store.GetDeviceInventory()
}

// GetNetVisorAPIs returns all APIs from the NetVisor database
func (a *App) GetNetVisorAPIs() ([]netvisor.APIDefinition, error) {
	if a.netvisorDB == nil {
		return nil, fmt.Errorf("NetVisor database not available")
	}
	return a.netvisorDB.GetAllAPIs()
}

// GetNetVisorAPIsByService returns APIs from NetVisor database filtered by service
func (a *App) GetNetVisorAPIsByService(service string) ([]netvisor.APIDefinition, error) {
	if a.netvisorDB == nil {
		return nil, fmt.Errorf("NetVisor database not available")
	}
	return a.netvisorDB.GetAPIsByService(service)
}

// SearchNetVisorAPIs searches APIs in the NetVisor database
func (a *App) SearchNetVisorAPIs(keyword string) ([]netvisor.APIDefinition, error) {
	if a.netvisorDB == nil {
		return nil, fmt.Errorf("NetVisor database not available")
	}
	return a.netvisorDB.SearchAPIs(keyword)
}

// GetNetVisorTables returns all tables in the NetVisor database
func (a *App) GetNetVisorTables() ([]string, error) {
	if a.netvisorDB == nil {
		return nil, fmt.Errorf("NetVisor database not available")
	}
	return a.netvisorDB.GetTables()
}

// TestAPIConnection tests a connection to an API endpoint by type
func (a *App) TestAPIConnection(apiType, url, username, password, token string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch apiType {
	case "eapi":
		if username == "" || password == "" {
			return map[string]interface{}{
				"success": false,
				"message": "Authentication required",
				"details": "Username and password are required for eAPI",
			}, nil
		}
		success, message, _, err := a.eapiClient.TestConnection(ctx, url, username, password)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"message": "Connection failed",
				"details": err.Error(),
			}, nil
		}
		return map[string]interface{}{
			"success": success,
			"message": message,
			"details": fmt.Sprintf("Connected to EOS eAPI at %s", url),
		}, nil
	case "cloudvision":
		if token == "" {
			return map[string]interface{}{
				"success": false,
				"message": "Authentication required",
				"details": "API token is required for CloudVision",
			}, nil
		}
		success, message, _, err := a.cvClient.TestConnection(ctx, url, token)
		if err != nil {
			return map[string]interface{}{
				"success": false,
				"message": "Connection failed",
				"details": err.Error(),
			}, nil
		}
		return map[string]interface{}{
			"success": success,
			"message": message,
			"details": fmt.Sprintf("Connected to CloudVision Portal at %s", url),
		}, nil
	case "eos_rest":
		if username == "" || password == "" {
			return map[string]interface{}{
				"success": false,
				"message": "Authentication required",
				"details": "Username and password are required for EOS REST API",
			}, nil
		}
		// For now, just validate URL format for EOS REST
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			return map[string]interface{}{
				"success": false,
				"message": "Invalid URL format",
				"details": "URL must start with http:// or https://",
			}, nil
		}
		// TODO: Implement actual EOS REST API test
		return map[string]interface{}{
			"success": false,
			"message": "EOS REST API test not implemented",
			"details": "EOS REST API connection testing will be implemented",
		}, nil
	default:
		return map[string]interface{}{
			"success": false,
			"message": "Unknown API type",
			"details": fmt.Sprintf("API type '%s' is not supported", apiType),
		}, nil
	}
}
