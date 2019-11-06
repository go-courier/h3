package h3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func createLinkedLoop(loop *LinkedGeoLoop, verts []GeoCoord, numVerts int) {
	*loop = LinkedGeoLoop{}

	for i := 0; i < numVerts; i++ {
		addLinkedCoord(loop, &verts[i])
	}
}

func Test_pointInsideLinkedGeoLoop(t *testing.T) {
	somewhere := GeoCoord{1, 2}
	inside := GeoCoord{0.659, -2.136}
	var loop LinkedGeoLoop
	createLinkedLoop(&loop, sfVerts, 6)

	var bbox BBox
	bboxFrom(&loop, &bbox)

	require.True(t, pointInside(&loop, &bbox, &inside), "contains exact4")
	require.True(t, !pointInside(&loop, &bbox, &somewhere), "contains somewhere else")

	destroyLinkedGeoLoop(&loop)
}

func Test_bboxFromGeofenceNoVertices(t *testing.T) {
	geofence := Geofence{}
	geofence.verts = nil
	geofence.numVerts = 0
	expected := BBox{0.0, 0.0, 0.0, 0.0}

	var result BBox
	bboxFrom(&geofence, &result)
	require.True(t, bboxEquals(&result, &expected), "Got expected bbox")
}

func Test_bboxFromLinkedGeoLoop(t *testing.T) {
	verts := []GeoCoord{{0.8, 0.3}, {0.7, 0.6}, {1.1, 0.7}, {1.0, 0.2}}
	var loop LinkedGeoLoop
	createLinkedLoop(&loop, verts, 4)

	expected := BBox{1.1, 0.7, 0.7, 0.2}
	var result BBox
	bboxFrom(&loop, &result)
	require.True(t, bboxEquals(&result, &expected), "Got expected bbox")
	destroyLinkedGeoLoop(&loop)
}

func Test_bboxFromLinkedGeoLoopNoVertices(t *testing.T) {
	loop := LinkedGeoLoop{}
	expected := BBox{0.0, 0.0, 0.0, 0.0}
	var result BBox
	bboxFrom(&loop, &result)

	require.True(t, bboxEquals(&result, &expected), "Got expected bbox")
	destroyLinkedGeoLoop(&loop)
}

func Test_isClockwiseLinkedGeoLoop(t *testing.T) {
	verts := []GeoCoord{{0.1, 0.1}, {0.2, 0.2}, {0.1, 0.2}}
	var loop LinkedGeoLoop
	createLinkedLoop(&loop, verts, 3)

	require.True(t, isClockwise(&loop), "Got true for clockwise loop")
	destroyLinkedGeoLoop(&loop)
}

func Test_isNotClockwiseLinkedGeoLoop(t *testing.T) {
	verts := []GeoCoord{{0, 0}, {0, 0.4}, {0.4, 0.4}, {0.4, 0}}
	var loop LinkedGeoLoop
	createLinkedLoop(&loop, verts, 4)
	require.True(t, !isClockwise(&loop), "Got false for counter-clockwise loop")
	destroyLinkedGeoLoop(&loop)
}

func Test_isClockwiseLinkedGeoLoopTransmeridian(t *testing.T) {
	verts := []GeoCoord{{0.4, M_PI - 0.1},
		{0.4, -M_PI + 0.1},
		{-0.4, -M_PI + 0.1},
		{-0.4, M_PI - 0.1}}
	var loop LinkedGeoLoop
	createLinkedLoop(&loop, verts, 4)
	require.True(t, isClockwise(&loop),
		"Got true for clockwise transmeridian loop")
	destroyLinkedGeoLoop(&loop)
}

