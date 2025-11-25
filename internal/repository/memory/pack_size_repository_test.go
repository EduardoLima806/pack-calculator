package memory

import (
	"reflect"
	"testing"
)

func TestMemoryPackSizeRepository_GetAll(t *testing.T) {
	initialSizes := []int{250, 500, 1000}
	repo := NewMemoryPackSizeRepository(initialSizes)
	
	sizes, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}
	
	if !reflect.DeepEqual(sizes, initialSizes) {
		t.Errorf("GetAll() = %v, want %v", sizes, initialSizes)
	}
	
	// Verify that modifying the returned slice doesn't affect the repository
	sizes[0] = 999
	sizes2, _ := repo.GetAll()
	if sizes2[0] != 250 {
		t.Errorf("GetAll() returned slice was modified, repository should return a copy")
	}
}

func TestMemoryPackSizeRepository_SetAll(t *testing.T) {
	repo := NewMemoryPackSizeRepository([]int{250, 500})
	
	newSizes := []int{100, 200, 300}
	err := repo.SetAll(newSizes)
	if err != nil {
		t.Fatalf("SetAll() error = %v", err)
	}
	
	sizes, _ := repo.GetAll()
	if !reflect.DeepEqual(sizes, newSizes) {
		t.Errorf("SetAll() failed, GetAll() = %v, want %v", sizes, newSizes)
	}
}

func TestMemoryPackSizeRepository_Add(t *testing.T) {
	repo := NewMemoryPackSizeRepository([]int{250, 500})
	
	err := repo.Add(1000)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	
	sizes, _ := repo.GetAll()
	expected := []int{250, 500, 1000}
	if !reflect.DeepEqual(sizes, expected) {
		t.Errorf("Add() failed, GetAll() = %v, want %v", sizes, expected)
	}
	
	// Adding duplicate should not error
	err = repo.Add(250)
	if err != nil {
		t.Errorf("Add() duplicate should not error, got %v", err)
	}
	
	// Adding invalid pack size should error
	err = repo.Add(0)
	if err != ErrInvalidPackSize {
		t.Errorf("Add() with 0 should return ErrInvalidPackSize, got %v", err)
	}
	
	err = repo.Add(-1)
	if err != ErrInvalidPackSize {
		t.Errorf("Add() with -1 should return ErrInvalidPackSize, got %v", err)
	}
}

func TestMemoryPackSizeRepository_Remove(t *testing.T) {
	repo := NewMemoryPackSizeRepository([]int{250, 500, 1000})
	
	err := repo.Remove(500)
	if err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	
	sizes, _ := repo.GetAll()
	expected := []int{250, 1000}
	if !reflect.DeepEqual(sizes, expected) {
		t.Errorf("Remove() failed, GetAll() = %v, want %v", sizes, expected)
	}
	
	// Removing non-existent should error
	err = repo.Remove(999)
	if err != ErrPackSizeNotFound {
		t.Errorf("Remove() non-existent should return ErrPackSizeNotFound, got %v", err)
	}
}

func TestMemoryPackSizeRepository_Exists(t *testing.T) {
	repo := NewMemoryPackSizeRepository([]int{250, 500, 1000})
	
	exists, err := repo.Exists(500)
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Errorf("Exists(500) = false, want true")
	}
	
	exists, err = repo.Exists(999)
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Errorf("Exists(999) = true, want false")
	}
}

func TestMemoryPackSizeRepository_ThreadSafety(t *testing.T) {
	repo := NewMemoryPackSizeRepository([]int{250, 500})
	
	// Test concurrent access
	done := make(chan bool)
	
	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				repo.GetAll()
			}
			done <- true
		}()
	}
	
	// Concurrent writes
	for i := 0; i < 5; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				repo.Add(1000 + j)
				repo.Remove(1000 + j)
			}
			done <- true
		}()
	}
	
	// Wait for all goroutines
	for i := 0; i < 15; i++ {
		<-done
	}
	
	// Verify repository is still in valid state
	sizes, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll() after concurrent access error = %v", err)
	}
	
	// Should still have original sizes (or close to it, depending on timing)
	if len(sizes) < 2 {
		t.Errorf("Repository state corrupted, sizes = %v", sizes)
	}
}

