// Concurrent example: Fetch weather for multiple locations in parallel
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/krupki/weathersync"
)

func main() {
	// Create a new client with custom timeout
	client := weathersync.New(
		weathersync.WithTimeout(5 * time.Second),
	)

	// Define multiple locations
	locations := []weathersync.Location{
		{Name: "Berlin", Latitude: 52.52, Longitude: 13.41},
		{Name: "Tokyo", Latitude: 35.68, Longitude: 139.75},
		{Name: "New York", Latitude: 40.71, Longitude: -74.01},
		{Name: "Sydney", Latitude: -33.87, Longitude: 151.21},
		{Name: "São Paulo", Latitude: -23.55, Longitude: -46.63},
	}

	// Fetch all weather data concurrently
	start := time.Now()
	results := client.FetchMultiple(context.Background(), locations)
	elapsed := time.Since(start)

	// Display results
	fmt.Println("Weather Report")
	fmt.Println("==============")
	for _, data := range results {
		if data.Error != nil {
			fmt.Printf("%s: %v\n", data.Location.Name, data.Error)
		} else {
			fmt.Printf("%s: %.1f°C (fetched in %.3fs)\n",
				data.Location.Name,
				data.Temperature,
				data.FetchDuration.Seconds())
		}
	}

	fmt.Printf("\nTotal time: %.3fs (parallel fetching)\n", elapsed.Seconds())
}
