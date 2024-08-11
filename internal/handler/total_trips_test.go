package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/residwi/chicago-taxi-trips-backend/internal/router"
	"github.com/residwi/chicago-taxi-trips-backend/model"
	"github.com/residwi/chicago-taxi-trips-backend/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TotalTripsTestSuite struct {
	suite.Suite
}

func TestTotalTripsSuite(t *testing.T) {
	suite.Run(t, new(TotalTripsTestSuite))
}

func (suite *TotalTripsTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (suite *TotalTripsTestSuite) TestTotalTripsWithValidParams() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/total_trips?start=2020-01-01&end=2020-01-02", nil)
	router := router.SetupRouter()

	startDate, _ := time.Parse(time.DateOnly, "2020-01-01")
	endDate, _ := time.Parse(time.DateOnly, "2020-01-02")
	expectedTotalTrips := []model.TotalTrips{
		{Date: "2020-01-01", TotalTrips: 111},
		{Date: "2020-01-02", TotalTrips: 222},
	}
	tripServiceMock := mocks.NewMockITrip(suite.T())
	tripServiceMock.EXPECT().GetTotalTripsByDateRange(startDate, endDate).Return(expectedTotalTrips, nil)
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"data": expectedTotalTrips})
	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}

func (suite *TotalTripsTestSuite) TestTotalTripsWithEmptyStartDateParam() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/total_trips?end=2020-01-02", nil)
	router := router.SetupRouter()

	tripServiceMock := mocks.NewMockITrip(suite.T())

	NewTrip(router, tripServiceMock)
	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"error": "start date param is required"})
	assert.Equal(suite.T(), http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}

func (suite *TotalTripsTestSuite) TestTotalTripsWithEmptyEndDateParam() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/total_trips?start=2020-01-01", nil)
	router := router.SetupRouter()

	tripServiceMock := mocks.NewMockITrip(suite.T())
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"error": "end date param is required"})
	assert.Equal(suite.T(), http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}

func (suite *TotalTripsTestSuite) TestTotalTripsWithInvalidStartDateParam() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/total_trips?start=2020/01/01&end=2020-01-02", nil)
	router := router.SetupRouter()

	tripServiceMock := mocks.NewMockITrip(suite.T())
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"error": "invalid start date: 2020/01/01 (use format: YYYY-MM-DD)"})
	assert.Equal(suite.T(), http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}

func (suite *TotalTripsTestSuite) TestTotalTripsWithInvalidEndDateParams() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/total_trips?start=2020-01-01&end=2020/01/02", nil)
	router := router.SetupRouter()

	tripServiceMock := mocks.NewMockITrip(suite.T())
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"error": "invalid end date: 2020/01/02 (use format: YYYY-MM-DD)"})
	assert.Equal(suite.T(), http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}

func (suite *TotalTripsTestSuite) TestTotalTripsWithInvalidDateRangeParams() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/total_trips?start=2020-01-02&end=2020-01-01", nil)
	router := router.SetupRouter()

	tripServiceMock := mocks.NewMockITrip(suite.T())
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"error": "start date must be before end date"})
	assert.Equal(suite.T(), http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}

func (suite *TotalTripsTestSuite) TestTotalTripsWithInternalServerError() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/total_trips?start=2020-01-01&end=2020-01-02", nil)
	router := router.SetupRouter()

	startDate, _ := time.Parse(time.DateOnly, "2020-01-01")
	endDate, _ := time.Parse(time.DateOnly, "2020-01-02")
	tripServiceMock := mocks.NewMockITrip(suite.T())
	tripServiceMock.EXPECT().GetTotalTripsByDateRange(startDate, endDate).Return([]model.TotalTrips{}, assert.AnError)
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"error": "internal server error"})
	assert.Equal(suite.T(), http.StatusInternalServerError, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}
