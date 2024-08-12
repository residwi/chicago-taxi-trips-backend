package service

import (
	"time"

	"github.com/residwi/chicago-taxi-trips-backend/model"
	"github.com/residwi/chicago-taxi-trips-backend/repository"

	log "github.com/sirupsen/logrus"
)

type ITrip interface {
	GetTotalTripsByDateRange(startDate time.Time, endDate time.Time) ([]model.TotalTrips, error)
	GetAverageSpeedByDate(date time.Time) ([]model.AverageSpeed, error)
}

type Trip struct {
	tripRepository repository.ITrip
}

func NewTrip(tripRepository repository.ITrip) *Trip {
	return &Trip{tripRepository: tripRepository}
}

func (t *Trip) GetTotalTripsByDateRange(startDate time.Time, endDate time.Time) ([]model.TotalTrips, error) {
	totalTrips, err := t.tripRepository.GetTotalTripsByDateRange(startDate, endDate)
	if err != nil {
		log.Error(err)

		return []model.TotalTrips{}, err
	}

	return totalTrips, nil
}

func (t *Trip) GetAverageSpeedByDate(date time.Time) ([]model.AverageSpeed, error) {
	averageSpeed, err := t.tripRepository.GetAverageSpeedByDate(date)
	if err != nil {
		log.Error(err)

		return []model.AverageSpeed{}, err
	}

	return averageSpeed, nil
}
