package h3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var sfVerts = []GeoCoord{
	{0.659966917655, -2.1364398519396}, {0.6595011102219, -2.1359434279405},
	{0.6583348114025, -2.1354884206045}, {0.6581220034068, -2.1382437718946},
	{0.6594479998527, -2.1384597563896}, {0.6599990002976, -2.1376771158464},
}

func Test_pointInsideGeofence(t *testing.T) {
	geofence := Geofence{6, sfVerts}

	inside := GeoCoord{0.659, -2.136}
	somewhere := GeoCoord{1, 2}

	var bbox BBox
	bboxFrom(&geofence, &bbox)

	require.True(t, !pointInside(&geofence, &bbox, &sfVerts[0]), "contains exact")
	require.True(t, pointInside(&geofence, &bbox, &sfVerts[4]), "contains exact 4")
	require.True(t, pointInside(&geofence, &bbox, &inside), "contains point inside")
	require.True(t, !pointInside(&geofence, &bbox, &somewhere), "contains somewhere else")
}

func Test_pointInsideGeofenceTransmeridian(t *testing.T) {
	verts := []GeoCoord{{0.01, -M_PI + 0.01},
		{0.01, M_PI - 0.01},
		{-0.01, M_PI - 0.01},
		{-0.01, -M_PI + 0.01},
	}

	transMeridianGeofence := Geofence{4, verts}
	eastPoint := GeoCoord{0.001, -M_PI + 0.001}
	eastPointOutside := GeoCoord{0.001, -M_PI + 0.1}
	westPoint := GeoCoord{0.001, M_PI - 0.001}
	westPointOutside := GeoCoord{0.001, M_PI - 0.1}

	var bbox BBox
	bboxFrom(&transMeridianGeofence, &bbox)
	require.True(t, pointInside(&transMeridianGeofence, &bbox, &westPoint),
		"contains point to the west of the antimeridian")
	require.True(t, pointInside(&transMeridianGeofence, &bbox, &eastPoint),
		"contains point to the east of the antimeridian")
	require.True(t, !pointInside(&transMeridianGeofence, &bbox, &westPointOutside),
		"does not contain outside point to the west of the antimeridian")
	require.True(t, !pointInside(&transMeridianGeofence, &bbox, &eastPointOutside),
		"does not contain outside point to the east of the antimeridian")
}

func Test_bboxesFromGeoPolygon(t *testing.T) {
	t.Run("no hole", func(t *testing.T) {
		verts := []GeoCoord{{0.8, 0.3}, {0.7, 0.6}, {1.1, 0.7}, {1.0, 0.2}}
		geofence := Geofence{4, verts}
		polygon := GeoPolygon{geofence, 0, nil}
		expected := BBox{1.1, 0.7, 0.7, 0.2}
		result := make([]BBox, 1)

		bboxesFromGeoPolygon(&polygon, result)
		require.True(t, bboxEquals(&result[0], &expected), "Got expected bbox")
	})

	t.Run("with hole", func(t *testing.T) {
		verts := []GeoCoord{{0.8, 0.3}, {0.7, 0.6}, {1.1, 0.7}, {1.0, 0.2}}
		geofence := Geofence{numVerts: 4, verts: verts}

		// not a real hole, but doesn't matter for the test
		holeVerts := []GeoCoord{{0.9, 0.3}, {0.9, 0.5}, {1.0, 0.7}, {0.9, 0.3}}
		holeGeofence := Geofence{numVerts: 4, verts: holeVerts}
		polygon := GeoPolygon{geofence, 1, []Geofence{holeGeofence}}
		expected := BBox{1.1, 0.7, 0.7, 0.2}
		expectedHole := BBox{1.0, 0.9, 0.7, 0.3}

		result := make([]BBox, 2)
		bboxesFromGeoPolygon(&polygon, result)

		require.True(t, bboxEquals(&result[0], &expected), "Got expected bbox")
		require.True(t, bboxEquals(&result[1], &expectedHole), "Got expected hole bbox")
	})
}

func Test_isClockwiseGeofence(t *testing.T) {
	verts := []GeoCoord{{0, 0}, {0.1, 0.1}, {0, 0.1}}
	geofence := Geofence{3, verts}

	require.True(t, isClockwise(&geofence), "Got true for clockwise geofence")
}

func Test_isClockwiseGeofenceTransmeridian(t *testing.T) {
	verts := []GeoCoord{
		{0.4, M_PI - 0.1},
		{0.4, -M_PI + 0.1},
		{-0.4, -M_PI + 0.1},
		{-0.4, M_PI - 0.1},
	}
	
	geofence := Geofence{4, verts};
	require.True(t, isClockwise(&geofence), "Got true for clockwise geofence");
}
