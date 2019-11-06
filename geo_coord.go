package h3

import (
	"math"
)

/** epsilon of ~0.1mm in degrees */
const EPSILON_DEG = .000000001

/** epsilon of ~0.1mm in radians */
const EPSILON_RAD = (EPSILON_DEG * M_PI_180)

/**
 * Normalizes radians to a value between 0.0 and two PI.
 *
 * @param rads The input radians value.
 * @return The normalized radians value.
 */
func _posAngleRads(rads float64) float64 {
	var tmp float64

	if rads < 0.0 {
		tmp = rads + M_2PI
	} else {
		tmp = rads
	}

	if rads >= M_2PI {
		tmp -= M_2PI
	}
	return tmp
}

/**
 * Determines if the components of two spherical coordinates are within some
 * threshold distance of each other.
 *
 * @param p1 The first spherical coordinates.
 * @param p2 The second spherical coordinates.
 * @param threshold The threshold distance.
 * @return Whether or not the two coordinates are within the threshold distance
 *         of each other.
 */
func geoAlmostEqualThreshold(p1 *GeoCoord, p2 *GeoCoord, threshold float64) bool {
	return math.Abs(p1.Lat-p2.Lat) < threshold && math.Abs(p1.Lon-p2.Lon) < threshold
}

/**
 * Determines if the components of two spherical coordinates are within our
 * standard epsilon distance of each other.
 *
 * @param p1 The first spherical coordinates.
 * @param p2 The second spherical coordinates.
 * @return Whether or not the two coordinates are within the epsilon distance
 *         of each other.
 */
func geoAlmostEqual(p1 *GeoCoord, p2 *GeoCoord) bool {
	return geoAlmostEqualThreshold(p1, p2, EPSILON_RAD)
}

/**
 * Set the components of spherical coordinates in decimal degrees.
 *
 * @param p The spherical coodinates.
 * @param latDegs The desired latitidue in decimal degrees.
 * @param lonDegs The desired longitude in decimal degrees.
 */
func setGeoDegs(p *GeoCoord, latDegs float64, lonDegs float64) {
	_setGeoRads(p, degsToRads(latDegs), degsToRads(lonDegs))
}

/**
 * Set the components of spherical coordinates in radians.
 *
 * @param p The spherical coodinates.
 * @param latRads The desired latitidue in decimal radians.
 * @param lonRads The desired longitude in decimal radians.
 */
func _setGeoRads(p *GeoCoord, latRads float64, lonRads float64) {
	p.Lat = latRads
	p.Lon = lonRads
}

/**
 * Convert from decimal degrees to radians.
 *
 * @param degrees The decimal degrees.
 * @return The corresponding radians.
 */
func degsToRads(degrees float64) float64 { return degrees * M_PI_180 }

/**
 * Convert from radians to decimal degrees.
 *
 * @param radians The radians.
 * @return The corresponding decimal degrees.
 */
func radsToDegs(radians float64) float64 { return radians * M_180_PI }

/**
 * constrainLat makes sure latitudes are in the proper bounds
 *
 * @param Lat The original Lat value
 * @return The corrected Lat value
 */
func constrainLat(lat float64) float64 {
	for lat > M_PI_2 {
		lat = lat - M_PI
	}
	return lat
}

/**
 * constrainLng makes sure longitudes are in the proper bounds
 *
 * @param lng The origin lng value
 * @return The corrected lng value
 */
func constrainLng(lng float64) float64 {
	for lng > M_PI {
		lng = lng - (2 * M_PI)
	}
	for lng < -M_PI {
		lng = lng + (2 * M_PI)
	}
	return lng
}

/**
 * Find the great circle distance in radians between two spherical coordinates.
 *
 * @param p1 The first spherical coordinates.
 * @param p2 The second spherical coordinates.
 * @return The great circle distance in radians between p1 and p2.
 */
func _geoDistRads(p1 *GeoCoord, p2 *GeoCoord) float64 {
	// use spherical triangle with p1 at A, p2 at B, and north pole at C
	bigC := math.Abs(p2.Lon - p1.Lon)
	if bigC > M_PI {
		// note that in this case they can't both be negative
		lon1 := p1.Lon
		if lon1 < 0.0 {
			lon1 += 2.0 * M_PI
		}
		lon2 := p2.Lon

		if lon2 < 0.0 {
			lon2 += 2.0 * M_PI
		}

		bigC = math.Abs(lon2 - lon1)
	}

	b := M_PI_2 - p1.Lat
	a := M_PI_2 - p2.Lat

	// use law of cosines to find c
	cosc := math.Cos(a)*math.Cos(b) + math.Sin(a)*math.Sin(b)*math.Cos(bigC)
	if cosc > 1.0 {
		cosc = 1.0
	}
	if cosc < -1.0 {
		cosc = -1.0
	}

	return math.Acos(cosc)
}

/**
 * Find the great circle distance in kilometers between two spherical
 * coordinates.
 *
 * @param p1 The first spherical coordinates.
 * @param p2 The second spherical coordinates.
 * @return The distance in kilometers between p1 and p2.
 */
func _geoDistKm(p1 *GeoCoord, p2 *GeoCoord) float64 {
	return EARTH_RADIUS_KM * _geoDistRads(p1, p2)
}

/**
 * Determines the azimuth to p2 from p1 in radians.
 *
 * @param p1 The first spherical coordinates.
 * @param p2 The second spherical coordinates.
 * @return The azimuth in radians from p1 to p2.
 */
