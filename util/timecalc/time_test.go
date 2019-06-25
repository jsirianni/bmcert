package timecalc

import (
    "testing"
    "strconv"
)

func TestSecondsMonthString(t *testing.T) {
    s := SecondsMonthString(1)
    seconds, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        t.Errorf(err.Error())
        return
    }

    if seconds != 2.628e+6 {
        t.Errorf("Expected SecondsMonth(1) to return 2.628e+6")
    }
}

func TestSecondsDayString(t *testing.T) {
    s := SecondsDayString(1)
    seconds, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        t.Errorf(err.Error())
        return
    }

    if seconds != 86400 {
        t.Errorf("Expected SecondsMonth(1) to return 86400")
    }
}
