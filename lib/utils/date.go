package utils

import (
	"fmt"
	"time"
)

// ConvertStringToDate Converts the string to date
func ConvertStringToDate(date string) time.Time {
	_time, _ := time.Parse("2006-01-02", date)
	return _time
}

// ConvertDateStringToYearMonth Converts the string to year month
func ConvertDateStringToYearMonth(date string) (int, int) {
	_time, _ := time.Parse("2006-01", date)
	return _time.Year(), int(_time.Month())
}

// ConvertRFCStringToDate Convert RFC date string to date
func ConvertRFCStringToDate(date string) time.Time {
	_time, err := time.Parse(time.RFC3339, date)
	fmt.Println(err)
	return _time
}

func GetJPLocation() *time.Location {
	loc, _ := time.LoadLocation("Japan")
	return loc
}

func GetCurrentTimeInJP() time.Time {
	return time.Now().In(GetJPLocation())
}

func GetCurrentTimeInJPString() string {
	return GetCurrentTimeInJP().Format("2006-01-02 15:04:05")
}

func MonthDate(dateString string) (time.Time, time.Time) {
	date := ConvertStringToDate(dateString)

	firstDay := time.Date(date.Year(), date.Month(), 1,
		0, 0, 0, 0, nil)
	return firstDay, firstDay.AddDate(0, 1, 0)
}
