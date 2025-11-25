package models

type WeatherResult struct {
	City                 string  `json:"city"`
	Temperature          float64 `json:"temperature"`
	Continent            string  `json:"continent" yaml:"continent"`
	Error                error   `json:"error"`
	FetchDurationSeconds float64 `json:"fetch_duration_seconds"`
}

type City struct {
	Name      string  `yaml:"name" json:"name"`
	Latitude  float64 `yaml:"latitude" json:"latitude"`
	Longitude float64 `yaml:"longitude" json:"longitude"`
}
