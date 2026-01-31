// Package weathersync provides a simple, concurrent weather data fetching library.
// It allows you to retrieve current weather information for multiple locations
// in parallel using goroutines.
//
// Basic usage:
//
//	client := weathersync.New()
//	location := weathersync.Location{
//		Name:      "Berlin",
//		Latitude:  52.52,
//		Longitude: 13.41,
//	}
//	data, err := client.FetchWeather(context.Background(), location)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("%s: %.1f°C\n", data.Location.Name, data.Temperature)
package weathersync

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Client is the main entry point for the weathersync library.
// It provides methods to fetch weather data for single or multiple locations.
type Client struct {
	apiURL     string
	httpClient *http.Client
	timeout    time.Duration
}

// Option is a function that configures a Client.
type Option func(*Client)

// WithTimeout sets a custom timeout for HTTP requests.
// Default is 10 seconds.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.timeout = d
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithAPIURL sets a custom weather API URL.
// Default is "https://api.open-meteo.com".
func WithAPIURL(url string) Option {
	return func(c *Client) {
		c.apiURL = url
	}
}

// New creates a new weathersync Client with the given options.
// If no options are provided, sensible defaults are used.
func New(opts ...Option) *Client {
	c := &Client{
		apiURL:     "https://api.open-meteo.com",
		httpClient: &http.Client{},
		timeout:    10 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}

	// Apply timeout to HTTP client
	c.httpClient.Timeout = c.timeout

	return c
}

// FetchWeather retrieves current weather data for a single location.
// It respects the context for cancellation and timeouts.
//
// Parameters:
//   - ctx: context for cancellation and timeout control
//   - location: the geographic location to fetch weather for
//
// Returns:
//   - *WeatherData containing temperature and metadata
//   - error if the request fails or data is invalid
func (c *Client) FetchWeather(ctx context.Context, location Location) (*WeatherData, error) {
	url := fmt.Sprintf("%s/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m",
		c.apiURL, location.Latitude, location.Longitude)

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var apiResp struct {
		Current struct {
			Temperature2M float64 `json:"temperature_2m"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &WeatherData{
		Location:      location,
		Temperature:   apiResp.Current.Temperature2M,
		FetchDuration: time.Since(start),
		Timestamp:     time.Now(),
	}, nil
}

// FetchMultiple retrieves weather data for multiple locations concurrently.
// All requests are performed in parallel using goroutines.
//
// Parameters:
//   - ctx: context for cancellation and timeout control
//   - locations: slice of locations to fetch weather for
//
// Returns:
//   - []WeatherData slice containing results for all locations
//   - Any errors are embedded in the individual WeatherData.Error field
func (c *Client) FetchMultiple(ctx context.Context, locations []Location) []WeatherData {
	results := make([]WeatherData, len(locations))
	var wg sync.WaitGroup

	for i, loc := range locations {
		wg.Add(1)
		go func(index int, location Location) {
			defer wg.Done()

			data, err := c.FetchWeather(ctx, location)
			if err != nil {
				results[index] = WeatherData{
					Location: location,
					Error:    err,
				}
				return
			}
			results[index] = *data
		}(i, loc)
	}

	wg.Wait()
	return results
}
