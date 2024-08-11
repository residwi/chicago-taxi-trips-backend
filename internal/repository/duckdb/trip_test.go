package duckdb

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/residwi/chicago-taxi-trips-backend/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TripRepositoryTestSuite struct {
	suite.Suite
	db     *sql.DB
	mockDB sqlmock.Sqlmock
}

func TestTripRepositorySuite(t *testing.T) {
	suite.Run(t, new(TripRepositoryTestSuite))
}

func (suite *TripRepositoryTestSuite) SetupSuite() {
	suite.T().Setenv("PARQUET_FILE_PATH", "stub-parquet-path.parquet")

	db, mockDB, err := sqlmock.New()
	if err != nil {
		suite.T().Fatal(err)
	}

	suite.db = db
	suite.mockDB = mockDB

	suite.T().Cleanup(func() {
		suite.db.Close()
	})
}

func (suite *TripRepositoryTestSuite) TestGetTotalTripsByDateRangeSuccess() {
	startDate, _ := time.Parse(time.DateOnly, "2020-01-01")
	endDate, _ := time.Parse(time.DateOnly, "2020-01-02")

	rows := sqlmock.NewRows([]string{"pickup_date", "total_trips"}).
		AddRow(startDate, 111).
		AddRow(endDate, 222)

	expectedQuery := `
		SELECT
			CAST(trip_start_timestamp AS DATE) AS pickup_date,
			COUNT(*) AS total_trips
		FROM 'stub-parquet-path.parquet'
		WHERE CAST(trip_start_timestamp AS DATE) BETWEEN $1 AND $2
		GROUP BY pickup_date
		ORDER BY pickup_date
    `

	suite.mockDB.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs("2020-01-01", "2020-01-02").
		WillReturnRows(rows).
		RowsWillBeClosed()

	repo := NewTripRepository(suite.db)
	result, err := repo.GetTotalTripsByDateRange(startDate, endDate)

	require.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)

	expectedResult := []model.TotalTrips{
		{Date: "2020-01-01", TotalTrips: 111},
		{Date: "2020-01-02", TotalTrips: 222},
	}
	assert.Equal(suite.T(), expectedResult, result)
}

func (suite *TripRepositoryTestSuite) TestGetTotalTripsByDateRangeErrorQuery() {
	startDate, _ := time.Parse(time.DateOnly, "2020-01-01")
	endDate, _ := time.Parse(time.DateOnly, "2020-01-02")

	expectedQuery := `
		SELECT
			CAST(trip_start_timestamp AS DATE) AS pickup_date,
			COUNT(*) AS total_trips
		FROM 'stub-parquet-path.parquet'
		WHERE CAST(trip_start_timestamp AS DATE) BETWEEN $1 AND $2
		GROUP BY pickup_date
		ORDER BY pickup_date
    `

	suite.mockDB.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs("2020-01-01", "2020-01-02").
		WillReturnError(assert.AnError).
		RowsWillBeClosed()

	repo := NewTripRepository(suite.db)
	result, err := repo.GetTotalTripsByDateRange(startDate, endDate)

	require.Error(suite.T(), err)
	assert.Len(suite.T(), result, 0)

	expectedResult := []model.TotalTrips{}
	assert.Equal(suite.T(), expectedResult, result)
}

func (suite *TripRepositoryTestSuite) TestGetTotalTripsByDateRangeErrorScan() {
	startDate, _ := time.Parse(time.DateOnly, "2020-01-01")
	endDate, _ := time.Parse(time.DateOnly, "2020-01-03")

	rows := sqlmock.NewRows([]string{"pickup_date", "total_trips"}).
		AddRow(startDate, 111).
		AddRow(nil, 222).
		AddRow(endDate, 333)

	expectedQuery := `
		SELECT
			CAST(trip_start_timestamp AS DATE) AS pickup_date,
			COUNT(*) AS total_trips
		FROM 'stub-parquet-path.parquet'
		WHERE CAST(trip_start_timestamp AS DATE) BETWEEN $1 AND $2
		GROUP BY pickup_date
		ORDER BY pickup_date
    `

	suite.mockDB.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs("2020-01-01", "2020-01-03").
		WillReturnRows(rows).
		RowsWillBeClosed()

	repo := NewTripRepository(suite.db)
	result, err := repo.GetTotalTripsByDateRange(startDate, endDate)

	require.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)

	expectedResult := []model.TotalTrips{
		{Date: "2020-01-01", TotalTrips: 111},
		{Date: "2020-01-03", TotalTrips: 333},
	}
	assert.Equal(suite.T(), expectedResult, result)
}
