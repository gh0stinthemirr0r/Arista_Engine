package store

import (
	"arista_engine/internal/core"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
)

// Store handles data persistence
type Store struct {
	db *bolt.DB
}

// NewStore creates a new store instance
func NewStore(dbPath string) (*Store, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Check if database exists
	dbExists := true
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dbExists = false
	}

	// Open database (BoltDB will create it if it doesn't exist)
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{
		Timeout: 5 * time.Second,
		NoGrowSync: false,
		ReadOnly: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &Store{db: db}

	// Initialize buckets
	if err := store.initBuckets(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize buckets: %w", err)
	}

	if !dbExists {
		// Log that we created a new database
		fmt.Printf("Created new database at: %s\n", dbPath)
	}

	return store, nil
}

// initBuckets initializes the database buckets
func (s *Store) initBuckets() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		buckets := []string{"endpoints", "query_log", "api_catalog", "device_inventory"}
		for _, bucket := range buckets {
			if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
				return fmt.Errorf("failed to create bucket %s: %w", bucket, err)
			}
		}
		return nil
	})
}

// SaveEndpoint saves an endpoint to the database
func (s *Store) SaveEndpoint(endpoint core.Endpoint) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("endpoints"))
		if bucket == nil {
			return fmt.Errorf("endpoints bucket not found")
		}

		data, err := json.Marshal(endpoint)
		if err != nil {
			return fmt.Errorf("failed to marshal endpoint: %w", err)
		}

		return bucket.Put([]byte(endpoint.ID), data)
	})
}

// GetEndpoint retrieves an endpoint by ID
func (s *Store) GetEndpoint(id string) (core.Endpoint, error) {
	var endpoint core.Endpoint

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("endpoints"))
		if bucket == nil {
			return fmt.Errorf("endpoints bucket not found")
		}

		data := bucket.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("endpoint not found")
		}

		return json.Unmarshal(data, &endpoint)
	})

	return endpoint, err
}

// GetEndpoints retrieves all endpoints
func (s *Store) GetEndpoints() ([]core.Endpoint, error) {
	var endpoints []core.Endpoint

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("endpoints"))
		if bucket == nil {
			return fmt.Errorf("endpoints bucket not found")
		}

		return bucket.ForEach(func(k, v []byte) error {
			var endpoint core.Endpoint
			if err := json.Unmarshal(v, &endpoint); err != nil {
				return err
			}
			endpoints = append(endpoints, endpoint)
			return nil
		})
	})

	return endpoints, err
}

// DeleteEndpoint deletes an endpoint by ID
func (s *Store) DeleteEndpoint(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("endpoints"))
		if bucket == nil {
			return fmt.Errorf("endpoints bucket not found")
		}

		return bucket.Delete([]byte(id))
	})
}

// SaveQueryRecord saves a query record to the log
func (s *Store) SaveQueryRecord(record core.APIQueryRecord) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("query_log"))
		if bucket == nil {
			return fmt.Errorf("query_log bucket not found")
		}

		data, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal query record: %w", err)
		}

		// Use timestamp as key for chronological ordering
		key := fmt.Sprintf("%d_%s", record.Timestamp.UnixNano(), record.ID)
		return bucket.Put([]byte(key), data)
	})
}

// GetQueryLog retrieves the query log
func (s *Store) GetQueryLog() ([]core.APIQueryRecord, error) {
	var records []core.APIQueryRecord

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("query_log"))
		if bucket == nil {
			return fmt.Errorf("query_log bucket not found")
		}

		return bucket.ForEach(func(k, v []byte) error {
			var record core.APIQueryRecord
			if err := json.Unmarshal(v, &record); err != nil {
				return err
			}
			records = append(records, record)
			return nil
		})
	})

	return records, err
}

// GetQueryLogByEndpoint retrieves query log for a specific endpoint
func (s *Store) GetQueryLogByEndpoint(endpointID string) ([]core.APIQueryRecord, error) {
	var records []core.APIQueryRecord

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("query_log"))
		if bucket == nil {
			return fmt.Errorf("query_log bucket not found")
		}

		return bucket.ForEach(func(k, v []byte) error {
			var record core.APIQueryRecord
			if err := json.Unmarshal(v, &record); err != nil {
				return err
			}
			if record.EndpointID == endpointID {
				records = append(records, record)
			}
			return nil
		})
	})

	return records, err
}

// SaveAPICatalog saves the API catalog
func (s *Store) SaveAPICatalog(catalog *core.APICatalog) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("api_catalog"))
		if bucket == nil {
			return fmt.Errorf("api_catalog bucket not found")
		}

		data, err := json.Marshal(catalog)
		if err != nil {
			return fmt.Errorf("failed to marshal API catalog: %w", err)
		}

		return bucket.Put([]byte("catalog"), data)
	})
}

// GetAPICatalog retrieves the API catalog
func (s *Store) GetAPICatalog() (*core.APICatalog, error) {
	var catalog core.APICatalog

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("api_catalog"))
		if bucket == nil {
			return fmt.Errorf("api_catalog bucket not found")
		}

		data := bucket.Get([]byte("catalog"))
		if data == nil {
			return fmt.Errorf("API catalog not found")
		}

		return json.Unmarshal(data, &catalog)
	})

	return &catalog, err
}

// Device Inventory Methods

// GetDeviceInventory retrieves all device inventory records
func (s *Store) GetDeviceInventory() ([]core.DeviceInventory, error) {
	var devices []core.DeviceInventory

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("device_inventory"))
		if bucket == nil {
			return fmt.Errorf("device_inventory bucket not found")
		}

		return bucket.ForEach(func(k, v []byte) error {
			var device core.DeviceInventory
			if err := json.Unmarshal(v, &device); err != nil {
				return err
			}
			devices = append(devices, device)
			return nil
		})
	})

	return devices, err
}

// AddDeviceToInventory adds a device to the inventory
func (s *Store) AddDeviceToInventory(device core.DeviceInventory) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("device_inventory"))
		if bucket == nil {
			return fmt.Errorf("device_inventory bucket not found")
		}

		data, err := json.Marshal(device)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(device.ID), data)
	})
}

// UpdateDeviceInventory updates a device in the inventory
func (s *Store) UpdateDeviceInventory(device core.DeviceInventory) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("device_inventory"))
		if bucket == nil {
			return fmt.Errorf("device_inventory bucket not found")
		}

		data, err := json.Marshal(device)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(device.ID), data)
	})
}

// DeleteDeviceFromInventory removes a device from the inventory
func (s *Store) DeleteDeviceFromInventory(id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("device_inventory"))
		if bucket == nil {
			return fmt.Errorf("device_inventory bucket not found")
		}

		return bucket.Delete([]byte(id))
	})
}

// GetDeviceInventoryByID retrieves a specific device from inventory
func (s *Store) GetDeviceInventoryByID(id string) (core.DeviceInventory, error) {
	var device core.DeviceInventory

	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("device_inventory"))
		if bucket == nil {
			return fmt.Errorf("device_inventory bucket not found")
		}

		data := bucket.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("device not found")
		}

		return json.Unmarshal(data, &device)
	})

	return device, err
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}
