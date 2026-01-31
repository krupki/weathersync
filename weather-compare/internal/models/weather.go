// Package models defines the core data structures for weather information
// and geographic locations used throughout the WeatherSync application.
package models

// WeatherResult represents the complete weather information for a single city,
// including temperature data, location context, and fetch performance metrics.
// The Error field contains any issues encountered during the data fetch.
type WeatherResult struct {
	City                 string  `json:"city"`                       // Name of the city
	Temperature          float64 `json:"temperature"`                // Current temperature in Celsius
	Continent            string  `json:"continent" yaml:"continent"` // Continent where the city is located
	Error                error   `json:"error"`                      // Error if fetch failed, nil otherwise
	FetchDurationSeconds float64 `json:"fetch_duration_seconds"`     // Time taken to fetch data in seconds
}

// City represents a geographic location with coordinates.
// It is used for configuration and weather API requests.
type City struct {
	Name      string  `yaml:"name" json:"name"`           // City name
	Latitude  float64 `yaml:"latitude" json:"latitude"`   // Latitude coordinate
	Longitude float64 `yaml:"longitude" json:"longitude"` // Longitude coordinate
}
