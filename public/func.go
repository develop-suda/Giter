package public

import (
	"strconv"
	"time"
)

// FormatDate returns the current date in YYYYMMDD format
func FormatDate() string {
	t := time.Now()

	y := strconv.Itoa(t.Year())
	m := strconv.Itoa(int(t.Month()))
	if len(m) == 1 {
		m = "0" + m
	}
	d := strconv.Itoa(t.Day())
	if len(d) == 1 {
		d = "0" + d
	}

	yyyymmdd := y + m + d
	return yyyymmdd
}
