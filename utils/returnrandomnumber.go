package utils

import (
	"math/rand"
)

func ReturnRandomNumberBetween2And5(r *rand.Rand) float64 {
	// Generate a random float between 0 and 1, then scale to 2-5 range
	return r.Float64()*3 + 2 // rand.Float64() returns 0.0–1.0, multiply by 3 and add 2 to get 2.0–5.0
}
