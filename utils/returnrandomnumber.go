package utils

import (
	"math/rand"
)

func ReturnRandomNumber(r *rand.Rand, min int, max int) float64 {
	return r.Float64()*(float64(max)-float64(min)+1) + float64(min)
}
