package timecalc

import (
    "strconv"
)

// SecondsYear returns number of seconds in a month
func SecondsMonth(duration int) string {
	t := duration * 2.628e+6
	return strconv.Itoa(t)
}

// SecondsDay returns number of seconds in a day
func SecondsDay(duration int) string {
	t := duration * 86400
	return strconv.Itoa(t)
}
