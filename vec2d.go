package h3

import (
	"fmt"
	"math"
)

/**
 *  @brief 2D floating-point vector
 */
type Vec2d struct {
	x float64
	y float64
}

func (v Vec2d) String() string {
	return fmt.Sprintf("%f,%f", v.x, v.y)
}

/**
 * Calculates the magnitude of a 2D cartesian vector.
 * @param v The 2D cartesian vector.
 * @return The magnitude of the vector.
 */
func _v2dMag(v *Vec2d) float64 { return math.Sqrt(v.x*v.x + v.y*v.y) }

/**
 * Finds the intersection between two lines. Assumes that the lines intersect
 * and that the intersection is not at an endpoint of either line.
 * @param p0 The first endpoint of the first line.
 * @param p1 The second endpoint of the first line.
 * @param p2 The first endpoint of the second line.
 * @param p3 The second endpoint of the second line.
 * @param inter The intersection point.
 */
func _v2dIntersect(p0 *Vec2d, p1 *Vec2d, p2 *Vec2d, p3 *Vec2d, inter *Vec2d) {
	var s1, s2 Vec2d

	s1.x = p1.x - p0.x
	s1.y = p1.y - p0.y
	s2.x = p3.x - p2.x
	s2.y = p3.y - p2.y

	t := (s2.x*(p0.y-p2.y) - s2.y*(p0.x-p2.x)) / (-s2.x*s1.y + s1.x*s2.y)

	inter.x = p0.x + (t * s1.x)
	inter.y = p0.y + (t * s1.y)
}

/**
 * Whether two 2D vectors are equal. Does not consider possible false
 * negatives due to floating-point errors.
 * @param v1 First vector to compare
 * @param v2 Second vector to compare
 * @return Whether the vectors are equal
 */
func _v2dEquals(v1 *Vec2d, v2 *Vec2d) bool {
	return floatEqual(v1.x, v2.x) && floatEqual(v1.y, v2.y)
}

func floatEqual(a float64, b float64) bool {
	return nearlyEqual(a, b, 1e-9)
}

func nearlyEqual(a float64, b float64, epsilon float64) bool {
	if a == b {
		return true
	}

	diff := math.Abs(a - b)

	if diff < epsilon {
		return true
	}

	return diff <= epsilon*math.Min(math.Abs(a), math.Abs(b))
}

func normalizeDegree(v float64, min float64, max float64) float64 {
	d := max - min
	r := (v - min) / d
	return min + (r-float64(int64(r)))*d
}
