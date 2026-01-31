// Package main implements a concurrent weather data aggregator that fetches
// and compares real-time weather information across multiple continents.
//
// The application demonstrates production-level Go patterns including:
//   - Concurrent data fetching using goroutines and channels
//   - Per-continent worker isolation for better resource utilization
//   - Configuration-driven city management via YAML
//   - Performance tracking with latency measurements
//
// Weather data is fetched in parallel for all cities, grouped by continent,
// with comprehensive error handling and performance metrics.
package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"weather-compare/internal/comparer"
	"weather-compare/internal/fetcher"
	"weather-compare/internal/models"

	"gopkg.in/yaml.v2"
)

func main() {
	cfg, err := loadCities("../../config/cities.yaml")
	if err != nil {
		log.Fatalf("Error loading cities: %v", err)
	}

	// prepare channels and counters
	total := 0
	for _, c := range cfg.Continents {
		total += len(c.Cities)
	}
	resultCh := make(chan models.WeatherResult, total)
	contDoneCh := make(chan struct {
		Name     string
		Duration time.Duration
	}, len(cfg.Continents))

	// launch per-continent workers
	for _, cont := range cfg.Continents {
		var wg sync.WaitGroup
		start := time.Now()
		for _, city := range cont.Cities {
			wg.Add(1)
			go func(ct models.City) {
				defer wg.Done()
				res, err := fetcher.FetchWeatherData(ct.Name, ct.Latitude, ct.Longitude)
				if err != nil {
					log.Printf("Error fetching %s: %v", ct.Name, err)
				}
				res.Continent = cont.Name
				resultCh <- res
			}(city)
		}

		// report when this continent's cities are done
		go func(name string, w *sync.WaitGroup, s time.Time) {
			w.Wait()
			contDoneCh <- struct {
				Name     string
				Duration time.Duration
			}{Name: name, Duration: time.Since(s)}
		}(cont.Name, &wg, start)
	}

	// collect results
	weatherData := make([]models.WeatherResult, 0, total)
	for i := 0; i < total; i++ {
		r := <-resultCh
		weatherData = append(weatherData, r)
	}

	// collect continent durations
	contDurations := make(map[string]time.Duration)
	for i := 0; i < len(cfg.Continents); i++ {
		d := <-contDoneCh
		contDurations[d.Name] = d.Duration
	}
	close(resultCh)
	close(contDoneCh)

	comparisonSummary := comparer.CompareWeatherData(weatherData, contDurations)
	fmt.Println(comparisonSummary)
}

// loadCities reads and parses the YAML configuration file containing city definitions.
// It returns a structure with continents and their associated cities.
//
// Parameters:
//   - filePath: path to the YAML configuration file
//
// Returns:
//   - A struct containing continents and their cities
//   - error if file reading or YAML parsing fails
func loadCities(filePath string) (struct {
	Continents []struct {
		Name   string        `yaml:"name"`
		Cities []models.City `yaml:"cities"`
	} `yaml:"continents"`
}, error) {
	var cfg struct {
		Continents []struct {
			Name   string        `yaml:"name"`
			Cities []models.City `yaml:"cities"`
		} `yaml:"continents"`
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
