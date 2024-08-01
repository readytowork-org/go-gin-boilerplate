package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCalculateTimeInterval(t *testing.T) {
	interval := CalculateTimeInterval("11:00", "10:00", "15:04")

	assert.Equal(t, interval, -1*time.Hour)
}
