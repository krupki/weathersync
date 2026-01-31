package weathersync

import (
	"testing"
	"time"
)

// TestLocationCreation tests Location struct creation
func TestLocationCreation(t *testing.T) {
	tests := []struct {
		name      string
		location  Location
		wantName  string
		wantLat   float64
		wantLon   float64
	}{
		{
			name:     "Berlin",
			location: Location{Name: "Berlin", Latitude: 52.52, Longitude: 13.41},
			wantName: "Berlin",
			wantLat:  52.52,
			wantLon:  13.41,
		},
		{
			name:     "Tokyo",
			location: Location{Name: "Tokyo", Latitude: 35.68, Longitude: 139.75},
			wantName: "Tokyo",
			wantLat:  35.68,
			wantLon:  139.75,
		},
		{
			name:     "Negative coordinates",
			location: Location{Name: "Sydney", Latitude: -33.87, Longitude: 151.21},
			wantName: "Sydney",
			wantLat:  -33.87,
			wantLon:  151.21,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.location.Name != tt.wantName {
				t.Errorf("Name = %s, want %s", tt.location.Name, tt.wantName)
			}
			if tt.location.Latitude != tt.wantLat {
				t.Errorf("Latitude = %f, want %f", tt.location.Latitude, tt.wantLat)
			}
			if tt.location.Longitude != tt.wantLon {
				t.Errorf("Longitude = %f, want %f", tt.location.Longitude, tt.wantLon)
			}
		})
	}
}

// TestWeatherDataCreation tests WeatherData struct creation
func TestWeatherDataCreation(t *testing.T) {
	location := Location{Name: "Berlin", Latitude: 52.52, Longitude: 13.41}
	timestamp := time.Now()

	data := WeatherData{
		Location:      location,
		Temperature:   15.3,
		FetchDuration: 125 * time.Millisecond,
		Timestamp:     timestamp,
		Error:         nil,
	}

	if data.Location.Name != "Berlin" {
		t.Errorf("Location.Name = %s, want Berlin", data.Location.Name)
	}

	if data.Temperature != 15.3 {
		t.Errorf("Temperature = %f, want 15.3", data.Temperature)
	}

	if data.FetchDuration != 125*time.Millisecond {
		t.Errorf("FetchDuration = %v, want 125ms", data.FetchDuration)
	}

	if data.Timestamp != timestamp {
		t.Errorf("Timestamp mismatch")
	}

	if data.Error != nil {
		t.Errorf("Error = %v, want nil", data.Error)
	}
}

// TestWeatherDataWithError tests WeatherData with error
func TestWeatherDataWithError(t *testing.T) {
	location := Location{Name: "Test", Latitude: 0, Longitude: 0}
	testErr := "test error"

	data := WeatherData{
		Location:      location,
		Temperature:   0,
		FetchDuration: 0,
		Timestamp:     time.Now(),
		Error:         &TestError{msg: testErr},
	}

	if data.Error == nil {
		t.Fatal("Expected error, got nil")
	}

	if data.Error.Error() != testErr {
		t.Errorf("Error message = %s, want %s", data.Error.Error(), testErr)
	}
}

// TestError helper for testing
type TestError struct {
	msg string
}

func (e *TestError) Error() string {
	return e.msg
}

// TestLocationBoundaries tests extreme coordinate values
func TestLocationBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		location Location
	}{
		{
			name:     "North Pole",
			location: Location{Name: "North Pole", Latitude: 90, Longitude: 0},
		},
		{
			name:     "South Pole",
			location: Location{Name: "South Pole", Latitude: -90, Longitude: 180},
		},
		{
			name:     "Prime Meridian",
			location: Location{Name: "Prime Meridian", Latitude: 51.5, Longitude: 0},
		},
		{
			name:     "International Date Line",
			location: Location{Name: "Date Line", Latitude: 0, Longitude: 180},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify the values are set correctly
			if tt.location.Latitude < -90 || tt.location.Latitude > 90 {
				t.Errorf("Latitude out of bounds: %f", tt.location.Latitude)
			}

			if tt.location.Longitude < -180 || tt.location.Longitude > 180 {
				t.Errorf("Longitude out of bounds: %f", tt.location.Longitude)
			}
		})
	}
}

// TestWeatherDataTemperatureRange tests realistic temperature values
func TestWeatherDataTemperatureRange(t *testing.T) {
	tests := []struct {
		name        string
		temperature float64
		valid       bool
	}{
		{name: "Freezing", temperature: 0, valid: true},
		{name: "Room temperature", temperature: 20, valid: true},
		{name: "Very hot", temperature: 50, valid: true},
		{name: "Very cold", temperature: -40, valid: true},
		{name: "Extreme hot", temperature: 60, valid: true},
		{name: "Extreme cold", temperature: -60, valid: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := WeatherData{
				Temperature: tt.temperature,
			}

			if data.Temperature != tt.temperature {
				t.Errorf("Temperature = %f, want %f", data.Temperature, tt.temperature)
			}
		})
	}
}
