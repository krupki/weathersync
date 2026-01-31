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

// WeatherData contains comprehensive weather information for a specific location.
// All fields are populated by the Open-Meteo API.
type WeatherData struct {
	// Location is the geographic location this data applies to
	Location Location

	// Temperature is the current temperature in Celsius
	Temperature float64

	// ApparentTemperature is how the temperature "feels" in Celsius
	ApparentTemperature float64

	// Humidity is the relative humidity as a percentage (0-100)
	Humidity float64

	// Precipitation is the rainfall in millimeters
	Precipitation float64

	// WeatherCode is the WMO weather interpretation code
	// See https://open-meteo.com/en/docs for code descriptions
	WeatherCode int

	// WindSpeed is the wind speed at 10 meters height in km/h
	WindSpeed float64

	// WindDirection is the wind direction at 10 meters height in degrees (0-360)
	WindDirection float64

	// WindGusts is the maximum wind gust speed in km/h
	WindGusts float64

	// CloudCover is the total cloud coverage as a percentage (0-100)
	CloudCover float64

	// Visibility is the visibility distance in meters
	Visibility float64

	// Pressure is the atmospheric pressure at mean sea level in hPa
	Pressure float64

	// FetchDuration is the time it took to fetch this data
	FetchDuration time.Duration

	// Timestamp is when this data was fetched
	Timestamp time.Time

	// Error contains any error that occurred during fetching
	// If nil, the fetch was successful
	Error error
}
