# Chicago Taxi Trips API

## Overview

A well-architected REST API built with **Go** and **DuckDB** for analyzing 2020 Chicago taxi trip data.
This project demonstrates clean architecture principles, comprehensive testing strategies, and modern Go development practices.

### Project Highlights

- **Clean Architecture**: Layered design with dependency inversion
- **High Performance**: DuckDB analytical database with Parquet columnar storage
- **Comprehensive Testing**: 85%+ code coverage with unit and integration tests
- **Deployment Ready**: Docker support, graceful error handling, structured logging
- **Geospatial Analysis**: S2 geometry library for location-based fare heatmaps
- **API Documentation**: OpenAPI 3.0 specification included

### Tech Stack

- Go 1.22.6 ([https://github.com/golang/go](https://github.com/golang/go))
- Gin framework ([https://github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)) - HTTP routing and middleware
- DuckDB database ([https://github.com/marcboeker/go-duckdb](https://github.com/marcboeker/go-duckdb)) - Analytical SQL database for fast queries on Parquet data
- S2 Geometry ([https://github.com/golang/geo](https://github.com/golang/geo)) - Geospatial indexing for clustering and location-based analysis
- SQLMock ([https://github.com/DATA-DOG/go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)) - Database interaction mocking for tests
- Mockery ([https://github.com/vektra/mockery](https://github.com/vektra/mockery)) - Automated mock generation for interfaces

## Features

### Analytics Endpoints

| Endpoint                    | Description                   | Use Case                                     |
| --------------------------- | ----------------------------- | -------------------------------------------- |
| **`/total_trips`**          | Daily trip count aggregation  | Analyze trip volume trends over time         |
| **`/average_speed_24hrs`**  | Hourly average speed analysis | Identify rush hour patterns and traffic flow |
| **`/average_fare_heatmap`** | Geospatial fare distribution  | Visualize pricing patterns across Chicago    |

### Technical Features

- **RESTful API Design** with Gin web framework
- **Parquet Columnar Storage** for efficient data compression and retrieval
- **Analytical Database** using DuckDB for fast Parquet queries
- **Geospatial Analysis** with S2 cell-based aggregation
- **Clean Code** with comprehensive unit and integration tests
- **Docker Support** for easy deployment
- **Mock Generation** with Mockery for testability
- **OpenAPI Spec** for API documentation

## Quick Start

### Prerequisites

- **Go** 1.22+
- **Docker** & Docker Compose

### Installation

```bash
# Clone the repository
git clone https://github.com/residwi/chicago-taxi-trips-backend.git
cd chicago-taxi-trips-backend

# Setup (downloads dataset, runs tests, builds image)
./bin/setup

# Run the server
./bin/run
```

The API will be available at `http://localhost:8080`

### Alternative: Docker Compose

```bash
# Start with Docker Compose
make up

# Stop
make down
```

## API Documentation

### Endpoints Overview

#### Total Trips Aggregation

```bash
GET /total_trips?start=2020-01-01&end=2020-01-31
```

**Response:**

```json
{
  "data": [
    { "date": "2020-01-01", "total_trips": 12534 },
    { "date": "2020-01-02", "total_trips": 15678 }
  ]
}
```

#### Average Speed Analysis

```bash
GET /average_speed_24hrs?date=2020-01-15
```

**Response:**

```json
{
  "data": [
    { "hour": 0, "average_speed": 24.15 },
    { "hour": 1, "average_speed": 26.73 },
    { "hour": 2, "average_speed": null }
  ]
}
```

**Note:** Speeds are calculated in km/h with 2 decimal precision

#### Fare Heatmap

```bash
GET /average_fare_heatmap?date=2020-01-15
```

**Response:**

```json
{
  "data": [
    { "s2_cell_id": "9759b77c", "average_fare": 12.45 },
    { "s2_cell_id": "9759b784", "average_fare": 15.8 }
  ]
}
```

**Note:** Uses S2 geometry for geospatial clustering

### Full API Specification

Complete OpenAPI 3.0 specification is available at [`docs/api-doc.yaml`](./docs/api-doc.yaml)

## Architecture

### Project Structure

```
chicago-taxi-trips-backend/
├── internal/              # Internal packages (unexported)
│   ├── handler/          # HTTP request handlers
│   ├── repository/       # Database implementations
│   └── router/           # Route definitions
├── service/              # Business logic layer
├── repository/           # Repository interfaces
├── model/                # Domain models
├── test/                 # Integration tests
├── docs/                 # API documentation
```

### Key Design Patterns

- **Repository Pattern**: Abstracted data access
- **Dependency Injection**: Testable, loosely coupled components
- **Interface Segregation**: Clean boundaries between layers
- **Clean Architecture**: Business logic independent of frameworks

## Development

### Running Tests

```bash
# Run all tests with coverage
make test

# Run specific tests
go test ./internal/handler/...
go test ./service/...
go test ./test/...

# View coverage report
make coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run static analysis
make vet

# Generate mocks
make mockery
```

## Testing

### Test Coverage

- **Unit Tests**: Handlers, Services, Repositories
- **Integration Tests**: End-to-end API testing
- **Mock Generation**: Automated with Mockery
- **Race Detection**: Enabled by default
- **Coverage**: 85%+ across all packages

### Test Structure

```go
// Example: Handler test with mocks
func TestGetTotalTrips_ValidDateRange_ReturnsSuccess(t *testing.T) {
    // Arrange
    mockService := new(mocks.ITrip)
    handler := NewTrip(mockService)

    // Act
    handler.GetTotalTrips(ctx)

    // Assert
    assert.Equal(t, http.StatusOK, response.Code)
}
```

## Design Decisions

### Why Go?

- **Performance**: Fast execution and low memory footprint
- **Simplicity**: Clean syntax and easy to maintain
- **Concurrency**: Built-in goroutines for scalability
- **Standard Library**: Excellent HTTP and testing support

### Why DuckDB?

- **Parquet Support**: Native, efficient columnar storage
- **Analytical Queries**: Optimized for aggregations
- **Embedded**: No separate database server required
- **SQL Support**: Full SQL query capabilities

### Why No Caching?

- **Performance**: DuckDB is already very fast for analytical queries
- **Simplicity**: Avoiding premature optimization
- **YAGNI Principle**: Only add complexity when needed

### Architecture Motivation

- **Independence**: Business logic is decoupled from database and framework implementations
- **Testability**: Each layer can be tested independently using mocks
- **Maintainability**: Clear separation of concerns

## Project Showcase

### Key Implementations

#### Clean Architecture

```go
// Dependency Inversion: Handlers depend on interfaces
type Trip struct {
    service service.ITrip  // Interface, not concrete implementation
}

// Easy to test with mocks
func TestHandler(t *testing.T) {
    mockService := new(mocks.ITrip)
    handler := NewTrip(mockService)
    // Test without real database
}
```

#### Geospatial Analysis

```go
// S2 cell-based fare aggregation
cellID := s2.CellIDFromLatLng(s2.LatLngFromDegrees(lat, lng))
token := cellID.ToToken()  // Hex string for client
```
