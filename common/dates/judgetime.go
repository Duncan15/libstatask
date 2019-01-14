package dates

import "time"

//IsInCurrentHour check the specified time is in the current hour
func IsInCurrentHour(t int64) bool {
	ct := time.Now().Unix()
	chStart := ct - ct%3600
	if t-chStart > 3600 || t-chStart < 0 {
		return false
	}
	return true
}
