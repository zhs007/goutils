package goutils

import "time"

func GetCurTimestamp() int64 {
	return time.Now().Unix()
}

func FormatUTCDayTs(t time.Time) int64 {
	y, m, d := t.Date()

	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix()
}
