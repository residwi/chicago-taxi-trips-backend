package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (t *Trip) AvarageFareHeatmap(c *gin.Context) {
	if c.Query("date") == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "date param is required"})
		return
	}

	date, err := validateDate(c.Query("date"))
	if err != nil {
		errorMessage := fmt.Sprintf("invalid date: %s (use format: YYYY-MM-DD)", c.Query("date"))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorMessage})

		return
	}

	averageFareHeatmaps, err := t.service.GetAverageFareHeatmapByDate(date)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"data": averageFareHeatmaps})
}
