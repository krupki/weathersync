// Context example: Using context for timeout and cancellation
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/krupki/weathersync"
)

func main() {
	client := weathersync.New()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	location := weathersync.Location{
		Name:      "Paris",
		Latitude:  48.85,
		Longitude: 2.35,
	}

	fmt.Println("Fetching with 2-second timeout...")

	data, err := client.FetchWeather(ctx, location)
	if err != nil {
		log.Fatalf("Failed to fetch weather: %v", err)
	}

	fmt.Printf("%s: %.1f°C\n", data.Location.Name, data.Temperature)
	fmt.Printf("Fetch completed in: %v\n", data.FetchDuration)
}
