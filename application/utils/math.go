package utils

import "strconv"

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func F2S(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
