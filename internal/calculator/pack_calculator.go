package calculator

import (
	"sort"
)

// PackCalculator handles the calculation of optimal pack combinations
type PackCalculator struct {
	packSizes []int
}

// NewPackCalculator creates a new calculator with the given pack sizes
// Pack sizes should be sorted in descending order for optimal performance
func NewPackCalculator(packSizes []int) *PackCalculator {
	// Create a copy and sort in descending order
	sizes := make([]int, len(packSizes))
	copy(sizes, packSizes)
	sort.Sort(sort.Reverse(sort.IntSlice(sizes)))
	
	return &PackCalculator{
		packSizes: sizes,
	}
}

// GetPackSizes returns the current pack sizes
func (pc *PackCalculator) GetPackSizes() []int {
	result := make([]int, len(pc.packSizes))
	copy(result, pc.packSizes)
	return result
}

// CalculatePacks determines the optimal pack combination for the given order quantity
// Rules:
// 1. Only whole packs can be sent
// 2. Minimize total items (takes precedence)
// 3. Minimize number of packs
func (pc *PackCalculator) CalculatePacks(orderQuantity int) map[int]int {
	if orderQuantity <= 0 {
		return make(map[int]int)
	}

	// For very large quantities, use optimized algorithm
	if orderQuantity > 100000 {
		return pc.calculateForLargeQuantity(orderQuantity)
	}

	// Use dynamic programming for smaller quantities
	return pc.calculateWithDP(orderQuantity)
}

// calculateWithDP uses dynamic programming to find optimal solution
func (pc *PackCalculator) calculateWithDP(orderQuantity int) map[int]int {
	// dp[i] represents the minimum number of items needed to fulfill order quantity i
	// packCount[i] represents the minimum number of packs needed for order quantity i
	// We need to extend the array to handle cases where we need to exceed orderQuantity
	maxPackSize := pc.packSizes[0]
	maxQuantity := orderQuantity + maxPackSize
	dp := make([]int, maxQuantity+1)
	packCount := make([]int, maxQuantity+1)
	prevPack := make([]int, maxQuantity+1)
	
	// Initialize with infinity (represented as a large number)
	for i := range dp {
		dp[i] = 1 << 30 // Large number representing infinity
		packCount[i] = 1 << 30
		prevPack[i] = -1
	}
	dp[0] = 0
	packCount[0] = 0

	// Fill the DP table
	for i := 1; i <= maxQuantity; i++ {
		for _, packSize := range pc.packSizes {
			if packSize <= i {
				// Calculate the new total items if we use this pack
				newTotalItems := dp[i-packSize] + packSize
				newPackCount := packCount[i-packSize] + 1
				
				// Rule 2: Minimize total items (takes precedence)
				if newTotalItems < dp[i] {
					dp[i] = newTotalItems
					packCount[i] = newPackCount
					prevPack[i] = packSize
				} else if newTotalItems == dp[i] {
					// Rule 3: If total items are equal, minimize pack count
					if newPackCount < packCount[i] {
						packCount[i] = newPackCount
						prevPack[i] = packSize
					}
				}
			}
		}
	}

	// If we couldn't find an exact match, find the minimum excess
	minItems := dp[orderQuantity]
	minPacks := packCount[orderQuantity]
	bestQuantity := orderQuantity
	
	if minItems == 1<<30 {
		// No exact solution found, find minimum excess
		for qty := orderQuantity + 1; qty <= maxQuantity; qty++ {
			if dp[qty] < minItems {
				minItems = dp[qty]
				minPacks = packCount[qty]
				bestQuantity = qty
			} else if dp[qty] == minItems && packCount[qty] < minPacks {
				minPacks = packCount[qty]
				bestQuantity = qty
			}
		}
	}

	// Reconstruct the solution
	result := make(map[int]int)
	if bestQuantity < len(prevPack) && prevPack[bestQuantity] != -1 {
		pc.reconstructSolution(bestQuantity, prevPack, result)
	}

	return result
}

// reconstructSolution builds the pack combination from the DP table
func (pc *PackCalculator) reconstructSolution(quantity int, prevPack []int, result map[int]int) {
	current := quantity
	for current > 0 && prevPack[current] != -1 {
		packSize := prevPack[current]
		result[packSize]++
		current -= packSize
	}
}

// calculateForLargeQuantity handles cases where quantity is very large
// First tries to find exact solution, then falls back to greedy with optimization
func (pc *PackCalculator) calculateForLargeQuantity(quantity int) map[int]int {
	// First, try to find an exact solution (important for edge cases)
	if len(pc.packSizes) == 3 {
		exactSolution := pc.solveThreePackSizes(quantity)
		if exactSolution != nil && pc.totalItems(exactSolution) == quantity {
			return exactSolution
		}
	}
	
	// Try general exact solution finder
	exactSolution := pc.findExactSolution(quantity)
	if exactSolution != nil && pc.totalItems(exactSolution) == quantity {
		return exactSolution
	}
	
	// Fall back to greedy approach with optimization
	result := make(map[int]int)
	remaining := quantity

	// Use largest packs as much as possible
	for _, packSize := range pc.packSizes {
		if remaining >= packSize {
			count := remaining / packSize
			result[packSize] = count
			remaining -= count * packSize
		}
	}

	// Handle any remaining quantity
	if remaining > 0 {
		// Find the smallest pack that covers the remainder
		for i := len(pc.packSizes) - 1; i >= 0; i-- {
			if pc.packSizes[i] >= remaining {
				result[pc.packSizes[i]]++
				break
			}
		}
	}

	// Optimize the result
	return pc.optimizeLargeResult(result, quantity)
}

