package handler

import (
	"time"

	"github.com/residwi/chicago-taxi-trips-backend/service"

	"github.com/gin-gonic/gin"
)

type Trip struct {
	service service.ITrip
}

func NewTrip(router *gin.Engine, service service.ITrip) {
	handler := &Trip{service: service}

	router.GET("/total_trips", handler.TotalTrips)
	router.GET("/average_speed_24hrs", handler.AverageSpeed)
	router.GET("/average_fare_heatmap", handler.AvarageFareHeatmap)
}

func validateDate(dateString string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, dateString)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
