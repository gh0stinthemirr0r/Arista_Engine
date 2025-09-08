package netvisor

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// NetVisorDB handles the NetVisor API database
type NetVisorDB struct {
	db *sql.DB
}

// APIDefinition represents an API definition from the NetVisor database
type APIDefinition struct {
	ID          string `json:"id"`
	Service     string `json:"service"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Tags        string `json:"tags"`
	Parameters  string `json:"parameters"`
	Example     string `json:"example"`
}

// NewNetVisorDB creates a new NetVisor database connection
func NewNetVisorDB(dbPath string) (*NetVisorDB, error) {
	// Ensure the database file exists
	if !filepath.IsAbs(dbPath) {
		dbPath = filepath.Join(".", dbPath)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open NetVisor database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping NetVisor database: %w", err)
	}

	return &NetVisorDB{db: db}, nil
}

// Close closes the database connection
func (n *NetVisorDB) Close() error {
	if n.db != nil {
		return n.db.Close()
	}
	return nil
}

// GetTables returns all tables in the database
func (n *NetVisorDB) GetTables() ([]string, error) {
	query := "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name"
	rows, err := n.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

// GetTableSchema returns the schema for a specific table
func (n *NetVisorDB) GetTableSchema(tableName string) ([]string, error) {
	query := "PRAGMA table_info(" + tableName + ")"
	rows, err := n.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query table schema: %w", err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull int
		var defaultValue sql.NullString
		var pk int

		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk); err != nil {
			return nil, fmt.Errorf("failed to scan column info: %w", err)
		}
		columns = append(columns, fmt.Sprintf("%s %s", name, dataType))
	}

	return columns, nil
}

// GetAllAPIs returns all API definitions from the database
func (n *NetVisorDB) GetAllAPIs() ([]APIDefinition, error) {
	// First, let's see what tables exist
	tables, err := n.GetTables()
	if err != nil {
		return nil, err
	}

	// Look for common API table names
	var apiTable string
	for _, table := range tables {
		if table == "apis" || table == "api_definitions" || table == "endpoints" || 
		   table == "api_catalog" || table == "netvisor_apis" {
			apiTable = table
			break
		}
	}

	if apiTable == "" {
		// If no obvious API table, return the first table's data
		if len(tables) > 0 {
			apiTable = tables[0]
		} else {
			return nil, fmt.Errorf("no tables found in NetVisor database")
		}
	}

	// Get the schema to understand the column structure
	columns, err := n.GetTableSchema(apiTable)
	if err != nil {
		return nil, err
	}

	// Build a dynamic query based on available columns
	query := fmt.Sprintf("SELECT * FROM %s LIMIT 100", apiTable)
	rows, err := n.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query APIs: %w", err)
	}
	defer rows.Close()

	var apis []APIDefinition
	for rows.Next() {
		// Create a slice to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan API row: %w", err)
		}

		// Convert to APIDefinition (we'll need to map columns dynamically)
		api := APIDefinition{}
		
		// Try to map common column names
		for i, col := range columns {
			if values[i] != nil {
				value := fmt.Sprintf("%v", values[i])
				switch {
				case col == "id" || col == "ID":
					api.ID = value
				case col == "service" || col == "Service":
					api.Service = value
				case col == "method" || col == "Method":
					api.Method = value
				case col == "path" || col == "Path":
					api.Path = value
				case col == "description" || col == "Description":
					api.Description = value
				case col == "category" || col == "Category":
					api.Category = value
				case col == "tags" || col == "Tags":
					api.Tags = value
				case col == "parameters" || col == "Parameters":
					api.Parameters = value
				case col == "example" || col == "Example":
					api.Example = value
				}
			}
		}

		// If we don't have an ID, generate one
		if api.ID == "" {
			api.ID = fmt.Sprintf("netvisor_%d", len(apis)+1)
		}

		apis = append(apis, api)
	}

	return apis, nil
}

// GetAPIsByService returns APIs filtered by service type
func (n *NetVisorDB) GetAPIsByService(service string) ([]APIDefinition, error) {
	apis, err := n.GetAllAPIs()
	if err != nil {
		return nil, err
	}

	var filtered []APIDefinition
	for _, api := range apis {
		if api.Service == service || service == "" {
			filtered = append(filtered, api)
		}
	}

	return filtered, nil
}

// SearchAPIs searches for APIs by keyword
func (n *NetVisorDB) SearchAPIs(keyword string) ([]APIDefinition, error) {
	apis, err := n.GetAllAPIs()
	if err != nil {
		return nil, err
	}

	var results []APIDefinition
	keyword = fmt.Sprintf("%%%s%%", keyword) // Add wildcards for LIKE query

	for _, api := range apis {
		if contains(api.Description, keyword) || 
		   contains(api.Path, keyword) || 
		   contains(api.Category, keyword) ||
		   contains(api.Tags, keyword) {
			results = append(results, api)
		}
	}

	return results, nil
}

// Helper function to check if a string contains a substring (case insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    len(s) > len(substr) && 
		    (s[:len(substr)] == substr || 
		     s[len(s)-len(substr):] == substr || 
		     indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
