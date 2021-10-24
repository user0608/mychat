package utils

import (
	"math/rand"
	"time"
)

func GetRandon() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
func GetRandonInt() int {
	return GetRandon().Int()
}
