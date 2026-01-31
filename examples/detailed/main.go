package main

import (
	"context"
	"fmt"

	"github.com/krupki/weathersync"
)

func main() {
	client := weathersync.New()

	locations := []weathersync.Location{
		{
			Name:      "Berlin",
			Latitude:  52.52,
			Longitude: 13.41,
		},
		{
			Name:      "London",
			Latitude:  51.51,
			Longitude: -0.13,
		},
		{
			Name:      "Sydney",
			Latitude:  -33.87,
			Longitude: 151.21,
		},
	}

	fmt.Println("🌤️  Detailed Weather Report")
	fmt.Println("============================\n")

	results := client.FetchMultiple(context.Background(), locations)

	for _, weather := range results {
		if weather.Error != nil {
			fmt.Printf("❌ %s: %v\n", weather.Location.Name, weather.Error)
			continue
		}

		fmt.Printf("📍 %s\n", weather.Location.Name)
		fmt.Printf("   Temperature:        %.1f°C\n", weather.Temperature)
		fmt.Printf("   Feels Like:          %.1f°C\n", weather.ApparentTemperature)
		fmt.Printf("   Humidity:            %.0f%%\n", weather.Humidity)
		fmt.Printf("   Wind Speed:          %.1f km/h\n", weather.WindSpeed)
		fmt.Printf("   Wind Direction:      %.0f°\n", weather.WindDirection)
		fmt.Printf("   Wind Gusts:          %.1f km/h\n", weather.WindGusts)
		fmt.Printf("   Precipitation:       %.1f mm\n", weather.Precipitation)
		fmt.Printf("   Cloud Cover:         %.0f%%\n", weather.CloudCover)
		fmt.Printf("   Visibility:          %.0f m\n", weather.Visibility)
		fmt.Printf("   Pressure:            %.0f hPa\n", weather.Pressure)
		fmt.Printf("   Weather Code:        %d\n", weather.WeatherCode)
		fmt.Printf("   Fetched:             %.3fs\n", weather.FetchDuration.Seconds())
		fmt.Println()
	}
}
