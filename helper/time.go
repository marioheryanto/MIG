package helper

import (
	"log"
	"time"
)

func TimeNowJakarta() time.Time {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println(err)
	}

	now := time.Now().In(loc)
	return now
}

func LoadLocationJakarta() *time.Location {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println(err)
	}

	return loc
}

func ParseTimeFromString(layout, value string) (time.Time, error) {
	loc := LoadLocationJakarta

	return time.ParseInLocation(layout, value, loc())
}

func GenerateRange(month, year int) (string, string) {
	loc := LoadLocationJakarta()
	timeNow := time.Now().In(loc)

	if month == 0 && year == 0 {
		month = int(timeNow.Month())
		year = timeNow.Year()
	}

	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	to := from.AddDate(0, 1, 0)

	return from.Format("2006-01-02"), to.Format("2006-01-02")
}
