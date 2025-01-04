package public

import (
	"strconv"
	"time"
)

// FormatDate returns the current date in YYYYMMDD formatt
func FormatDate() string {
	t := time.Now()

	y := strconv.Itoa(t.Year())
	m := strconv.Itoa(int(t.Month()))
	d := strconv.Itoa(t.Day())

	yyyymmdd := y + m + d
	return yyyymmdd
}
