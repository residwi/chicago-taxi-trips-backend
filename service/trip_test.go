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
