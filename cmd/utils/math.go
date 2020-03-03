package utils

import "math"

// Min returns the minimum of two numbers
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Max returns the maximum of two numbers
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Mod returns the reminder of x/y
func Mod(x, y int) int {
	return int(math.Mod(float64(x), float64(y)))
}
