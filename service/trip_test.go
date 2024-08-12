package service

import (
	"testing"
	"time"

	"github.com/residwi/chicago-taxi-trips-backend/model"
	"github.com/residwi/chicago-taxi-trips-backend/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TripTestSuite struct {
	suite.Suite
}

func TestTripSuite(t *testing.T) {
	suite.Run(t, new(TripTestSuite))
}

func (suite *TripTestSuite) TestGetTotalTripsByDateRangeSuccess() {
	startDate, _ := time.Parse(time.DateOnly, "2020-01-01")
	endDate, _ := time.Parse(time.DateOnly, "2020-01-02")

	expectedTotalTrips := []model.TotalTrips{
		{Date: "2020-01-01", TotalTrips: 111},
		{Date: "2020-01-02", TotalTrips: 222},
	}
	tripRepositoryMock := mocks.NewMockITrip(suite.T())
	tripRepositoryMock.EXPECT().GetTotalTripsByDateRange(startDate, endDate).Return(expectedTotalTrips, nil)

	tripService := NewTrip(tripRepositoryMock)
	result, err := tripService.GetTotalTripsByDateRange(startDate, endDate)

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedTotalTrips, result)
}

func (suite *TripTestSuite) TestGetTotalTripsByDateRangeFailed() {
	startDate, _ := time.Parse(time.DateOnly, "2020-01-01")
	endDate, _ := time.Parse(time.DateOnly, "2020-01-02")

	tripRepositoryMock := mocks.NewMockITrip(suite.T())
	tripRepositoryMock.EXPECT().GetTotalTripsByDateRange(startDate, endDate).Return([]model.TotalTrips{}, assert.AnError)

	tripService := NewTrip(tripRepositoryMock)
	result, err := tripService.GetTotalTripsByDateRange(startDate, endDate)

	require.EqualError(suite.T(), err, assert.AnError.Error())
	assert.Equal(suite.T(), []model.TotalTrips{}, result)
}

func (suite *TripTestSuite) TestGetAverageSpeedByDateSuccess() {
	date, _ := time.Parse(time.DateOnly, "2020-01-01")

	expectedAverageSpeed := []model.AverageSpeed{{AverageSpeed: 100.1}}
	tripRepositoryMock := mocks.NewMockITrip(suite.T())
	tripRepositoryMock.EXPECT().GetAverageSpeedByDate(date).Return(expectedAverageSpeed, nil)

	tripService := NewTrip(tripRepositoryMock)
	result, err := tripService.GetAverageSpeedByDate(date)

	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedAverageSpeed, result)
}

func (suite *TripTestSuite) TestGetAverageSpeedByDateFailed() {
	date, _ := time.Parse(time.DateOnly, "2020-01-01")

	tripRepositoryMock := mocks.NewMockITrip(suite.T())
	tripRepositoryMock.EXPECT().GetAverageSpeedByDate(date).Return([]model.AverageSpeed{}, assert.AnError)

	tripService := NewTrip(tripRepositoryMock)
	result, err := tripService.GetAverageSpeedByDate(date)

	require.EqualError(suite.T(), err, assert.AnError.Error())
	assert.Equal(suite.T(), []model.AverageSpeed{}, result)
}

func (suite *TripTestSuite) TestGetAverageFareHeatmapByDateSuccess() {
	date, _ := time.Parse(time.DateOnly, "2020-01-01")

	farePerLocations := []model.FarePerLocation{
		{Latitude: 41.92276062, Longitude: -87.699155343, Fare: 45.0},
		{Latitude: 41.92276062, Longitude: -87.699155343, Fare: 5},
		{Latitude: 41.92276062, Longitude: -87.699155343, Fare: 9.5},
	}
	tripRepositoryMock := mocks.NewMockITrip(suite.T())
	tripRepositoryMock.EXPECT().GetPickupLocationFareByDate(date).Return(farePerLocations, nil)

	tripService := NewTrip(tripRepositoryMock)
	result, err := tripService.GetAverageFareHeatmapByDate(date)

	expectedAverageFareHeatmaps := []model.AverageFareHeatmap{{S2ID: "880fcd625", AverageFare: 19.83}}
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedAverageFareHeatmaps, result)
}

func (suite *TripTestSuite) TestGetAverageFareHeatmapByDateFailed() {
	date, _ := time.Parse(time.DateOnly, "2020-01-01")

	tripRepositoryMock := mocks.NewMockITrip(suite.T())
	tripRepositoryMock.EXPECT().GetPickupLocationFareByDate(date).Return([]model.FarePerLocation{}, assert.AnError)

	tripService := NewTrip(tripRepositoryMock)
	result, err := tripService.GetAverageFareHeatmapByDate(date)

	require.EqualError(suite.T(), err, assert.AnError.Error())
	assert.Equal(suite.T(), []model.AverageFareHeatmap{}, result)
}
