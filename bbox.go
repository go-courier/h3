package h3

import "math"

/**
 *  @brief  Geographic bounding box with coordinates defined in radians
 */
type BBox struct {
	north float64 ///< north latitude
	south float64 ///< south latitude
	east  float64 ///< east longitude
	west  float64 ///< west longitude
}

/**
 * Whether the given bounding box crosses the antimeridian
 * @param  bbox Bounding box to inspect
 * @return      is transmeridian
 */
func bboxIsTransmeridian(bbox *BBox) bool {
	return bbox.east < bbox.west
}

/**
 * Get the center of a bounding box
 * @param bbox   Input bounding box
 * @param center Output center coordinate
 */
func bboxCenter(bbox *BBox, center *GeoCoord) {
	center.Lat = (bbox.north + bbox.south) / 2.0
	// If the bbox crosses the antimeridian, shift east 360 degrees
	east := bbox.east
	if bboxIsTransmeridian(bbox) {
		east = bbox.east + M_2PI
	}
	center.Lon = constrainLng((east + bbox.west) / 2.0)
}

/**
 * Whether the bounding box contains a given point
 * @param  bbox  Bounding box
 * @param  point Point to test
 * @return       Whether the point is contained
 */
func bboxContains(bbox *BBox, point *GeoCoord) bool {
	return point.Lat >= bbox.south && point.Lat <= bbox.north && func() bool {
		if bboxIsTransmeridian(bbox) {
			return point.Lon >= bbox.west || point.Lon <= bbox.east
		}
		return point.Lon >= bbox.west && point.Lon <= bbox.east
	}()
}

/**
 * Whether two bounding boxes are strictly equal
 * @param  b1 Bounding box 1
 * @param  b2 Bounding box 2
 * @return    Whether the boxes are equal
 */
func bboxEquals(b1 *BBox, b2 *BBox) bool {
	return b1.north == b2.north && b1.south == b2.south &&
		b1.east == b2.east && b1.west == b2.west
}

/**
 * _hexRadiusKm returns the radius of a given hexagon in Km
 *
 * @param h3Index the index of the hexagon
 * @return the radius of the hexagon in Km
 */
func _hexRadiusKm(h3Index H3Index) float64 {
	// There is probably a cheaper way to determine the radius of a
	// hexagon, but this way is conceptually simple
	var h3Center GeoCoord
	var h3Boundary GeoBoundary
	h3ToGeo(h3Index, &h3Center)
	h3ToGeoBoundary(h3Index, &h3Boundary)
	return _geoDistKm(&h3Center, &h3Boundary.Verts[0])
}

/**
* bboxHexEstimate returns an estimated number of hexagons that fit
*                 within the cartesian-projected bounding box
*
* @param bbox the bounding box to estimate the hexagon fill level
* @param res the resolution of the H3 hexagons to fill the bounding box
* @return the estimated number of hexagons to fill the bounding box
 */
func bboxHexEstimate(bbox *BBox, res int) int {
	// Get the area of the pentagon as the maximally-distorted area possible
	pentagons := make([]H3Index, 0)
	getPentagonIndexes(res, &pentagons)

	pentagonRadiusKm := _hexRadiusKm(pentagons[0])
	// Area of a regular hexagon is 3/2*sqrt(3) * r * r
	// The pentagon has the most distortion (smallest edges) and shares its
	// edges with hexagons, so the most-distorted hexagons have this area
	pentagonAreaKm2 := 2.59807621135 * pentagonRadiusKm * pentagonRadiusKm

	// Then get the area of the bounding box of the geofence in question
	var p1, p2 GeoCoord
	p1.Lat = bbox.north
	p1.Lon = bbox.east
	p2.Lat = bbox.south
	p2.Lon = bbox.east
	h := _geoDistKm(&p1, &p2)
	p2.Lat = bbox.north
	p2.Lon = bbox.west
	w := _geoDistKm(&p1, &p2)

	// Divide the two to get an estimate of the number of hexagons needed
	estimate := int(math.Ceil(w * h / pentagonAreaKm2))
	if estimate == 0 {
		estimate = 1
	}
	return estimate
}

/**
* lineHexEstimate returns an estimated number of hexagons that trace
*                 the cartesian-projected line
*
*  @param origin the origin coordinates
*  @param destination the destination coordinates
*  @param res the resolution of the H3 hexagons to trace the line
*  @return the estimated number of hexagons required to trace the line
 */
func lineHexEstimate(origin *GeoCoord, destination *GeoCoord, res int) int {
	// Get the area of the pentagon as the maximally-distorted area possible
	pentagons := make([]H3Index, 0)
	getPentagonIndexes(res, &pentagons)

	pentagonRadiusKm := _hexRadiusKm(pentagons[0])

	dist := _geoDistKm(origin, destination)

	estimate := int(math.Ceil(dist / (2 * pentagonRadiusKm)))
	if estimate == 0 {
		estimate = 1
	}
	return estimate
}
