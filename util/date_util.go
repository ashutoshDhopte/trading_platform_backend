package util

import "time"

func GetDateTimeString(time time.Time) string {
	return time.Format("01-02-2006 15:04:05")
}
