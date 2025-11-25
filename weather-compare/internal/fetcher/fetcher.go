package fetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"weather-compare/internal/models"
)

func FetchWeatherData(city string, lat float64, lon float64) (models.WeatherResult, error) {
	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&hourly=temperature_2m&forecast_hours=1", lat, lon)
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return models.WeatherResult{City: city, Error: fmt.Errorf("HTTP error: %w", err), FetchDurationSeconds: time.Since(start).Seconds()}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.WeatherResult{City: city, Error: fmt.Errorf("API status %d", resp.StatusCode), FetchDurationSeconds: time.Since(start).Seconds()}, fmt.Errorf("API status %d", resp.StatusCode)
	}

	var respBody struct {
		Hourly struct {
			Temperature2M []float64 `json:"temperature_2m"`
		} `json:"hourly"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return models.WeatherResult{City: city, Error: fmt.Errorf("JSON error: %w", err), FetchDurationSeconds: time.Since(start).Seconds()}, err
	}
	if len(respBody.Hourly.Temperature2M) == 0 {
		return models.WeatherResult{City: city, Error: fmt.Errorf("no temperature data"), FetchDurationSeconds: time.Since(start).Seconds()}, fmt.Errorf("no temperature data")
	}

	return models.WeatherResult{City: city, Temperature: respBody.Hourly.Temperature2M[0], FetchDurationSeconds: time.Since(start).Seconds()}, nil
}
