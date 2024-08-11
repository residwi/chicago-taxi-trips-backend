package main

import (
	"github.com/residwi/chicago-taxi-trips-backend/internal/handler"
	"github.com/residwi/chicago-taxi-trips-backend/internal/repository/duckdb"
	"github.com/residwi/chicago-taxi-trips-backend/internal/router"
	"github.com/residwi/chicago-taxi-trips-backend/service"
)

func main() {
	server := router.SetupRouter()
	database := duckdb.CreateConnection()

	tripRepository := duckdb.NewTripRepository(database)
	tripService := service.NewTrip(tripRepository)
	handler.NewTrip(server, tripService)

	err := server.Run()
	if err != nil {
		panic(err)
	}
}
