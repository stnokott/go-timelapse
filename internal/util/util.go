// Package util contains helper functions
package util

import (
	"fmt"
	"time"
)

// GetClosestPastTimeRange expects two strings in the form of "hh:mm" and returns the closest time range in the past relative to rel.
func GetClosestPastTimeRange(fromStr, toStr string, rel time.Time) (from time.Time, to time.Time, err error) {
	// Parse the "from" and "to" strings as time values on today's date.
	from, err = time.Parse("15:04", fromStr)
	if err != nil {
		err = fmt.Errorf("failed to parse from time: %s", err)
		return
	}
	to, err = time.Parse("15:04", toStr)
	if err != nil {
		err = fmt.Errorf("failed to parse to time: %s", err)
		return
	}
	relDay := rel.Truncate(24 * time.Hour)
	from = time.Date(relDay.Year(), relDay.Month(), relDay.Day(), from.Hour(), from.Minute(), 0, 0, rel.Location())
	to = time.Date(relDay.Year(), relDay.Month(), relDay.Day(), to.Hour(), to.Minute(), 0, 0, rel.Location())

	// If the "to" time is before the "from" time, it must be for the next day, so add a day to the "to" time.
	if to.Before(from) {
		to = to.Add(24 * time.Hour)
	}
	// if "to" is equal to "from", move "from" one day back.
	if to.Equal(from) {
		from = from.Add(-24 * time.Hour)
	}

	// If the time range is after the current time, subtract a day from the "from" time and a day from the "to" time.
	if to.After(rel) {
		from = from.Add(-24 * time.Hour)
		to = to.Add(-24 * time.Hour)
	}
	return
}
