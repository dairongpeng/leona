package time

import (
	"fmt"
	"testing"
	"time"
)

func TestGetWeekDay(t *testing.T) {
	// tm := time.Now()
	tm := time.Date(2022, 01, 10, 0, 0, 0, 0, time.Local)
	fmt.Println(GetWeekDay(tm))
	fmt.Println(GetMonthDay(tm))
	fmt.Println(GetQuarterDay(tm))
}
