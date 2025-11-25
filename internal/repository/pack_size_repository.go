package repository

// PackSizeRepository defines the interface for pack size persistence operations
// This abstraction allows for easy swapping between in-memory and database implementations
type PackSizeRepository interface {
	// GetAll retrieves all pack sizes
	GetAll() ([]int, error)
	
	// SetAll replaces all pack sizes with the provided ones
	SetAll(packSizes []int) error
	
	// Add adds a new pack size
	Add(packSize int) error
	
	// Remove removes a pack size
	Remove(packSize int) error
	
	// Exists checks if a pack size exists
	Exists(packSize int) (bool, error)
}

