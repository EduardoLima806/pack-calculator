package memory

import (
	"sync"

	"pack-calculator/internal/repository"
)

// MemoryPackSizeRepository is an in-memory implementation of PackSizeRepository
// It uses a mutex to ensure thread-safe operations
type MemoryPackSizeRepository struct {
	packSizes []int
	mu        sync.RWMutex
}

// NewMemoryPackSizeRepository creates a new in-memory pack size repository
// with the provided initial pack sizes
func NewMemoryPackSizeRepository(initialPackSizes []int) repository.PackSizeRepository {
	// Create a copy of the initial pack sizes to avoid external modifications
	packSizes := make([]int, len(initialPackSizes))
	copy(packSizes, initialPackSizes)
	
	return &MemoryPackSizeRepository{
		packSizes: packSizes,
	}
}

// GetAll retrieves all pack sizes
func (r *MemoryPackSizeRepository) GetAll() ([]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Return a copy to prevent external modifications
	result := make([]int, len(r.packSizes))
	copy(result, r.packSizes)
	return result, nil
}

// SetAll replaces all pack sizes with the provided ones
func (r *MemoryPackSizeRepository) SetAll(packSizes []int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Create a copy to avoid external modifications
	r.packSizes = make([]int, len(packSizes))
	copy(r.packSizes, packSizes)
	return nil
}

// Add adds a new pack size if it doesn't already exist
func (r *MemoryPackSizeRepository) Add(packSize int) error {
	if packSize <= 0 {
		return ErrInvalidPackSize
	}
	
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Check if pack size already exists
	for _, size := range r.packSizes {
		if size == packSize {
			return nil // Already exists, no error
		}
	}
	
	// Add the new pack size
	r.packSizes = append(r.packSizes, packSize)
	return nil
}

// Remove removes a pack size if it exists
func (r *MemoryPackSizeRepository) Remove(packSize int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Find and remove the pack size
	for i, size := range r.packSizes {
		if size == packSize {
			// Remove by creating a new slice without this element
			r.packSizes = append(r.packSizes[:i], r.packSizes[i+1:]...)
			return nil
		}
	}
	
	return ErrPackSizeNotFound
}

// Exists checks if a pack size exists
func (r *MemoryPackSizeRepository) Exists(packSize int) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	for _, size := range r.packSizes {
		if size == packSize {
			return true, nil
		}
	}
	
	return false, nil
}

// Errors
var (
	ErrInvalidPackSize  = &RepositoryError{Message: "pack size must be greater than zero"}
	ErrPackSizeNotFound = &RepositoryError{Message: "pack size not found"}
)

// RepositoryError represents a repository operation error
type RepositoryError struct {
	Message string
}

func (e *RepositoryError) Error() string {
	return e.Message
}

