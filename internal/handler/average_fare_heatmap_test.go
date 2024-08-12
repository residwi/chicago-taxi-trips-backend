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

type AverageFareHeatmapTestSuite struct {
	suite.Suite
}

func TestAverageFareHeatmapSuite(t *testing.T) {
	suite.Run(t, new(AverageFareHeatmapTestSuite))
}

func (suite *AverageFareHeatmapTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
}

func (suite *AverageFareHeatmapTestSuite) TestAverageFareHeatmapWithValidParams() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/average_fare_heatmap?date=2020-01-01", nil)
	router := router.SetupRouter()

	date, _ := time.Parse(time.DateOnly, "2020-01-01")
	expecetAverageFareHeatmaps := []model.AverageFareHeatmap{{S2ID: "stub-s2-id", AverageFare: 10}}
	tripServiceMock := mocks.NewMockITrip(suite.T())
	tripServiceMock.EXPECT().GetAverageFareHeatmapByDate(date).Return(expecetAverageFareHeatmaps, nil)
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"data": expecetAverageFareHeatmaps})
	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}

func (suite *AverageFareHeatmapTestSuite) TestAverageFareHeatmapWithEmptyDateParam() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/average_fare_heatmap", nil)
	router := router.SetupRouter()

	tripServiceMock := mocks.NewMockITrip(suite.T())
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"error": "date param is required"})
	assert.Equal(suite.T(), http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}

func (suite *AverageFareHeatmapTestSuite) TestAverageFareHeatmapWithInvalidDateParam() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/average_fare_heatmap?date=stub-date", nil)
	router := router.SetupRouter()

	tripServiceMock := mocks.NewMockITrip(suite.T())
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"error": "invalid date: stub-date (use format: YYYY-MM-DD)"})
	assert.Equal(suite.T(), http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}

func (suite *AverageFareHeatmapTestSuite) TestAverageFareHeatmapWithInternalServerError() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/average_fare_heatmap?date=2020-01-01", nil)
	router := router.SetupRouter()

	date, _ := time.Parse(time.DateOnly, "2020-01-01")
	tripServiceMock := mocks.NewMockITrip(suite.T())
	tripServiceMock.EXPECT().GetAverageFareHeatmapByDate(date).Return([]model.AverageFareHeatmap{}, assert.AnError)
	NewTrip(router, tripServiceMock)

	router.ServeHTTP(responseRecorder, request)

	expectedResponse, _ := json.Marshal(gin.H{"error": "internal server error"})
	assert.Equal(suite.T(), http.StatusInternalServerError, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}
