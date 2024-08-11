package duckdb

import (
	"database/sql"
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
