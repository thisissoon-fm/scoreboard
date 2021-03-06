package db

import (
	"fmt"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
)

// Returns weekday (mon-fri) time range for the given date. If the given
// date is a Satuday or Sunday this date is adjusted to return the previous
// weeks time range
func weekdayTimeRange(date time.Time) (mon time.Time, fri time.Time) {
	wd := date.Weekday()
	// If we are not a mon-fri date, adjust date to the previous week
	if wd == time.Saturday {
		date = date.AddDate(0, 0, -1) // friday
		wd = date.Weekday()
	}
	if wd == time.Sunday {
		date = date.AddDate(0, 0, -2) // friday
		wd = date.Weekday()
	}
	mon = date
	for mon.Weekday() != time.Monday {
		mon = mon.AddDate(0, 0, -1)
	}
	fri = date
	for fri.Weekday() != time.Friday {
		fri = fri.AddDate(0, 0, 1)
	}
	return mon, fri.AddDate(0, 0, 1)
}

var scoresBetweenDateQry = `SELECT SUM("value") AS total
FROM "scores"
WHERE time >= '%s'
	AND time <= '%s'
GROUP BY "user";`

func ScoresByWeek(q Queryer, t time.Time) ([]influxdb.Result, error) {
	date := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	mon, fri := weekdayTimeRange(date)
	qry := fmt.Sprintf(scoresBetweenDateQry, mon.Format(time.RFC3339), fri.Format(time.RFC3339))
	return q.Query(qry)
}

func ScoresByYear(q Queryer, t time.Time) ([]influxdb.Result, error) {
	jan := time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
	dec := time.Date(t.Year(), time.December, 31, 24, 0, 0, 0, time.Local)
	qry := fmt.Sprintf(
		scoresBetweenDateQry,
		jan.Format(time.RFC3339),
		dec.Format(time.RFC3339),
	)
	return q.Query(qry)
}
