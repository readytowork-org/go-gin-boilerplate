package utils

import (
	"fmt"
	"time"
)

func ConvertStringToTime(timeString, layout string) (time.Time, error) {
	// Parse time strings into time.Time objects
	t1, err := time.Parse(layout, timeString)
	if err != nil {
		fmt.Println("Error parsing time:", err)
	}

	return t1, err
}

func CalculateTimeInterval(startTime, endTime, layout string) time.Duration {
	t1, err := ConvertStringToTime(startTime, layout)
	t2, err := ConvertStringToTime(endTime, layout)
	if err != nil {
		fmt.Println("Error converting time:", err)
	}

	// Calculate the duration between the two times
	duration := t2.Sub(t1)

	return duration
}

func IntervalGreaterThanCondition(reqStartDate time.Time, condition time.Duration) bool {
	currentDate := GetCurrentTimeInJP().Truncate(24 * time.Hour)
	requestedDate := reqStartDate.Truncate(24 * time.Hour)

	if currentDate.After(requestedDate){
		return false
	}
	if currentDate.Before(requestedDate) {
		return true
	}
	if currentDate.Equal(requestedDate) {
		currentTime := GetCurrentTimeInJP().Format("15:04:00")
		timeInterval := CalculateTimeInterval(currentTime, reqStartDate.Format("15:04:00"), "15:04:00")
		if timeInterval > condition {
			return true
		}
	}
	return false
}
