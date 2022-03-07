package util

import "math"

func SinScale(i int, n int) float32 {
	return float32(math.Sin((float64(i) - 0.5) / float64(n) * 180.0))
}