func _geoAzimuthRads(p1 *GeoCoord, p2 *GeoCoord) float64 {
	return math.Atan2(math.Cos(p2.Lat)*math.Sin(p2.Lon-p1.Lon),
		math.Cos(p1.Lat)*math.Sin(p2.Lat)-
			math.Sin(p1.Lat)*math.Cos(p2.Lat)*math.Cos(p2.Lon-p1.Lon))
}

/**
 * Computes the point on the sphere a specified azimuth and distance from
 * another point.
 *
 * @param p1 The first spherical coordinates.
 * @param az The desired azimuth from p1.
 * @param distance The desired distance from p1, must be non-negative.
 * @param p2 The spherical coordinates at the desired azimuth and distance from
 * p1.
 */
func _geoAzDistanceRads(p1 *GeoCoord, az float64, distance float64, p2 *GeoCoord) {
	if distance < EPSILON {
		*p2 = *p1
		return
	}

	var sinlat, sinlon, coslon float64

	az = _posAngleRads(az)

	// check for due north/south azimuth
	if az < EPSILON || math.Abs(az-M_PI) < EPSILON {
		if az < EPSILON {
			// due north
			p2.Lat = p1.Lat + distance
		} else {
			// due south
			p2.Lat = p1.Lat - distance
		}

		// north pole
		if math.Abs(p2.Lat-M_PI_2) < EPSILON {
			p2.Lat = M_PI_2
			p2.Lon = 0.0
		} else if math.Abs(p2.Lat+M_PI_2) < EPSILON {
			// south pole
			p2.Lat = -M_PI_2
			p2.Lon = 0.0
		} else {
			p2.Lon = constrainLng(p1.Lon)
		}
	} else {
		// not due north or south
		sinlat = math.Sin(p1.Lat)*math.Cos(distance) +
			math.Cos(p1.Lat)*math.Sin(distance)*math.Cos(az)
		if sinlat > 1.0 {
			sinlat = 1.0
		}
		if sinlat < -1.0 {
			sinlat = -1.0
		}
		p2.Lat = math.Asin(sinlat)
		if math.Abs(p2.Lat-M_PI_2) < EPSILON {
			// north pole
			p2.Lat = M_PI_2
			p2.Lon = 0.0
		} else if math.Abs(p2.Lat+M_PI_2) < EPSILON {
			// south pole
			p2.Lat = -M_PI_2
			p2.Lon = 0.0
		} else {
			sinlon = math.Sin(az) * math.Sin(distance) / math.Cos(p2.Lat)
			coslon = (math.Cos(distance) - math.Sin(p1.Lat)*math.Sin(p2.Lat)) /
				math.Cos(p1.Lat) / math.Cos(p2.Lat)
			if sinlon > 1.0 {
				sinlon = 1.0
			}
			if sinlon < -1.0 {
				sinlon = -1.0
			}
			if coslon > 1.0 {
				coslon = 1.0
			}
			if coslon < -1.0 {
				coslon = -1.0
			}
			p2.Lon = constrainLng(p1.Lon + math.Atan2(sinlon, coslon))
		}
	}
}

/*
 * The following functions provide meta information about the H3 hexagons at
 * each zoom level. Since there are only 16 total levels, these are current
 * handled with hardwired static values, but it may be worthwhile to put these
 * static values into another file that can be autogenerated by source code in
 * the future.
 */

var areaKM2s = []float64{
	4250546.848, 607220.9782, 86745.85403, 12392.26486,
	1770.323552, 252.9033645, 36.1290521, 5.1612932,
	0.7373276, 0.1053325, 0.0150475, 0.0021496,
	0.0003071, 0.0000439, 0.0000063, 0.0000009,
}

func hexAreaKm2(res int) float64 {
	return areaKM2s[res]
}

var areas = []float64{
	4.25055e+12, 6.07221e+11, 86745854035, 12392264862,
	1770323552, 252903364.5, 36129052.1, 5161293.2,
	737327.6, 105332.5, 15047.5, 2149.6,
	307.1, 43.9, 6.3, 0.9,
}

func hexAreaM2(res int) float64 {
	return areas[res]
}

var edgeLenOfKMs = []float64{
	1107.712591, 418.6760055, 158.2446558, 59.81085794,
	22.6063794, 8.544408276, 3.229482772, 1.220629759,
	0.461354684, 0.174375668, 0.065907807, 0.024910561,
	0.009415526, 0.003559893, 0.001348575, 0.000509713}

func edgeLengthKm(res int) float64 {
	return edgeLenOfKMs[res]
}

var edgeLenOfMs = []float64{
	1107712.591, 418676.0055, 158244.6558, 59810.85794,
	22606.3794, 8544.408276, 3229.482772, 1220.629759,
	461.3546837, 174.3756681, 65.90780749, 24.9105614,
	9.415526211, 3.559893033, 1.348574562, 0.509713273,
}

func edgeLengthM(res int) float64 {
	return edgeLenOfMs[res]
}

var nums = []int64{
	122,
	842,
	5882,
	41162,
	288122,
	2016842,
	14117882,
	98825162,
	691776122,
	4842432842,
	33897029882,
	237279209162,
	1660954464122,
	11626681248842,
	81386768741882,
	569707381193162,
}

/** @brief Number of unique valid H3Indexes at given resolution. */
func numHexagons(res int) int64 {
	return nums[res]
}
