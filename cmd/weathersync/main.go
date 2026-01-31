// Package main demonstrates how to use the weathersync library
// to fetch and compare weather data from multiple cities.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/krupki/weathersync"
	"gopkg.in/yaml.v2"
)

type config struct {
	Continents []struct {
		Name   string                 `yaml:"name"`
		Cities []weathersync.Location `yaml:"cities"`
	} `yaml:"continents"`
}

func main() {
	// Load configuration
	cfg, err := loadConfig("config/cities.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Create weathersync client
	client := weathersync.New(
		weathersync.WithTimeout(5 * time.Second),
	)

	// Prepare all locations
	var allLocations []weathersync.Location
	continentMap := make(map[string]string) // city name -> continent

	for _, cont := range cfg.Continents {
		for _, city := range cont.Cities {
			allLocations = append(allLocations, city)
			continentMap[city.Name] = cont.Name
		}
	}

	// Fetch weather data for all locations concurrently
	fmt.Println("Fetching weather data...")
	ctx := context.Background()
	results := client.FetchMultiple(ctx, allLocations)

	// Display results grouped by continent
	displayResults(results, continentMap)
}

func loadConfig(path string) (*config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func displayResults(results []weathersync.WeatherData, continentMap map[string]string) {
	// Group by continent
	groups := make(map[string][]weathersync.WeatherData)
	for _, r := range results {
		continent := continentMap[r.Location.Name]
		groups[continent] = append(groups[continent], r)
	}

	fmt.Println("\n========================================")
	fmt.Println("Weather Comparison Summary")
	fmt.Println("========================================")

	for continent, data := range groups {
		fmt.Printf("\nContinent: %s\n", continent)
		fmt.Println("----------------------------------------")

		var tempSum float64
		var validCount int

		for _, d := range data {
			if d.Error != nil {
				fmt.Printf("   %s: ERROR - %v (%.3fs)\n",
					d.Location.Name, d.Error, d.FetchDuration.Seconds())
			} else {
				fmt.Printf("   %s: %.1f°C (%.3fs)\n",
					d.Location.Name, d.Temperature, d.FetchDuration.Seconds())
				tempSum += d.Temperature
				validCount++
			}
		}

		if validCount > 0 {
			fmt.Printf("\n   Average: %.1f°C (%d cities)\n",
				tempSum/float64(validCount), validCount)
		}
	}

	fmt.Println("\n========================================")
}
