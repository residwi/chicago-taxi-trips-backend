package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/residwi/chicago-taxi-trips-backend/service"

	"github.com/gin-gonic/gin"
)

type Trip struct {
	service service.ITrip
}

func NewTrip(router *gin.Engine, service service.ITrip) {
	handler := &Trip{service: service}

	router.GET("/total_trips", handler.TotalTrips)
}

func (t *Trip) TotalTrips(c *gin.Context) {
	if c.Query("start") == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "start date param is required"})
		return
	}

	if c.Query("end") == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "end date param is required"})
		return
	}

	startDate, err := validateDate(c.Query("start"))
	if err != nil {
		errorMessage := fmt.Sprintf("invalid start date: %s (use format: YYYY-MM-DD)", c.Query("start"))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorMessage})

		return
	}

	endDate, err := validateDate(c.Query("end"))
	if err != nil {
		errorMessage := fmt.Sprintf("invalid end date: %s (use format: YYYY-MM-DD)", c.Query("end"))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorMessage})

		return
	}

	if startDate.After(endDate) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "start date must be before end date"})

		return
	}

	totalTrips, err := t.service.GetTotalTripsByDateRange(startDate, endDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"data": totalTrips})
}

func validateDate(dateString string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, dateString)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
