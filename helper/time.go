package helper

import (
	"time"
)

func ParseTimeFromString(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

func GenerateRange(month, year int) (string, string) {
	timeNow := time.Now()

	if month == 0 && year == 0 {
		month = int(timeNow.Month())
		year = timeNow.Year()
	}

	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location())
	to := from.AddDate(0, 1, 0)

	return from.Format("2006-01-02"), to.Format("2006-01-02")
}