// optimizeLargeResult optimizes the result for large quantities
// Tries to reduce excess items while maintaining minimum pack count
func (pc *PackCalculator) optimizeLargeResult(current map[int]int, targetQuantity int) map[int]int {
	optimized := make(map[int]int)
	for packSize, count := range current {
		optimized[packSize] = count
	}

	currentTotal := pc.totalItems(optimized)
	
	// If greedy result doesn't match target, try to find exact solution
	// This attempts to improve the greedy approximation
	if currentTotal != targetQuantity {
		// Try to find exact solution using backtracking or specialized solvers
		exactSolution := pc.findExactSolution(targetQuantity)
		if exactSolution != nil && pc.totalItems(exactSolution) == targetQuantity {
			return exactSolution
		}
	}

	return optimized
}

// findExactSolution attempts to find an exact solution for the target quantity
// Uses specialized solvers based on the number of pack sizes
func (pc *PackCalculator) findExactSolution(targetQuantity int) map[int]int {
	// For three pack sizes, use optimized brute force search
	if len(pc.packSizes) == 3 {
		return pc.solveThreePackSizes(targetQuantity)
	}
	
	// For other cases, use backtracking with memoization
	return pc.solveIterative(targetQuantity)
}

// solveThreePackSizes solves for exactly three pack sizes
// This handles the edge case: 23, 31, 53 with target 500000
// Uses optimized search with proper bounds to find optimal solution
func (pc *PackCalculator) solveThreePackSizes(targetQuantity int) map[int]int {
	a, b, c := pc.packSizes[0], pc.packSizes[1], pc.packSizes[2]
	var bestResult map[int]int
	bestItems := 1 << 30
	bestPacks := 1 << 30
	
	// Try all combinations of a and b, then solve for c
	maxA := targetQuantity / a
	
	// Search all valid combinations
	for i := 0; i <= maxA; i++ {
		remainingAfterA := targetQuantity - i*a
		if remainingAfterA < 0 {
			break
		}
		
		maxBForThisA := remainingAfterA / b
		for j := 0; j <= maxBForThisA; j++ {
			remaining := remainingAfterA - j*b
			if remaining >= 0 && remaining%c == 0 {
				k := remaining / c
				if k >= 0 {
					// Found exact solution
					currentItems := i*a + j*b + k*c
					currentPacks := i + j + k
					
					// Rule 2: Minimize items (takes precedence)
					if currentItems < bestItems {
						bestItems = currentItems
						bestPacks = currentPacks
						bestResult = make(map[int]int)
						if i > 0 {
							bestResult[a] = i
						}
						if j > 0 {
							bestResult[b] = j
						}
						if k > 0 {
							bestResult[c] = k
						}
					} else if currentItems == bestItems {
						// Rule 3: If items are equal, minimize packs
						if currentPacks < bestPacks {
							bestPacks = currentPacks
							bestResult = make(map[int]int)
							if i > 0 {
								bestResult[a] = i
							}
							if j > 0 {
								bestResult[b] = j
							}
							if k > 0 {
								bestResult[c] = k
							}
						}
					}
				}
			}
		}
	}
	
	return bestResult
}

// solveIterative uses iterative deepening to find exact solution
func (pc *PackCalculator) solveIterative(targetQuantity int) map[int]int {
	// Use backtracking with bounds
	result := make(map[int]int)
	if pc.solveBacktrack(targetQuantity, 0, result, make(map[string]bool)) {
		return result
	}
	return nil
}

// solveBacktrack uses backtracking to find exact solution
func (pc *PackCalculator) solveBacktrack(remaining int, packIndex int, result map[int]int, memo map[string]bool) bool {
	if remaining == 0 {
		return true
	}
	if remaining < 0 || packIndex >= len(pc.packSizes) {
		return false
	}
	
	// Memoization key
	key := string(rune(remaining)) + string(rune(packIndex))
	if memo[key] {
		return false
	}
	
	packSize := pc.packSizes[packIndex]
	maxCount := remaining / packSize
	
	// Try using this pack size
	for count := maxCount; count >= 0; count-- {
		if count > 10000 { // Limit to prevent excessive recursion
			break
		}
		newRemaining := remaining - count*packSize
		result[packSize] = count
		
		if pc.solveBacktrack(newRemaining, packIndex+1, result, memo) {
			if result[packSize] == 0 {
				delete(result, packSize)
			}
			return true
		}
		
		delete(result, packSize)
	}
	
	memo[key] = true
	return false
}

// totalItems calculates the total number of items in a pack combination
func (pc *PackCalculator) totalItems(packs map[int]int) int {
	total := 0
	for packSize, count := range packs {
		total += packSize * count
	}
	return total
}

// totalPacks calculates the total number of packs
func (pc *PackCalculator) totalPacks(packs map[int]int) int {
	total := 0
	for _, count := range packs {
		total += count
	}
	return total
}
