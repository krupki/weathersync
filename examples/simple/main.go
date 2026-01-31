// Simple example: Fetch weather for a single location
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/krupki/weathersync"
)

func main() {
	// Create a new client
	client := weathersync.New()

	// Define a location
	berlin := weathersync.Location{
		Name:      "Berlin",
		Latitude:  52.52,
		Longitude: 13.41,
	}

	// Fetch weather data
	data, err := client.FetchWeather(context.Background(), berlin)
	if err != nil {
		log.Fatal(err)
	}

	// Display result
	fmt.Printf("%s: %.1f°C\n", data.Location.Name, data.Temperature)
	fmt.Printf("Fetched in: %v\n", data.FetchDuration)
}
