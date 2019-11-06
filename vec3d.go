package h3

import "math"

type Vec3d struct {
	x float64 /// < x component
	y float64 /// < y component
	z float64 /// < z component
}

/**
 * Square of a number
 *
 * @param x The input number.
 * @return The square of the input number.
 */
func _square(x float64) float64 { return x * x }

/**
 * Calculate the square of the distance between two 3D coordinates.
 *
 * @param v1 The first 3D coordinate.
 * @param v2 The second 3D coordinate.
 * @return The square of the distance between the given points.
 */
func _pointSquareDist(v1 *Vec3d, v2 *Vec3d) float64 {
	return _square(v1.x-v2.x) + _square(v1.y-v2.y) + _square(v1.z-v2.z)
}

/**
 * Calculate the 3D coordinate on unit sphere from the latitude and longitude.
 *
 * @param geo The latitude and longitude of the point.
 * @param v The 3D coordinate of the point.
 */
func _geoToVec3d(geo *GeoCoord, v *Vec3d) {
	r := math.Cos(geo.Lat)
	v.z = math.Sin(geo.Lat)
	v.x = math.Cos(geo.Lon) * r
	v.y = math.Sin(geo.Lon) * r
}
