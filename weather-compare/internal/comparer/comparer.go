package comparer

import (
	"fmt"
	"math"
	"strings"
	"time"

	"weather-compare/internal/models"
)

func CompareWeatherData(weatherData []models.WeatherResult, contDurations map[string]time.Duration) string {
	if len(weatherData) == 0 {
		return "No weather data to compare."
	}

	groups := make(map[string][]models.WeatherResult)
	for _, r := range weatherData {
		groups[r.Continent] = append(groups[r.Continent], r)
	}

	var sb strings.Builder
	sb.WriteString("Weather comparison summary\n")
	sb.WriteString("========================================\n")

	for cont, items := range groups {
		sb.WriteString(fmt.Sprintf("\nContinent: %s\n", cont))
		sb.WriteString("----------------------------------------\n")

		var tempSum float64
		tempCount := 0
		minTemp := math.Inf(1)
		maxTemp := math.Inf(-1)

		var fetchSum float64
		fetchCount := 0
		minFetch := math.Inf(1)
		maxFetch := math.Inf(-1)

		for _, it := range items {
			if it.Error == nil {
				tempSum += it.Temperature
				tempCount++
				if it.Temperature < minTemp {
					minTemp = it.Temperature
				}
				if it.Temperature > maxTemp {
					maxTemp = it.Temperature
				}
			}
			fd := it.FetchDurationSeconds
			fetchSum += fd
			fetchCount++
			if fd < minFetch {
				minFetch = fd
			}
			if fd > maxFetch {
				maxFetch = fd
			}
		}

		if d, ok := contDurations[cont]; ok {
			sb.WriteString(fmt.Sprintf("Continent total wall time: %s\n", d.Truncate(time.Millisecond)))
		}

		if tempCount > 0 {
			sb.WriteString(fmt.Sprintf("Cities with valid temperature: %d, Avg: %.2f °C, Min: %.1f °C, Max: %.1f °C\n", tempCount, tempSum/float64(tempCount), minTemp, maxTemp))
		} else {
			sb.WriteString("No valid temperature data for this continent.\n")
		}

		if fetchCount > 0 {
			sb.WriteString(fmt.Sprintf("Fetch times (s): Avg: %.3f, Min: %.3f, Max: %.3f\n", fetchSum/float64(fetchCount), minFetch, maxFetch))
		}

		sb.WriteString("\nPer-city results:\n")
		for _, it := range items {
			if it.Error != nil {
				sb.WriteString(fmt.Sprintf("  %s: ERROR: %v (fetch %.3fs)\n", it.City, it.Error, it.FetchDurationSeconds))
			} else {
				sb.WriteString(fmt.Sprintf("  %s: %.1f °C (fetch %.3fs)\n", it.City, it.Temperature, it.FetchDurationSeconds))
			}
		}
	}

	return sb.String()
}
