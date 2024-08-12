package duckdb

import (
	"database/sql"
	"math"
	"os"
	"time"

	"github.com/residwi/chicago-taxi-trips-backend/model"

	log "github.com/sirupsen/logrus"
)

type TripRepository struct {
	conn *sql.DB
}

func NewTripRepository(conn *sql.DB) *TripRepository {
	return &TripRepository{conn: conn}
}

func (t *TripRepository) GetTotalTripsByDateRange(startDate time.Time, endDate time.Time) (totalTrips []model.TotalTrips, err error) {
	parquetFilepath := os.Getenv("PARQUET_FILE_PATH")

	query := `
		SELECT
			CAST(trip_start_timestamp AS DATE) AS pickup_date,
			COUNT(*) AS total_trips
		FROM '` + parquetFilepath + `'
		WHERE CAST(trip_start_timestamp AS DATE) BETWEEN $1 AND $2
		GROUP BY pickup_date
		ORDER BY pickup_date
    `

	rows, err := t.conn.Query(query, startDate.Format(time.DateOnly), endDate.Format(time.DateOnly))
	if err != nil {
		log.Error(err)

		return []model.TotalTrips{}, err
	}
	defer rows.Close()

	var pickupDate time.Time
	var totalTripsInt int
	for rows.Next() {
		err := rows.Scan(&pickupDate, &totalTripsInt)
		if err != nil {
			log.Error(err)

			continue
		}

		totalTrips = append(totalTrips, model.TotalTrips{
			Date:       pickupDate.Format(time.DateOnly),
			TotalTrips: totalTripsInt,
		})
	}

	return totalTrips, nil
}

func (t *TripRepository) GetAverageSpeedByDate(date time.Time) (averageSpeed []model.AverageSpeed, err error) {
	parquetFilepath := os.Getenv("PARQUET_FILE_PATH")

	query := `
		SELECT
			AVG((((trip_miles * 1.60934) / trip_seconds) * 3600)) AS average_speed
		FROM '` + parquetFilepath + `'
		WHERE CAST(trip_end_timestamp AS DATE) BETWEEN $1::date - INTERVAL '24 hour' AND $1
    `

	var avgSpeed float64
	err = t.conn.QueryRow(query, date.Format(time.DateOnly)).Scan(&avgSpeed)
	if err != nil {
		log.Error(err)

		return []model.AverageSpeed{}, err
	}

	avgSpeedOnlyTwoDecimal := math.Floor(avgSpeed*100) / 100
	averageSpeed = append(averageSpeed, model.AverageSpeed{
		AverageSpeed: avgSpeedOnlyTwoDecimal,
	})

	return averageSpeed, nil
}
