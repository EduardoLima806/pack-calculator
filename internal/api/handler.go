package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"pack-calculator/internal/calculator"
	"pack-calculator/internal/repository"
)

// Handler handles HTTP requests for the pack calculator API
type Handler struct {
	packSizeRepo repository.PackSizeRepository
	calc         *calculator.PackCalculator
}

// NewHandler creates a new API handler with a pack size repository
func NewHandler(packSizeRepo repository.PackSizeRepository) *Handler {
	// Get initial pack sizes from repository
	packSizes, err := packSizeRepo.GetAll()
	if err != nil {
		log.Printf("Error getting pack sizes from repository: %v", err)
		// Fallback to empty slice if repository fails
		packSizes = []int{}
	}
	
	return &Handler{
		packSizeRepo: packSizeRepo,
		calc:         calculator.NewPackCalculator(packSizes),
	}
}

// getCalculator returns a calculator with current pack sizes from repository
func (h *Handler) getCalculator() *calculator.PackCalculator {
	packSizes, err := h.packSizeRepo.GetAll()
	if err != nil {
		log.Printf("Error getting pack sizes from repository: %v", err)
		// Return existing calculator if repository fails
		return h.calc
	}
	
	// Create new calculator with current pack sizes
	return calculator.NewPackCalculator(packSizes)
}

// CalculateRequest represents a calculation request
type CalculateRequest struct {
	Quantity int   `json:"quantity"`
	PackSizes []int `json:"pack_sizes,omitempty"` // Optional: override default pack sizes
}

// CalculateResponse represents the calculation result
type CalculateResponse struct {
	Quantity    int            `json:"quantity"`
	Packs       map[int]int    `json:"packs"`
	TotalItems  int            `json:"total_items"`
	TotalPacks  int            `json:"total_packs"`
	PackSizes   []int          `json:"pack_sizes"`
}

// Calculate handles POST /api/calculate
func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Quantity < 0 {
		http.Error(w, "Quantity must be non-negative", http.StatusBadRequest)
		return
	}

	// Use custom pack sizes if provided, otherwise use pack sizes from repository
	var calc *calculator.PackCalculator
	if len(req.PackSizes) > 0 {
		calc = calculator.NewPackCalculator(req.PackSizes)
	} else {
		calc = h.getCalculator()
	}

	result := calc.CalculatePacks(req.Quantity)
	
	// Calculate totals
	totalItems := 0
	totalPacks := 0
	for packSize, count := range result {
		totalItems += packSize * count
		totalPacks += count
	}

	// Get pack sizes used
	packSizes := make([]int, 0, len(result))
	for packSize := range result {
		packSizes = append(packSizes, packSize)
	}

	response := CalculateResponse{
		Quantity:   req.Quantity,
		Packs:      result,
		TotalItems: totalItems,
		TotalPacks: totalPacks,
		PackSizes:  packSizes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPackSizes handles GET /api/pack-sizes
func (h *Handler) GetPackSizes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get pack sizes from repository
	packSizes, err := h.packSizeRepo.GetAll()
	if err != nil {
		http.Error(w, "Error retrieving pack sizes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"pack_sizes": packSizes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Health handles GET /health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// CalculateQuery handles GET /api/calculate?quantity=XXX
func (h *Handler) CalculateQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	quantityStr := r.URL.Query().Get("quantity")
	if quantityStr == "" {
		http.Error(w, "quantity parameter is required", http.StatusBadRequest)
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity < 0 {
		http.Error(w, "Invalid quantity parameter", http.StatusBadRequest)
		return
	}

	// Get calculator with current pack sizes from repository
	calc := h.getCalculator()
	result := calc.CalculatePacks(quantity)
	
	// Calculate totals
	totalItems := 0
	totalPacks := 0
	for packSize, count := range result {
		totalItems += packSize * count
		totalPacks += count
	}

	// Get pack sizes used
	packSizes := make([]int, 0, len(result))
	for packSize := range result {
		packSizes = append(packSizes, packSize)
	}

	response := CalculateResponse{
		Quantity:   quantity,
		Packs:      result,
		TotalItems: totalItems,
		TotalPacks: totalPacks,
		PackSizes:  packSizes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

