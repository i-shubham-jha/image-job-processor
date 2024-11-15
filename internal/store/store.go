package store

import (
	"encoding/csv"
	"os"
)

// StoreManager manages the store IDs
type StoreManager struct {
	storeIDs map[string]struct{}
}

// NewStoreManager creates a new instance of StoreManager and loads store IDs from a CSV file
func NewStoreManager(filePath string) (*StoreManager, error) {
	sm := &StoreManager{
		storeIDs: make(map[string]struct{}),
	}

	if err := sm.loadStoreIDs(filePath); err != nil {
		return nil, err
	}

	return sm, nil
}

// loadStoreIDs loads store IDs from a CSV file into the map
func (sm *StoreManager) loadStoreIDs(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Check if there are records and skip the header
	if len(records) > 0 {
		for _, record := range records[1:] { // Start from the second record to skip the header
			if len(record) > 2 { // Ensure there are enough columns
				storeID := record[2] // StoreID is in the third column (index 2)
				sm.storeIDs[storeID] = struct{}{}
			}
		}
	}

	return nil
}

// StoreIDExists checks if a store ID exists in the map
func (sm *StoreManager) StoreIDExists(storeID string) bool {
	_, exists := sm.storeIDs[storeID]
	return exists
}
