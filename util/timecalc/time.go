package timecalc

import (
    "strconv"
)

const BASE int            = 10
const SECONDS_DAY int64   = 86400
const SECONDS_MONTH int64 = 2.628e+6

// SecondsDay returns number of seconds in a day
// as an int64
func SecondsDay(duration int64) int64 {
	return duration * SECONDS_DAY
}

// SecondsDayString returns number of seconds in a day
// as a string
func SecondsDayString(duration int64) string {
	return strconv.FormatInt(SecondsDay(duration), BASE)
}

// SecondsYear returns number of seconds in a month
// as an int64
func SecondsMonth(duration int64) int64 {
	return duration * SECONDS_MONTH
}

// SecondsMonthString returns number of seconds in a month
// as a string
func SecondsMonthString(duration int64) string {
	return strconv.FormatInt(SecondsMonth(duration), BASE)
}