func Test_isNotClockwiseLinkedGeoLoopTransmeridian(t *testing.T) {
	verts := []GeoCoord{{0.4, M_PI - 0.1},
		{-0.4, M_PI - 0.1},
		{-0.4, -M_PI + 0.1},
		{0.4, -M_PI + 0.1}}
	var loop LinkedGeoLoop
	createLinkedLoop(&loop, verts, 4)
	require.True(t, !isClockwise(&loop), "Got false for counter-clockwise transmeridian loop")
	destroyLinkedGeoLoop(&loop)
}

func Test_normalizeMultiPolygonSingle(t *testing.T) {
	verts := []GeoCoord{{0, 0}, {0, 1}, {1, 1}}
	outer := &LinkedGeoLoop{}
	createLinkedLoop(outer, verts, 3)
	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, outer)
	result := normalizeMultiPolygon(&polygon)
	require.True(t, result == NORMALIZATION_SUCCESS, "No error code returned")
	require.True(t, countLinkedPolygons(&polygon) == 1, "Polygon count correct")
	require.True(t, countLinkedLoops(&polygon) == 1, "Loop count correct")
	require.True(t, polygon.first == outer, "Got expected loop")
	destroyLinkedPolygon(&polygon)
}

func Test_normalizeMultiPolygonTwoOuterLoops(t *testing.T) {
	verts1 := []GeoCoord{{0, 0}, {0, 1}, {1, 1}}
	outer1 := &LinkedGeoLoop{}
	createLinkedLoop(outer1, verts1, 3)
	verts2 := []GeoCoord{{2, 2}, {2, 3}, {3, 3}}
	outer2 := &LinkedGeoLoop{}
	createLinkedLoop(outer2, verts2, 3)
	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, outer1)
	addLinkedLoop(&polygon, outer2)
	result := normalizeMultiPolygon(&polygon)
	require.True(t, result == NORMALIZATION_SUCCESS, "No error code returned")
	require.True(t, countLinkedPolygons(&polygon) == 2, "Polygon count correct")
	require.True(t, countLinkedLoops(&polygon) == 1,
		"Loop count on first polygon correct")
	require.True(t, countLinkedLoops(polygon.next) == 1,
		"Loop count on second polygon correct")
	destroyLinkedPolygon(&polygon)
}

func Test_normalizeMultiPolygonOneHole(t *testing.T) {
	verts := []GeoCoord{{0, 0}, {0, 3}, {3, 3}, {3, 0}}
	outer := &LinkedGeoLoop{}
	createLinkedLoop(outer, verts, 4)
	verts2 := []GeoCoord{{1, 1}, {2, 2}, {1, 2}}
	inner := &LinkedGeoLoop{}
	createLinkedLoop(inner, verts2, 3)
	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, inner)
	addLinkedLoop(&polygon, outer)
	result := normalizeMultiPolygon(&polygon)
	require.True(t, result == NORMALIZATION_SUCCESS, "No error code returned")
	require.True(t, countLinkedPolygons(&polygon) == 1, "Polygon count correct")
	require.True(t, countLinkedLoops(&polygon) == 2,
		"Loop count on first polygon correct")
	require.True(t, polygon.first == outer, "Got expected outer loop")
	require.True(t, polygon.first.next == inner, "Got expected inner loop")
	destroyLinkedPolygon(&polygon)
}

func Test_normalizeMultiPolygonTwoHoles(t *testing.T) {
	verts := []GeoCoord{{0, 0}, {0, 0.4}, {0.4, 0.4}, {0.4, 0}}
	outer := &LinkedGeoLoop{}
	require.True(t, outer != nil)
	createLinkedLoop(outer, verts, 4)
	verts2 := []GeoCoord{{0.1, 0.1}, {0.2, 0.2}, {0.1, 0.2}}
	inner1 := &LinkedGeoLoop{}
	require.True(t, inner1 != nil)
	createLinkedLoop(inner1, verts2, 3)
	verts3 := []GeoCoord{{0.2, 0.2}, {0.3, 0.3}, {0.2, 0.3}}
	inner2 := &LinkedGeoLoop{}

	createLinkedLoop(inner2, verts3, 3)
	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, inner2)
	addLinkedLoop(&polygon, outer)
	addLinkedLoop(&polygon, inner1)
	result := normalizeMultiPolygon(&polygon)
	require.True(t, result == NORMALIZATION_SUCCESS, "No error code returned")
	require.True(t, countLinkedPolygons(&polygon) == 1,
		"Polygon count correct for 2 holes")
	require.True(t, polygon.first == outer, "Got expected outer loop")
	require.True(t, countLinkedLoops(&polygon) == 3,
		"Loop count on first polygon correct")
	destroyLinkedPolygon(&polygon)
}

