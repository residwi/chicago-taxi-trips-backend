package service

import (
	"math"
	"time"

	"github.com/golang/geo/s2"
	"github.com/residwi/chicago-taxi-trips-backend/model"
	"github.com/residwi/chicago-taxi-trips-backend/repository"

	log "github.com/sirupsen/logrus"
)

type ITrip interface {
	GetTotalTripsByDateRange(startDate time.Time, endDate time.Time) ([]model.TotalTrips, error)
	GetAverageSpeedByDate(date time.Time) ([]model.AverageSpeed, error)
	GetAverageFareHeatmapByDate(date time.Time) ([]model.AverageFareHeatmap, error)
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

		return nil, err
	}

	return totalTrips, nil
}

func (t *Trip) GetAverageSpeedByDate(date time.Time) ([]model.AverageSpeed, error) {
	averageSpeed, err := t.tripRepository.GetAverageSpeedByDate(date)
	if err != nil {
		log.Error(err)

		return nil, err
	}

	return averageSpeed, nil
}

func (t *Trip) GetAverageFareHeatmapByDate(date time.Time) ([]model.AverageFareHeatmap, error) {
	farePerPickupLocations, err := t.tripRepository.GetPickupLocationFareByDate(date)
	if err != nil {
		log.Error(err)

		return nil, err
	}

	fareSums := make(map[s2.CellID]float64)
	fareCounts := make(map[s2.CellID]int)
	s2Level := 16

	for _, farePerPickupLocation := range farePerPickupLocations {
		cellID := s2.CellFromLatLng(s2.LatLngFromDegrees(farePerPickupLocation.Latitude, farePerPickupLocation.Longitude)).
			ID().
			Parent(s2Level)

		fareSums[cellID] += farePerPickupLocation.Fare
		fareCounts[cellID]++
	}

	var averageFareHeatmaps []model.AverageFareHeatmap
	for cellID, sum := range fareSums {
		averageFare := sum / float64(fareCounts[cellID])
		averageFareHeatmaps = append(averageFareHeatmaps, model.AverageFareHeatmap{
			S2ID:        cellID.ToToken(),
			AverageFare: math.Floor(averageFare*100) / 100,
		})
	}

	return averageFareHeatmaps, nil
}
