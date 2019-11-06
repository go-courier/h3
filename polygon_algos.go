package h3

import "math"

/** Macro: Normalize longitude, dealing with transmeridian arcs */
func NORMALIZE_LON(lon float64, isTransmeridian bool) float64 {
	if isTransmeridian && lon < 0 {
		return lon + M_2PI
	}
	return lon
}

type GeoIterator interface {
	// zero
	IsZero() bool
	// true for continue
	NewIterate() func(vertexA *GeoCoord, vertexB *GeoCoord) bool
}

func pointInside(iterator GeoIterator, bbox *BBox, coord *GeoCoord) bool {
	// fail fast if we're outside the bounding box
	if !bboxContains(bbox, coord) {
		return false
	}

	isTransmeridian := bboxIsTransmeridian(bbox)
	contains := false

	lat := coord.Lat
	lng := NORMALIZE_LON(coord.Lon, isTransmeridian)

	var a GeoCoord
	var b GeoCoord

	iterate := iterator.NewIterate()

	for {
		if !iterate(&a, &b) {
			break
		}

		// Ray casting algo requires the second point to always be higher
		// than the first, so swap if needed
		if a.Lat > b.Lat {
			tmp := a
			a = b
			b = tmp
		}

		// If we're totally above or below the latitude ranges, the test
		// ray cannot intersect the line segment, so let's move on
		if lat < a.Lat || lat > b.Lat {
			continue
		}

		aLng := NORMALIZE_LON(a.Lon, isTransmeridian)
		bLng := NORMALIZE_LON(b.Lon, isTransmeridian)

		// Rays are cast in the longitudinal direction, in case a point
		// exactly matches, to decide tiebreakers, bias westerly
		if aLng == lng || bLng == lng {
			lng -= EPSILON
		}

		// For the latitude of the point, compute the longitude of the
		// point that lies on the line segment defined by a and b
		// This is done by computing the percent above a the Lat is,
		// and traversing the same percent in the longitudinal direction
		// of a to b
		ratio := (lat - a.Lat) / (b.Lat - a.Lat)
		testLng := NORMALIZE_LON(aLng+(bLng-aLng)*ratio, isTransmeridian)

		// Intersection of the ray
		if testLng > lng {
			contains = !contains
		}
	}

	return contains
}

/**
 * Create a bounding box from a simple polygon loop.
 * Known limitations:
 * - Does not support polygons with two adjacent points > 180 degrees of
 *   longitude apart. These will be interpreted as crossing the antimeridian.
 * - Does not currently support polygons containing a pole.
 * @param loop     Loop of coordinates
 * @param bbox     Output bbox
 */
func bboxFrom(iterator GeoIterator, bbox *BBox) {
	// Early exit if there are no vertices
	if iterator.IsZero() {
		*bbox = BBox{}
		return
	}

	bbox.south = math.MaxFloat64
	bbox.west = math.MaxFloat64
	bbox.north = -math.MaxFloat64
	bbox.east = -math.MaxFloat64
	minPosLon := math.MaxFloat64
	maxNegLon := -math.MaxFloat64

	isTransmeridian := false

	var lon, lat float64
	var coord, next GeoCoord

	iterate := iterator.NewIterate()

	for {
		if !iterate(&coord, &next) {
			break
		}

		lat = coord.Lat
		lon = coord.Lon

		if lat < bbox.south {
			bbox.south = lat
		}
		if lon < bbox.west {
			bbox.west = lon
		}
		if lat > bbox.north {
			bbox.north = lat
		}
		if lon > bbox.east {
			bbox.east = lon
		}

		// Save the min positive and max negative longitude for
		// use in the transmeridian case
		if lon > 0 && lon < minPosLon {
			minPosLon = lon
		}
		if lon < 0 && lon > maxNegLon {
			maxNegLon = lon
		}
		// check for arcs > 180 degrees longitude, flagging as transmeridian
		if math.Abs(lon-next.Lon) > M_PI {
			isTransmeridian = true
		}
	}

	// Swap east and west if transmeridian
	if isTransmeridian {
		bbox.east = maxNegLon
		bbox.west = minPosLon
	}
}

/**
 * Whether the winding order of a given loop is clockwise, with normalization
 * for loops crossing the antimeridian.
 * @param loop              The loop to check
 * @param isTransmeridian   Whether the loop crosses the antimeridian
 * @return                  Whether the loop is clockwise
 */
func isClockwiseNormalized(iterator GeoIterator, isTransmeridian bool) bool {
	sum := float64(0)
	var a, b GeoCoord

	iterate := iterator.NewIterate()

	for {
		if !iterate(&a, &b) {
			break
		}
		// If we identify a transmeridian arc (> 180 degrees longitude),
		// start over with the transmeridian flag set
		if !isTransmeridian && math.Abs(a.Lon-b.Lon) > M_PI {
			return isClockwiseNormalized(iterator, true)
		}

		sum += ((NORMALIZE_LON(b.Lon, isTransmeridian) - NORMALIZE_LON(a.Lon, isTransmeridian)) * (b.Lat + a.Lat))
	}

	return sum > 0
}

/**
 * Whether the winding order of a given loop is clockwise. In GeoJSON,
 * clockwise loops are always inner loops (holes).
 * @param loop  The loop to check
 * @return      Whether the loop is clockwise
 */
func isClockwise(iterator GeoIterator) bool {
	return isClockwiseNormalized(iterator, false)
}
