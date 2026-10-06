package clock

import "time"

var WIB *time.Location

func init() {
	var err error
	WIB, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		WIB = time.FixedZone("WIB", 7*60*60)
	}
}

func NowMs() int64 {
	return time.Now().UnixMilli()
}

func DateKey(ms int64) string {
	return time.UnixMilli(ms).In(WIB).Format("2006-01-02")
}

func DueMs(date string, hour, min int) int64 {
	d, _ :=  time.ParseInLocation("2006-01-02", date, WIB)
	return time.Date(d.Year(), d.Month(), d.Day(), hour, min, 0, 0, WIB).UnixMilli()	
}

const (
	Second = int64(1_000)
	Minute = 60 * Second
)