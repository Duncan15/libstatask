package dates

import "time"

func GetStartTimeOfCurrentDay(t time.Time) time.Time {
	uix := GetStartUnixOfCurrentDay(t.Unix())
	return time.Unix(uix, 0)
}
func GetStartUnixOfCurrentDay(t int64) int64 {
	t -= t % (60 * 60 * 24)
	t -= 60 * 60 * 8
	return t
}
