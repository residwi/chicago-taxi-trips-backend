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

type AverageSpeedTestSuite struct {
	suite.Suite
	tripService service.ITrip
}

func TestAverageSpeedSuite(t *testing.T) {
	suite.Run(t, new(AverageSpeedTestSuite))
}

func (suite *AverageSpeedTestSuite) SetupTest() {
	suite.T().Setenv("PARQUET_FILE_PATH", "../dataset/chicago_taxi_trips_2020.parquet")
	gin.SetMode(gin.TestMode)
	connDB := duckdb.CreateConnection()
	tripRepository := duckdb.NewTripRepository(connDB)
	suite.tripService = service.NewTrip(tripRepository)
}

func (suite *AverageSpeedTestSuite) TestAverageSpeed24hoursByDate() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/average_speed_24hrs?date=2020-01-01", nil)
	router := router.SetupRouter()
	handler.NewTrip(router, suite.tripService)

	router.ServeHTTP(responseRecorder, request)

	expectedAverageSpeed := []model.AverageSpeed{
		{AverageSpeed: 23.31},
	}

	expectedResponse, _ := json.Marshal(gin.H{"data": expectedAverageSpeed})
	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}

func (suite *AverageSpeedTestSuite) TestAverageSpeed24hoursByDateWhenNoData() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/average_speed_24hrs?date=2027-01-01", nil)
	router := router.SetupRouter()
	handler.NewTrip(router, suite.tripService)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"data": nil})
	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}