func Test_normalizeMultiPolygonTwoDonuts(t *testing.T) {
	verts := []GeoCoord{{0, 0}, {0, 3}, {3, 3}, {3, 0}}
	outer := &LinkedGeoLoop{}
	createLinkedLoop(outer, verts, 4)
	verts2 := []GeoCoord{{1, 1}, {2, 2}, {1, 2}}
	inner := &LinkedGeoLoop{}
	createLinkedLoop(inner, verts2, 3)
	verts3 := []GeoCoord{{0, 0}, {0, -3}, {-3, -3}, {-3, 0}}
	outer2 := &LinkedGeoLoop{}
	createLinkedLoop(outer2, verts3, 4)
	verts4 := []GeoCoord{{-1, -1}, {-2, -2}, {-1, -2}}
	inner2 := &LinkedGeoLoop{}
	createLinkedLoop(inner2, verts4, 3)
	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, inner2)
	addLinkedLoop(&polygon, inner)
	addLinkedLoop(&polygon, outer)
	addLinkedLoop(&polygon, outer2)
	result := normalizeMultiPolygon(&polygon)
	require.True(t, result == NORMALIZATION_SUCCESS, "No error code returned")
	require.True(t, countLinkedPolygons(&polygon) == 2, "Polygon count correct")
	require.True(t, countLinkedLoops(&polygon) == 2,
		"Loop count on first polygon correct")
	require.True(t, countLinkedCoords(polygon.first) == 4,
		"Got expected outer loop")
	require.True(t, countLinkedCoords(polygon.first.next) == 3,
		"Got expected inner loop")
	require.True(t, countLinkedLoops(polygon.next) == 2,
		"Loop count on second polygon correct")
	require.True(t, countLinkedCoords(polygon.next.first) == 4,
		"Got expected outer loop")
	require.True(t, countLinkedCoords(polygon.next.first.next) == 3,
		"Got expected inner loop")
	destroyLinkedPolygon(&polygon)
}

func Test_normalizeMultiPolygonNestedDonuts(t *testing.T) {
	verts := []GeoCoord{{0.2, 0.2}, {0.2, -0.2}, {-0.2, -0.2}, {-0.2, 0.2}}
	outer := &LinkedGeoLoop{}
	createLinkedLoop(outer, verts, 4)
	verts2 := []GeoCoord{
		{0.1, 0.1}, {-0.1, 0.1}, {-0.1, -0.1}, {0.1, -0.1}}
	inner := &LinkedGeoLoop{}
	createLinkedLoop(inner, verts2, 4)
	verts3 := []GeoCoord{
		{0.6, 0.6}, {0.6, -0.6}, {-0.6, -0.6}, {-0.6, 0.6}}
	outerBig := &LinkedGeoLoop{}
	createLinkedLoop(outerBig, verts3, 4)
	verts4 := []GeoCoord{
		{0.5, 0.5}, {-0.5, 0.5}, {-0.5, -0.5}, {0.5, -0.5}}
	innerBig := &LinkedGeoLoop{}
	createLinkedLoop(innerBig, verts4, 4)
	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, inner)
	addLinkedLoop(&polygon, outerBig)
	addLinkedLoop(&polygon, innerBig)
	addLinkedLoop(&polygon, outer)
	result := normalizeMultiPolygon(&polygon)
	require.True(t, result == NORMALIZATION_SUCCESS, "No error code returned")
	require.True(t, countLinkedPolygons(&polygon) == 2, "Polygon count correct")
	require.True(t, countLinkedLoops(&polygon) == 2,
		"Loop count on first polygon correct")
	require.True(t, polygon.first == outerBig, "Got expected outer loop")
	require.True(t, polygon.first.next == innerBig, "Got expected inner loop")
	require.True(t, countLinkedLoops(polygon.next) == 2,
		"Loop count on second polygon correct")
	require.True(t, polygon.next.first == outer, "Got expected outer loop")
	require.True(t, polygon.next.first.next == inner, "Got expected inner loop")
	destroyLinkedPolygon(&polygon)
}

