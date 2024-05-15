package utils

import (
	"math/rand"
)

// func init() {
// 	rand.Seed(time.Now().UnixNano())
// }

func SelectRandIdx(bound int, ratio float32) (pre []int, left []int) {
	ints := rand.Perm(bound)
	pivot := int(ratio * float32(bound))
	if pivot >= bound {
		pivot = bound - 1
	}
	if pivot == 0 {
		pivot = 1
	}
	pre = ints[:pivot]
	left = ints[pivot:]
	return
}

func GetRandomStr(l int) string {
	const str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, l)
	for i := 0; i < l; i++ {
		result[i] = str[rand.Intn(len(str))]
	}
	return string(result)
}

func GetRandomIntStr(l int) string {
	const str = "0123456789"
	result := make([]byte, l)
	for i := 0; i < l; i++ {
		result[i] = str[rand.Intn(len(str))]
	}
	return string(result)
}

func GetRandomInt() uint64 {
	var result uint64
	for i := 0; i < 16; i++ {
		d := rand.Intn(10)
		if i == 0 && d == 0 {
			d = 1
		}
		result = result*10 + uint64(d)
	}
	return result
}
