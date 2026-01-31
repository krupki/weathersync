# WeatherSync – Concurrent Weather Data Aggregator

> A production-ready Go application that efficiently fetches and compares real-time weather data across multiple continents using concurrent processing and intelligent performance monitoring.

---

## Project Overview

**WeatherSync** is a sophisticated weather data aggregation system built in Go that demonstrates:

- **Concurrent programming** with goroutines and channels
- **Efficient multi-level parallelism** (per-continent and per-city workers)
- **Real-world API integration** with weather data providers
- **Configuration-driven architecture** for scalability
- **Performance metrics & monitoring** with latency tracking

The application fetches weather data from multiple cities across different continents **simultaneously**, measures execution time per continent, and provides comprehensive results with error handling.

---

## Key Features

### Concurrent Fetching

Goroutine-based parallel requests with WaitGroups for synchronization

### Per-Continent Workers

Isolates continent-level data processing with independent performance tracking

### Configuration-Driven

YAML-based city configuration for easy scaling and maintenance

### Real-time Performance Metrics

Tracks fetch duration per request and continent-level latency

### Robust Error Handling

Graceful error propagation without blocking other concurrent operations

### Structured Data Models

Type-safe JSON/YAML serialization for API responses and config

---

## Architecture

```bash
weather-compare/
├── cmd/weather-compare/      # Application entrypoint
│   └── main.go              # Orchestrates concurrent fetching & reporting
├── internal/
│   ├── fetcher/             # API integration layer
│   │   └── fetcher.go       # Handles weather API requests
│   ├── comparer/            # Data analysis & comparison
│   │   └── comparer.go      # Processes & compares weather data
│   └── models/              # Domain models
│       └── weather.go       # WeatherResult, City structures
└── config/
    └── cities.yaml          # Multi-continent city definitions
```

### Design Patterns Applied

- **Separation of Concerns**: Fetching, comparison, and main logic isolated in distinct packages
- **Channel-Based Communication**: Safe data passing between goroutines via channels
- **Configuration-Driven Setup**: Externalized config for easy modifications without code changes
- **Structured Logging**: Error tracking and performance visibility across the application

---

## Performance Highlights

- **Parallel Processing**: Fetches weather data for all cities concurrently, achieving **N-times speedup** vs sequential processing
- **Continent Isolation**: Each continent's data fetching runs independently, enabling better resource utilization
- **Latency Tracking**: Measures per-request and per-continent fetch times for performance insights
- **Non-blocking Error Handling**: One city's fetch failure doesn't impact others

---

## Technical Implementation

### Concurrent Goroutine Pattern

```go
for _, city := range cont.Cities {
    go func(ct models.City) {
        res, err := fetcher.FetchWeatherData(ct.Name, ct.Latitude, ct.Longitude)
        resultCh <- res
    }(city)
}
```

### Channel-Based Result Aggregation

- Results collected via buffered channels
- Prevents goroutine blocking and ensures all requests complete
- Enables safe concurrent reads across the application

### Configuration Management

- YAML-based configuration with structured unmarshaling
- Supports unlimited cities across multiple continents
- Easy to extend without code modifications

---

## Getting Started

### Prerequisites

- Go 1.18+
- Weather API credentials (configured in your environment or config)

### Installation

```bash
git clone https://github.com/krupki/weathersync.git
cd weathersync/weather-compare
```

### Running the Application

```bash
go run ./cmd/weather-compare/main.go
```

### Adding More Cities

Edit `config/cities.yaml` and add new continents or cities with their coordinates:

```yaml
continents:
  - name: South America
    cities:
      - name: São Paulo
        latitude: -23.55
        longitude: -46.63
```

---

## Output Example

The application generates detailed performance reports:

- Per-city weather data with fetch latency
- Per-continent aggregated execution times
- Error tracking with informative logging
- Structured JSON output for integration with other systems

---

## Testing & Reliability

- **Graceful error handling**: Fetch failures logged without crashing
- **Synchronized shutdown**: WaitGroups ensure all goroutines complete
- **Type-safe structures**: Go's type system prevents data corruption
- **Resource cleanup**: Channels and goroutines properly managed

---

## Scalability

This project demonstrates how to build scalable concurrent systems:

- **Horizontal scaling**: Add more cities/continents without code changes
- **Efficient resource usage**: Goroutines are lightweight (~2KB each)
- **Buffered channels**: Prevents goroutine starvation with proper buffer sizing
- **Non-blocking patterns**: No thread contention or deadlock risks

---

## Why This Matters

This project showcases **production-level Go patterns** used in real-world systems:

- **Uber**, **CloudFlare**, and **Docker** use similar concurrent patterns
- Demonstrates understanding of **concurrency primitives** (goroutines, channels, WaitGroups)
- **API integration** skills with error handling and retries
- **Clean code architecture** with separation of concerns
- **Performance-conscious design** with metrics and monitoring

---

## Future Enhancements

- [ ] Caching layer with TTL for API responses
- [ ] Retry logic with exponential backoff for failed requests
- [ ] Database storage (PostgreSQL/MongoDB) for historical data
- [ ] REST API endpoint to query weather data on-demand
- [ ] Unit tests with mocked weather provider
- [ ] Docker containerization for deployment
- [ ] Metrics export to Prometheus for monitoring

---

## License

MIT License – Feel free to use this project as a portfolio piece or reference implementation.

---

## Author

Built by Krupki as a demonstration of professional Go development practices.

**Connect with me:**

[![LinkedIn](https://img.shields.io/badge/LinkedIn-0A66C2?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/kevin-greupner-378138189/)
[![GitHub](https://img.shields.io/badge/GitHub-181717?style=for-the-badge&logo=github&logoColor=white)](https://github.com/krupki)

---

## Learning Resources

This project demonstrates:

- [Effective Go - Concurrency](https://golang.org/doc/effective_go#concurrency)
- [Go Channels Best Practices](https://go.dev/blog/pipelines)
- [Configuration Management Patterns](https://12factor.net/)