func Test_normalizeMultiPolygonNoOuterLoops(t *testing.T) {
	verts1 := []GeoCoord{{0, 0}, {1, 1}, {0, 1}}
	outer1 := &LinkedGeoLoop{}
	createLinkedLoop(outer1, verts1, 3)
	verts2 := []GeoCoord{{2, 2}, {3, 3}, {2, 3}}
	outer2 := &LinkedGeoLoop{}
	createLinkedLoop(outer2, verts2, 3)
	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, outer1)
	addLinkedLoop(&polygon, outer2)
	result := normalizeMultiPolygon(&polygon)
	require.True(t, result == NORMALIZATION_ERR_UNASSIGNED_HOLES,
		"Expected error code returned")
	require.True(t, countLinkedPolygons(&polygon) == 1, "Polygon count correct")
	require.True(t, countLinkedLoops(&polygon) == 0,
		"Loop count as expected with invalid input")
	destroyLinkedPolygon(&polygon)
}

func Test_normalizeMultiPolygonAlreadyNormalized(t *testing.T) {
	verts1 := []GeoCoord{{0, 0}, {0, 1}, {1, 1}}
	outer1 := &LinkedGeoLoop{}
	createLinkedLoop(outer1, verts1, 3)
	verts2 := []GeoCoord{{2, 2}, {2, 3}, {3, 3}}
	outer2 := &LinkedGeoLoop{}
	createLinkedLoop(outer2, verts2, 3)
	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, outer1)
	next := addNewLinkedPolygon(&polygon)
	addLinkedLoop(next, outer2)

	// Should be a no-op
	result := normalizeMultiPolygon(&polygon)
	require.True(t, result == NORMALIZATION_ERR_MULTIPLE_POLYGONS,
		"Expected error code returned")
	require.True(t, countLinkedPolygons(&polygon) == 2, "Polygon count correct")
	require.True(t, countLinkedLoops(&polygon) == 1,
		"Loop count on first polygon correct")
	require.True(t, polygon.first == outer1, "Got expected outer loop")
	require.True(t, countLinkedLoops(polygon.next) == 1,
		"Loop count on second polygon correct")
	require.True(t, polygon.next.first == outer2, "Got expected outer loop")
	destroyLinkedPolygon(&polygon)
}

func Test_normalizeMultiPolygon_unassignedHole(t *testing.T) {
	verts := []GeoCoord{{0, 0}, {0, 1}, {1, 1}, {1, 0}}
	outer := &LinkedGeoLoop{}
	createLinkedLoop(outer, verts, 4)
	verts2 := []GeoCoord{{2, 2}, {3, 3}, {2, 3}}
	inner := &LinkedGeoLoop{}
	createLinkedLoop(inner, verts2, 3)
	polygon := LinkedGeoPolygon{}
	addLinkedLoop(&polygon, inner)
	addLinkedLoop(&polygon, outer)
	result := normalizeMultiPolygon(&polygon)
	require.True(t, result == NORMALIZATION_ERR_UNASSIGNED_HOLES, "Expected error code returned")
	destroyLinkedPolygon(&polygon)
}
