package util

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomInt [min, max)
func RandomInt(min, max int) int {
	return seededRand.Intn(max-min) + min
}

// RandomFloat64 [0, 1)
func RandomFloat64() float64 {
	return seededRand.Float64()
}
