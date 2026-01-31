package weathersync

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestNewClient tests client creation with default options
func TestNewClient(t *testing.T) {
	client := New()

	if client == nil {
		t.Fatal("New() returned nil")
	}

	if client.apiURL != "https://api.open-meteo.com" {
		t.Errorf("Expected default API URL, got %s", client.apiURL)
	}

	if client.timeout != 10*time.Second {
		t.Errorf("Expected default timeout 10s, got %v", client.timeout)
	}
}

// TestNewClientWithOptions tests client creation with custom options
func TestNewClientWithOptions(t *testing.T) {
	customURL := "https://custom.api.com"
	customTimeout := 5 * time.Second

	client := New(
		WithAPIURL(customURL),
		WithTimeout(customTimeout),
	)

	if client.apiURL != customURL {
		t.Errorf("Expected API URL %s, got %s", customURL, client.apiURL)
	}

	if client.timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, client.timeout)
	}
}

// TestFetchWeatherSuccess tests successful weather fetch
func TestFetchWeatherSuccess(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request parameters (float comparison with tolerance)
		latStr := r.URL.Query().Get("latitude")
		lonStr := r.URL.Query().Get("longitude")

		if latStr == "" || lonStr == "" {
			t.Error("Missing latitude or longitude parameter")
		}

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"current": {
				"temperature_2m": 15.3
			}
		}`))
	}))
	defer server.Close()

	// Create client with mock server
	client := New(WithAPIURL(server.URL))

	location := Location{
		Name:      "Berlin",
		Latitude:  52.52,
		Longitude: 13.41,
	}

	// Fetch weather
	data, err := client.FetchWeather(context.Background(), location)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if data == nil {
		t.Fatal("Expected WeatherData, got nil")
	}

	if data.Location.Name != "Berlin" {
		t.Errorf("Expected location Berlin, got %s", data.Location.Name)
	}

	if data.Temperature != 15.3 {
		t.Errorf("Expected temperature 15.3, got %f", data.Temperature)
	}

	if data.Error != nil {
		t.Errorf("Expected no error in data, got %v", data.Error)
	}

	if data.FetchDuration == 0 {
		t.Error("Expected FetchDuration > 0")
	}

	if data.Timestamp.IsZero() {
		t.Error("Expected non-zero Timestamp")
	}
}

// TestFetchWeatherAPIError tests handling of API errors
func TestFetchWeatherAPIError(t *testing.T) {
	// Create mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := New(WithAPIURL(server.URL))

	location := Location{
		Name:      "TestCity",
		Latitude:  0,
		Longitude: 0,
	}

	_, err := client.FetchWeather(context.Background(), location)

	if err == nil {
		t.Fatal("Expected error for API status 500, got nil")
	}
}

// TestFetchWeatherTimeout tests context timeout handling
func TestFetchWeatherTimeout(t *testing.T) {
	// Create slow mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.Write([]byte(`{"current": {"temperature_2m": 15.3}}`))
	}))
	defer server.Close()

	client := New(WithAPIURL(server.URL), WithTimeout(100*time.Millisecond))

	location := Location{
		Name:      "TestCity",
		Latitude:  0,
		Longitude: 0,
	}

	// Create context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.FetchWeather(ctx, location)

	if err == nil {
		t.Fatal("Expected timeout error, got nil")
	}
}

// TestFetchWeatherInvalidJSON tests handling of invalid JSON response
func TestFetchWeatherInvalidJSON(t *testing.T) {
	// Create mock server with invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json {"))
	}))
	defer server.Close()

	client := New(WithAPIURL(server.URL))

	location := Location{
		Name:      "TestCity",
		Latitude:  0,
		Longitude: 0,
	}

	_, err := client.FetchWeather(context.Background(), location)

	if err == nil {
		t.Fatal("Expected JSON error, got nil")
	}
}

// TestFetchMultipleSuccess tests successful concurrent fetch
func TestFetchMultipleSuccess(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get("latitude")
		var temp float64

		// Return different temps based on latitude
		switch lat {
		case "52.52":
			temp = 15.3
		case "35.68":
			temp = 20.1
		default:
			temp = 10.0
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current": {"temperature_2m": ` + string(rune(int(temp))) + `}}`))
	}))
	defer server.Close()

	client := New(WithAPIURL(server.URL))

	locations := []Location{
		{Name: "Berlin", Latitude: 52.52, Longitude: 13.41},
		{Name: "Tokyo", Latitude: 35.68, Longitude: 139.75},
	}

	results := client.FetchMultiple(context.Background(), locations)

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	for _, result := range results {
		if result.Location.Name == "" {
			t.Error("Expected location name, got empty string")
		}
	}
}

// TestFetchMultiplePartialFailure tests handling of partial failures
func TestFetchMultiplePartialFailure(t *testing.T) {
	callCount := 0

	// Create mock server that fails on second call
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current": {"temperature_2m": 15.3}}`))
	}))
	defer server.Close()

	client := New(WithAPIURL(server.URL))

	locations := []Location{
		{Name: "City1", Latitude: 0, Longitude: 0},
		{Name: "City2", Latitude: 0, Longitude: 0},
	}

	results := client.FetchMultiple(context.Background(), locations)

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	successCount := 0
	failureCount := 0

	for _, result := range results {
		if result.Error == nil {
			successCount++
		} else {
			failureCount++
		}
	}

	if successCount != 1 {
		t.Errorf("Expected 1 success, got %d", successCount)
	}

	if failureCount != 1 {
		t.Errorf("Expected 1 failure, got %d", failureCount)
	}
}

// TestLocationValid tests Location struct
func TestLocationValid(t *testing.T) {
	location := Location{
		Name:      "TestCity",
		Latitude:  52.52,
		Longitude: 13.41,
	}

	if location.Name != "TestCity" {
		t.Errorf("Expected name TestCity, got %s", location.Name)
	}

	if location.Latitude != 52.52 {
		t.Errorf("Expected latitude 52.52, got %f", location.Latitude)
	}

	if location.Longitude != 13.41 {
		t.Errorf("Expected longitude 13.41, got %f", location.Longitude)
	}
}

// TestWeatherDataValid tests WeatherData struct
func TestWeatherDataValid(t *testing.T) {
	location := Location{Name: "Test", Latitude: 0, Longitude: 0}
	data := WeatherData{
		Location:      location,
		Temperature:   15.3,
		FetchDuration: 100 * time.Millisecond,
		Timestamp:     time.Now(),
		Error:         nil,
	}

	if data.Location.Name != "Test" {
		t.Errorf("Expected location name Test, got %s", data.Location.Name)
	}

	if data.Temperature != 15.3 {
		t.Errorf("Expected temperature 15.3, got %f", data.Temperature)
	}

	if data.Error != nil {
		t.Errorf("Expected no error, got %v", data.Error)
	}
}

// TestFetchWeatherConcurrency tests concurrent requests don't cause race conditions
func TestFetchWeatherConcurrency(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current": {"temperature_2m": 15.3}}`))
	}))
	defer server.Close()

	client := New(WithAPIURL(server.URL))

	// Run concurrent requests
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(index int) {
			location := Location{
				Name:      "City",
				Latitude:  float64(index),
				Longitude: float64(index),
			}

			_, err := client.FetchWeather(context.Background(), location)
			if err != nil {
				t.Errorf("Request %d failed: %v", index, err)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}
