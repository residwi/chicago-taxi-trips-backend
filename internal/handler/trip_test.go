package handler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateDate(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		date, err := validateDate("2020-01-01")
		assert.NoError(t, err)
		assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), date)
	})

	t.Run("invalid date", func(t *testing.T) {
		_, err := validateDate("invalid-date")
		assert.Error(t, err)
	})
}
