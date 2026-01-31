package weathersync

import "time"

// Location represents a geographic location with coordinates.
type Location struct {
	// Name is the display name of the location (e.g., "Berlin", "Tokyo")
	Name string

	// Latitude is the geographic latitude in decimal degrees
	Latitude float64

	// Longitude is the geographic longitude in decimal degrees
	Longitude float64
}

// WeatherData contains weather information for a specific location.
type WeatherData struct {
	// Location is the geographic location this data applies to
	Location Location

	// Temperature is the current temperature in Celsius
	Temperature float64

	// FetchDuration is the time it took to fetch this data
	FetchDuration time.Duration

	// Timestamp is when this data was fetched
	Timestamp time.Time

	// Error contains any error that occurred during fetching
	// If nil, the fetch was successful
	Error error
}
