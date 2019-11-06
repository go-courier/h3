package h3

/**
 * Create a bounding box from a GeoPolygon
 * @param polygon Input GeoPolygon
 * @param bboxes  Output bboxes, one for the outer loop and one for each hole
 */
func bboxesFromGeoPolygon(polygon *GeoPolygon, bboxes []BBox) {
	bboxFrom(&polygon.geofence, &bboxes[0])
	for i := 0; i < polygon.numHoles; i++ {
		bboxFrom(&polygon.holes[i], &bboxes[i+1])
	}
}

/**
 * pointInsidePolygon takes a given GeoPolygon data structure and
 * checks if it contains a given geo coordinate.
 *
 * @param geoPolygon The geofence and holes defining the relevant area
 * @param bboxes     The bboxes for the main geofence and each of its holes
 * @param coord      The coordinate to check
 * @return           Whether the point is contained
 */
func pointInsidePolygon(geoPolygon *GeoPolygon, bboxes []BBox, coord *GeoCoord) bool {
	// Start with contains state of primary geofence
	contains := pointInside(&(geoPolygon.geofence), &bboxes[0], coord)

	// If the point is contained in the primary geofence, but there are holes in
	// the geofence iterate through all holes and return false if the point is
	// contained in any hole
	if contains && geoPolygon.numHoles > 0 {
		for i := 0; i < geoPolygon.numHoles; i++ {
			if pointInside(&(geoPolygon.holes[i]), &bboxes[i+1], coord) {
				return false
			}
		}
	}

	return contains
}
