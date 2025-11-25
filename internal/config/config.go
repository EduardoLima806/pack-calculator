package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds the application configuration
type Config struct {
	Port      string
	PackSizes []int
}

// LoadConfig loads configuration from environment variables
// Default pack sizes: 250, 500, 1000, 2000, 5000
func LoadConfig() *Config {
	cfg := &Config{
		Port:      getEnv("PORT", "8080"),
		PackSizes: []int{250, 500, 1000, 2000, 5000},
	}

	// Override pack sizes from environment if provided
	if packSizesStr := os.Getenv("PACK_SIZES"); packSizesStr != "" {
		packSizes := parsePackSizes(packSizesStr)
		if len(packSizes) > 0 {
			cfg.PackSizes = packSizes
		}
	}

	return cfg
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parsePackSizes parses a comma-separated string of pack sizes
func parsePackSizes(s string) []int {
	parts := strings.Split(s, ",")
	packSizes := make([]int, 0, len(parts))
	
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if size, err := strconv.Atoi(part); err == nil && size > 0 {
			packSizes = append(packSizes, size)
		}
	}
	
	return packSizes
}

