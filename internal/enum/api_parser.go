package enum

import (
	"arista_engine/internal/core"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

// APIParser handles parsing and organizing the complete Arista API surface
type APIParser struct {
	catalog *core.APICatalog
}

// NewAPIParser creates a new API parser
func NewAPIParser() *APIParser {
	return &APIParser{
		catalog: &core.APICatalog{
			EAPI:        make(map[string]core.APIDefinition),
			CloudVision: make(map[string]core.APIDefinition),
			EOSREST:     make(map[string]core.APIDefinition),
			Telemetry:   make(map[string]core.APIDefinition),
			LastUpdated: time.Now(),
		},
	}
}

// ParseEnumeratedAPI parses the complete enumerated API from the markdown file
func (p *APIParser) ParseEnumeratedAPI(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open API file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read API file: %w", err)
	}

	return p.parseContent(string(content))
}

// parseContent parses the API content and extracts endpoints
func (p *APIParser) parseContent(content string) error {
	lines := strings.Split(content, "\n")
	
	var currentCategory string
	
	for i, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines and headers
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "Arista Networks") {
			continue
		}
		
		// Detect category headers (e.g., "AaaTacacs", "AccessLists")
		if isCategoryHeader(line) {
			currentCategory = line
			continue
		}
		
		// Parse API endpoints
		if strings.Contains(line, " /") {
			endpoint := p.parseEndpoint(line, currentCategory)
			if endpoint != nil {
				// Determine service type based on path patterns
				service := p.determineService(endpoint.Path)
				endpoint.Service = service
				endpoint.Category = currentCategory
				
				// Add to appropriate catalog
				key := fmt.Sprintf("%s_%s_%s", service, currentCategory, endpoint.Method)
				switch service {
				case "eapi":
					p.catalog.EAPI[key] = *endpoint
				case "cloudvision":
					p.catalog.CloudVision[key] = *endpoint
				case "eos_rest":
					p.catalog.EOSREST[key] = *endpoint
				case "telemetry":
					p.catalog.Telemetry[key] = *endpoint
				}
			}
		}
		
		// Progress indicator for large files
		if i%1000 == 0 && i > 0 {
			fmt.Printf("Parsed %d lines...\n", i)
		}
	}
	
	return nil
}

// isCategoryHeader determines if a line is a category header
func isCategoryHeader(line string) bool {
	// Category headers are typically single words or short phrases
	// and don't contain HTTP methods or paths
	if strings.Contains(line, " /") || strings.Contains(line, "get ") || 
	   strings.Contains(line, "post ") || strings.Contains(line, "put ") || 
	   strings.Contains(line, "delete ") {
		return false
	}
	
	// Must be a reasonable length and not contain special characters
	if len(line) > 50 || strings.Contains(line, " ") && len(line) > 30 {
		return false
	}
	
	return true
}

// parseEndpoint parses a single API endpoint line
func (p *APIParser) parseEndpoint(line, category string) *core.APIDefinition {
	// Extract HTTP method and path
	re := regexp.MustCompile(`(get|post|put|delete)\s+(/[^\s]+)`)
	matches := re.FindStringSubmatch(line)
	
	if len(matches) < 3 {
		return nil
	}
	
	method := strings.ToUpper(matches[1])
	path := matches[2]
	
	// Generate description based on path and category
	description := p.generateDescription(category, path, method)
	
	// Extract parameters from path
	params := p.extractParameters(path)
	
	// Generate tags
	tags := p.generateTags(category, path, method)
	
	return &core.APIDefinition{
		ID:          fmt.Sprintf("%s_%s_%s", category, method, strings.ReplaceAll(path, "/", "_")),
		Service:     "", // Will be set by caller
		Method:      method,
		Path:        path,
		Description: description,
		Params:      params,
		Category:    category,
		Tags:        tags,
	}
}

// determineService determines the service type based on path patterns
func (p *APIParser) determineService(path string) string {
	// CloudVision patterns
	if strings.Contains(path, "/api/") || strings.Contains(path, "/resources/") {
		return "cloudvision"
	}
	
	// Telemetry patterns
	if strings.Contains(path, "/telemetry/") || strings.Contains(path, "/streaming/") {
		return "telemetry"
	}
	
	// eAPI patterns (JSON-RPC over /command-api)
	if strings.Contains(path, "/command-api") {
		return "eapi"
	}
	
	// EOS REST API patterns (all the /vRest endpoints)
	if strings.HasPrefix(path, "/") && !strings.Contains(path, "/api/") {
		return "eos_rest"
	}
	
	// Default to eos_rest for most other patterns
	return "eos_rest"
}

