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

type AverageSpeedTestSuite struct {
	suite.Suite
}

func TestAverageSpeedSuite(t *testing.T) {
	suite.Run(t, new(AverageSpeedTestSuite))
}

func (suite *AverageSpeedTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (suite *AverageSpeedTestSuite) TestAverageSpeedWithValidParams() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/average_speed_24hrs?date=2020-01-01", nil)
	router := router.SetupRouter()

	date, _ := time.Parse(time.DateOnly, "2020-01-01")
	expecetAverageSpeed := []model.AverageSpeed{{AverageSpeed: 100.1}}
	tripServiceMock := mocks.NewMockITrip(suite.T())
	tripServiceMock.EXPECT().GetAverageSpeedByDate(date).Return(expecetAverageSpeed, nil)
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"data": expecetAverageSpeed})
	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}

func (suite *AverageSpeedTestSuite) TestAverageSpeedWithEmptyDateParam() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/average_speed_24hrs", nil)
	router := router.SetupRouter()

	tripServiceMock := mocks.NewMockITrip(suite.T())
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"error": "date param is required"})
	assert.Equal(suite.T(), http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}

func (suite *AverageSpeedTestSuite) TestAverageSpeedWithInvalidDateParam() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/average_speed_24hrs?date=stub-date", nil)
	router := router.SetupRouter()

	tripServiceMock := mocks.NewMockITrip(suite.T())
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"error": "invalid date: stub-date (use format: YYYY-MM-DD)"})
	assert.Equal(suite.T(), http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}

func (suite *AverageSpeedTestSuite) TestAverageSpeedWithInternalServerError() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/average_speed_24hrs?date=2020-01-01", nil)
	router := router.SetupRouter()

	date, _ := time.Parse(time.DateOnly, "2020-01-01")
	tripServiceMock := mocks.NewMockITrip(suite.T())
	tripServiceMock.EXPECT().GetAverageSpeedByDate(date).Return([]model.AverageSpeed{}, assert.AnError)
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"error": "internal server error"})
	assert.Equal(suite.T(), http.StatusInternalServerError, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}
