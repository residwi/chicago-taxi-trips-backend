package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/residwi/chicago-taxi-trips-backend/internal/handler"
	"github.com/residwi/chicago-taxi-trips-backend/internal/repository/duckdb"
	"github.com/residwi/chicago-taxi-trips-backend/internal/router"
	"github.com/residwi/chicago-taxi-trips-backend/model"
	"github.com/residwi/chicago-taxi-trips-backend/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TotalTripsTestSuite struct {
	suite.Suite
	tripService service.ITrip
}

func TestTotalTripsSuite(t *testing.T) {
	suite.Run(t, new(TotalTripsTestSuite))
}

func (suite *TotalTripsTestSuite) SetupTest() {
	suite.T().Setenv("PARQUET_FILE_PATH", "../chicago_taxi_trips_2020.parquet")
	gin.SetMode(gin.TestMode)
	connDB := duckdb.CreateConnection()
	tripRepository := duckdb.NewTripRepository(connDB)
	suite.tripService = service.NewTrip(tripRepository)
}

func (suite *TotalTripsTestSuite) TestTotalTripsPerDayBetweenDateRange() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/total_trips?start=2020-01-01&end=2020-01-02", nil)
	router := router.SetupRouter()
	handler.NewTrip(router, suite.tripService)

	router.ServeHTTP(responseRecorder, request)

	expectedTotalTrips := []model.TotalTrips{
		{Date: "2020-01-01", TotalTrips: 20798},
		{Date: "2020-01-02", TotalTrips: 26302},
	}

	expectedResponse, _ := json.Marshal(gin.H{"data": expectedTotalTrips})
	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}