// generateDescription creates a human-readable description
func (p *APIParser) generateDescription(category, path, method string) string {
	// Convert category to readable format
	categoryReadable := strings.ReplaceAll(category, "-", " ")
	categoryReadable = strings.ReplaceAll(categoryReadable, "_", " ")
	
	// Convert path to readable format
	pathReadable := strings.ReplaceAll(path, "/", " ")
	pathReadable = strings.ReplaceAll(pathReadable, "-", " ")
	pathReadable = strings.ReplaceAll(pathReadable, "_", " ")
	pathReadable = strings.TrimSpace(pathReadable)
	
	return fmt.Sprintf("%s %s for %s", strings.Title(method), pathReadable, categoryReadable)
}

// extractParameters extracts parameter names from the path
func (p *APIParser) extractParameters(path string) []string {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(path, -1)
	
	var params []string
	for _, match := range matches {
		if len(match) > 1 {
			params = append(params, match[1])
		}
	}
	
	return params
}

// generateTags creates relevant tags for the endpoint
func (p *APIParser) generateTags(category, path, method string) []string {
	var tags []string
	
	// Add method tag
	tags = append(tags, method)
	
	// Add category tag
	tags = append(tags, strings.ToLower(category))
	
	// Add path-based tags
	if strings.Contains(path, "stats") {
		tags = append(tags, "statistics")
	}
	if strings.Contains(path, "config") {
		tags = append(tags, "configuration")
	}
	if strings.Contains(path, "status") {
		tags = append(tags, "status")
	}
	if strings.Contains(path, "clear") {
		tags = append(tags, "maintenance")
	}
	
	return tags
}

// GetCatalog returns the parsed API catalog
func (p *APIParser) GetCatalog() *core.APICatalog {
	return p.catalog
}

// SaveCatalog saves the catalog to a JSON file
func (p *APIParser) SaveCatalog(filePath string) error {
	data, err := json.MarshalIndent(p.catalog, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal catalog: %w", err)
	}
	
	return os.WriteFile(filePath, data, 0644)
}

// LoadCatalog loads a catalog from a JSON file
func (p *APIParser) LoadCatalog(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read catalog file: %w", err)
	}
	
	return json.Unmarshal(data, &p.catalog)
}

// GetEndpointsByService returns endpoints for a specific service
func (p *APIParser) GetEndpointsByService(service string) map[string]core.APIDefinition {
	switch service {
	case "eapi":
		return p.catalog.EAPI
	case "cloudvision":
		return p.catalog.CloudVision
	case "eos_rest":
		return p.catalog.EOSREST
	case "telemetry":
		return p.catalog.Telemetry
	default:
		return make(map[string]core.APIDefinition)
	}
}

// GetEndpointsByCategory returns endpoints for a specific category
func (p *APIParser) GetEndpointsByCategory(category string) []core.APIDefinition {
	var endpoints []core.APIDefinition
	
	// Search through all services
	for _, endpoint := range p.catalog.EAPI {
		if endpoint.Category == category {
			endpoints = append(endpoints, endpoint)
		}
	}
	for _, endpoint := range p.catalog.CloudVision {
		if endpoint.Category == category {
			endpoints = append(endpoints, endpoint)
		}
	}
	for _, endpoint := range p.catalog.EOSREST {
		if endpoint.Category == category {
			endpoints = append(endpoints, endpoint)
		}
	}
	for _, endpoint := range p.catalog.Telemetry {
		if endpoint.Category == category {
			endpoints = append(endpoints, endpoint)
		}
	}
	
	return endpoints
}

// SearchEndpoints searches for endpoints by query
func (p *APIParser) SearchEndpoints(query string) []core.APIDefinition {
	var results []core.APIDefinition
	query = strings.ToLower(query)
	
	// Search through all services
	for _, endpoint := range p.catalog.EAPI {
		if p.matchesQuery(endpoint, query) {
			results = append(results, endpoint)
		}
	}
	for _, endpoint := range p.catalog.CloudVision {
		if p.matchesQuery(endpoint, query) {
			results = append(results, endpoint)
		}
	}
	for _, endpoint := range p.catalog.EOSREST {
		if p.matchesQuery(endpoint, query) {
			results = append(results, endpoint)
		}
	}
	for _, endpoint := range p.catalog.Telemetry {
		if p.matchesQuery(endpoint, query) {
			results = append(results, endpoint)
		}
	}
	
	return results
}

// matchesQuery checks if an endpoint matches the search query
func (p *APIParser) matchesQuery(endpoint core.APIDefinition, query string) bool {
	// Search in description, path, category, and tags
	searchText := strings.ToLower(fmt.Sprintf("%s %s %s %s %s", 
		endpoint.Description, endpoint.Path, endpoint.Category, 
		endpoint.Method, strings.Join(endpoint.Tags, " ")))
	
	return strings.Contains(searchText, query)
}
