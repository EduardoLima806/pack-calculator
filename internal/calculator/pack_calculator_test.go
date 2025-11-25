package calculator

import (
	"reflect"
	"testing"
)

func TestPackCalculator_CalculatePacks(t *testing.T) {
	// Default pack sizes from the problem
	defaultPacks := []int{250, 500, 1000, 2000, 5000}
	
	tests := []struct {
		name           string
		packSizes      []int
		orderQuantity  int
		expectedResult map[int]int
		expectedItems  int
		expectedPacks  int
	}{
		{
			name:          "Order of 250 - should return 1x250",
			packSizes:     defaultPacks,
			orderQuantity: 250,
			expectedResult: map[int]int{250: 1},
			expectedItems: 250,
			expectedPacks: 1,
		},
		{
			name:          "Order of 251 - should return 1x500 (Rule 2: minimize items)",
			packSizes:     defaultPacks,
			orderQuantity: 251,
			expectedResult: map[int]int{500: 1},
			expectedItems: 500,
			expectedPacks: 1,
		},
		{
			name:          "Order of 501 - should return 1x500 + 1x250",
			packSizes:     defaultPacks,
			orderQuantity: 501,
			expectedResult: map[int]int{500: 1, 250: 1},
			expectedItems: 750,
			expectedPacks: 2,
		},
		{
			name:          "Order of 1000 - should return 1x1000",
			packSizes:     defaultPacks,
			orderQuantity: 1000,
			expectedResult: map[int]int{1000: 1},
			expectedItems: 1000,
			expectedPacks: 1,
		},
		{
			name:          "Order of 12001 - should return 2x5000 + 1x2000 + 1x250",
			packSizes:     defaultPacks,
			orderQuantity: 12001,
			expectedResult: map[int]int{5000: 2, 2000: 1, 250: 1},
			expectedItems: 12250,
			expectedPacks: 4,
		},
		{
			name:          "Order of 1 - should return smallest pack",
			packSizes:     defaultPacks,
			orderQuantity: 1,
			expectedResult: map[int]int{250: 1},
			expectedItems: 250,
			expectedPacks: 1,
		},
		{
			name:          "Order of 0 - should return empty",
			packSizes:     defaultPacks,
			orderQuantity: 0,
			expectedResult: map[int]int{},
			expectedItems: 0,
			expectedPacks: 0,
		},
		{
			name:          "Edge case: Pack sizes 23, 31, 53, Order 500000",
			packSizes:     []int{23, 31, 53},
			orderQuantity: 500000,
			expectedResult: map[int]int{23: 2, 31: 7, 53: 9429},
			expectedItems: 500000,
			expectedPacks: 9438,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewPackCalculator(tt.packSizes)
			result := calc.CalculatePacks(tt.orderQuantity)

			// Verify the result structure
			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("CalculatePacks() = %v, want %v", result, tt.expectedResult)
			}

			// Verify total items
			totalItems := calc.totalItems(result)
			if totalItems != tt.expectedItems {
				t.Errorf("Total items = %d, want %d", totalItems, tt.expectedItems)
			}

			// Verify total packs
			totalPacks := calc.totalPacks(result)
			if totalPacks != tt.expectedPacks {
				t.Errorf("Total packs = %d, want %d", totalPacks, tt.expectedPacks)
			}

			// Verify that we meet or exceed the order quantity
			if totalItems < tt.orderQuantity {
				t.Errorf("Total items %d is less than order quantity %d", totalItems, tt.orderQuantity)
			}
		})
	}
}

func TestPackCalculator_EdgeCase(t *testing.T) {
	// Specific edge case from requirements
	packSizes := []int{23, 31, 53}
	orderQuantity := 500000
	expectedResult := map[int]int{23: 2, 31: 7, 53: 9429}

	calc := NewPackCalculator(packSizes)
	result := calc.CalculatePacks(orderQuantity)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Edge case failed: got %v, want %v", result, expectedResult)
	}

	// Verify the calculation
	totalItems := 23*2 + 31*7 + 53*9429
	if totalItems != 500000 {
		t.Errorf("Edge case verification failed: total items = %d, want 500000", totalItems)
	}
}

func TestNewPackCalculator(t *testing.T) {
	packSizes := []int{5000, 250, 1000, 500, 2000} // Unsorted
	calc := NewPackCalculator(packSizes)

	// Verify that pack sizes are sorted in descending order
	expectedOrder := []int{5000, 2000, 1000, 500, 250}
	if !reflect.DeepEqual(calc.packSizes, expectedOrder) {
		t.Errorf("Pack sizes not sorted correctly: got %v, want %v", calc.packSizes, expectedOrder)
	}
}

func TestPackCalculator_Rule2Precedence(t *testing.T) {
	// Test that Rule 2 (minimize items) takes precedence over Rule 3 (minimize packs)
	packSizes := []int{250, 500}
	calc := NewPackCalculator(packSizes)
	
	// Order of 251: 1x500 (500 items, 1 pack) is better than 2x250 (500 items, 2 packs)
	result := calc.CalculatePacks(251)
	
	if result[500] != 1 {
		t.Errorf("Rule 2 precedence failed: should prefer 1x500 over 2x250, got %v", result)
	}
	
	totalItems := calc.totalItems(result)
	if totalItems != 500 {
		t.Errorf("Total items should be 500, got %d", totalItems)
	}
}

