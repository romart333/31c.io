package numberUtil

import (
	"math/rand"
	"strconv"
)

func GetRandomNumber(min, max int) int {
	if min >= max {
		return min // Avoid invalid range
	}
	return rand.Intn(max-min) + min
}

func IntToString(number int) string {
	return strconv.Itoa(number)
}
