package repository

import (
	"time"

	"github.com/residwi/chicago-taxi-trips-backend/model"
)

type ITrip interface {
	GetTotalTripsByDateRange(startDate time.Time, endDate time.Time) ([]model.TotalTrips, error)
	GetAverageSpeedByDate(date time.Time) (averageSpeed []model.AverageSpeed, err error)
}
