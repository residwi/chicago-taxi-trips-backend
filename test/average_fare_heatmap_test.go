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

type AverageFareHeatmapTestSuite struct {
	suite.Suite
	tripService service.ITrip
}

func TestAverageFareHeatmapSuite(t *testing.T) {
	suite.Run(t, new(AverageFareHeatmapTestSuite))
}

func (suite *AverageFareHeatmapTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	connDB := duckdb.CreateConnection()
	tripRepository := duckdb.NewTripRepository(connDB)
	suite.tripService = service.NewTrip(tripRepository)

	connDB.Exec("CREATE TABLE test AS SELECT * FROM '../chicago_taxi_trips_2020.parquet' ORDER BY __index_level_0__ ASC LIMIT 10")
	suite.T().Setenv("PARQUET_FILE_PATH", "test")
}

func (suite *AverageFareHeatmapTestSuite) TestAverageFareHeatmapByDate() {
	responseRecorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/average_fare_heatmap?date=2020-02-23", nil)
	router := router.SetupRouter()
	handler.NewTrip(router, suite.tripService)

	router.ServeHTTP(responseRecorder, request)

	expectedAverageFareHeatmap := []model.AverageFareHeatmap{{S2ID: "880fcd625", AverageFare: 19.83}}
	expectedResponse, _ := json.Marshal(gin.H{"data": expectedAverageFareHeatmap})
	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code)
	assert.Equal(suite.T(), "application/json; charset=utf-8", responseRecorder.Header().Get("Content-Type"))
	assert.JSONEq(suite.T(), string(expectedResponse), responseRecorder.Body.String())
}
