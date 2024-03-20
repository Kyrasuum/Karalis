package lmath

import (
	"math"
)

// This file holds common functions which are used through-out the package

const (
	epsilon = 0.000000001
)

// Checks if two floats are equal. Doing a comparision using a small epsilon value
func closeEq(a, b, eps float64) bool {
	if a > b {
		return ((a - b) < eps)
	} else {
		return ((b - a) < eps)
	}
}

// Calculate the determinant of a 2x2 matrix.
// Values are givein in Row-Major order
func det2x2(x, y, z, w float64) float64 {
	return x*w - y*z
}

// Calculate the determinant of a 3x3 matrix.
// Values are given in Row-Major order
func det3x3(a1, a2, a3, b1, b2, b3, c1, c2, c3 float64) float64 {
	// a1 a2 a3
	// b1 b2 b3
	// c1 c2 c3
	return (a1*det2x2(b2, b3, c2, c3) -
		b1*det2x2(a2, a3, c2, c3) +
		c1*det2x2(a2, a3, b2, b3))
}

// Convert from a degree to a radian
func Radians(a float64) float64 {
	return a * math.Pi / 180.0
}

// Convert radian to degree
func Degrees(a float64) float64 {
	return a * 180.0 / math.Pi
}

// Clamp the a value between the lower and upper.
// Therefore the value returned is between the range [lower,upper] (inclusive)
func Clamp(a, lower, upper float64) float64 {
	if a < lower {
		return lower
	} else if a > upper {
		return upper
	}
	return a
}

// Linearly interpolates between the start and end values.
//
//	inc is specified between the range 0 -1
//	Lerp(0,2,0) ==> 0
//	Lerp(0,2,0.5) ==> 1
//	Lerp(0,2,1) ==> 2
func Lerp(start, end, inc float64) float64 {
	return (1-inc)*start + inc*end
}
