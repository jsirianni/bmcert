package timecalc

import (
    "strconv"
)

// BASE const is the unit base for rounding functions
const BASE int            = 10

// SECONDSDAY is the amount of seconds in a day
const SECONDSDAY int64   = 86400

// SECONDSMONTH is the amount of seconds in a month
const SECONDSMONTH int64 = 2.628e+6

// SecondsDay returns number of seconds in a day
// as an int64
func SecondsDay(duration int64) int64 {
	return duration * SECONDSDAY
}

// SecondsDayString returns number of seconds in a day
// as a string
func SecondsDayString(duration int64) string {
	return strconv.FormatInt(SecondsDay(duration), BASE)
}

// SecondsMonth returns number of seconds in a month
// as an int64
func SecondsMonth(duration int64) int64 {
	return duration * SECONDSMONTH
}

// SecondsMonthString returns number of seconds in a month
// as a string
func SecondsMonthString(duration int64) string {
	return strconv.FormatInt(SecondsMonth(duration), BASE)
}
