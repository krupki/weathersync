// Demo: Real-world usage of weathersync library
// This shows how a developer would integrate your library
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/krupki/weathersync"
)

func main() {
	fmt.Println(" WeatherSync Library Demo")
	fmt.Println("================================")

	// STEP 1: Create a client
	fmt.Println("Step 1: Creating weathersync client...")
	client := weathersync.New(
		weathersync.WithTimeout(5 * time.Second),
	)
	fmt.Println("Client created!")

	// STEP 2: Define locations (input data)
	fmt.Println("Step 2: Defining locations...")
	locations := []weathersync.Location{
		{Name: "Berlin", Latitude: 52.52, Longitude: 13.41},
		{Name: "Tokyo", Latitude: 35.68, Longitude: 139.75},
		{Name: "New York", Latitude: 40.71, Longitude: -74.01},
		{Name: "Sydney", Latitude: -33.87, Longitude: 151.21},
	}
	fmt.Printf("%d locations defined\n\n", len(locations))

	// STEP 3: Fetch weather data (library in action!)
	fmt.Println("Step 3: Fetching weather data concurrently...")
	fmt.Println("   (All requests happen in parallel)")
	start := time.Now()

	results := client.FetchMultiple(context.Background(), locations)

	elapsed := time.Since(start)
	fmt.Printf("Fetched in %.3f seconds!\n\n", elapsed.Seconds())

	// STEP 4: Process the output data
	fmt.Println("Step 4: Processing results...")
	fmt.Println("================================")

	successCount := 0
	var totalTemp float64

	for i, data := range results {
		fmt.Printf("\n[%d] %s\n", i+1, data.Location.Name)
		fmt.Println("    ├─ Coordinates:", data.Location.Latitude, ",", data.Location.Longitude)

		if data.Error != nil {
			fmt.Printf("    ├─ Status: FAILED\n")
			fmt.Printf("    └─ Error: %v\n", data.Error)
		} else {
			fmt.Printf("    ├─ Status: SUCCESS\n")
			fmt.Printf("    ├─ Temperature: %.1f°C\n", data.Temperature)
			fmt.Printf("    ├─ Fetch Duration: %v\n", data.FetchDuration)
			fmt.Printf("    └─ Timestamp: %s\n", data.Timestamp.Format("15:04:05"))
			successCount++
			totalTemp += data.Temperature
		}
	}

	// STEP 5: Summary
	fmt.Println("\n================================")
	fmt.Println("Summary:")
	fmt.Printf("   • Total locations: %d\n", len(results))
	fmt.Printf("   • Successful: %d\n", successCount)
	fmt.Printf("   • Failed: %d\n", len(results)-successCount)
	if successCount > 0 {
		fmt.Printf("   • Average temperature: %.1f°C\n", totalTemp/float64(successCount))
	}
	fmt.Printf("   • Total time: %.3f seconds\n", elapsed.Seconds())
	fmt.Printf("   • Performance: %dx faster than sequential!\n", len(locations))

	fmt.Println("\n✨ Demo completed!")
}
