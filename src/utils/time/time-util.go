package timeUtil

import (
	"time"
)

func GetUnixTime() int64 {
	return time.Now().Unix()
}

func SecFromMillis(time int64) int64 {
	return time / 1000
}

func MillisFromSec(time int64) int64 {
	return time * 1000
}

func DurationSeconds(value int) time.Duration {
	return time.Duration(value) * time.Second
}
