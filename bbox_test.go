package h3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func assertBBox(t *testing.T, geofence *Geofence, expected *BBox, inside *GeoCoord, outside *GeoCoord) {
	var result BBox
	bboxFrom(geofence, &result)
	require.True(t, bboxEquals(&result, expected), "Got expected bbox")
	require.True(t, bboxContains(&result, inside), "Contains expected inside point")
	require.True(t, !bboxContains(&result, outside),
		"Does not contain expected outside point")
}

func Test_posLatPosLon(t *testing.T) {
	verts := []GeoCoord{{0.8, 0.3}, {0.7, 0.6}, {1.1, 0.7}, {1.0, 0.2}}
	geofence := Geofence{4, verts}
	expected := BBox{1.1, 0.7, 0.7, 0.2}
	inside := GeoCoord{0.9, 0.4}
	outside := GeoCoord{0.0, 0.0}
	assertBBox(t, &geofence, &expected, &inside, &outside)
}

func Test_negLatPosLon(t *testing.T) {
	verts := []GeoCoord{{-0.3, 0.6}, {-0.4, 0.9}, {-0.2, 0.8}, {-0.1, 0.6}}
	geofence := Geofence{4, verts}
	expected := BBox{-0.1, -0.4, 0.9, 0.6}
	inside := GeoCoord{-0.3, 0.8}
	outside := GeoCoord{0.0, 0.0}
	assertBBox(t, &geofence, &expected, &inside, &outside)
}

func Test_posLatNegLon(t *testing.T) {
	verts := []GeoCoord{{0.7, -1.4}, {0.8, -0.9}, {1.0, -0.8}, {1.1, -1.3}}
	geofence := Geofence{4, verts}
	expected := BBox{1.1, 0.7, -0.8, -1.4}
	inside := GeoCoord{0.9, -1.0}
	outside := GeoCoord{0.0, 0.0}
	assertBBox(t, &geofence, &expected, &inside, &outside)
}

func Test_negLatNegLon(t *testing.T) {
	verts := []GeoCoord{
		{-0.4, -1.4}, {-0.3, -1.1}, {-0.1, -1.2}, {-0.2, -1.4}}
	geofence := Geofence{4, verts}
	expected := BBox{-0.1, -0.4, -1.1, -1.4}
	inside := GeoCoord{-0.3, -1.2}
	outside := GeoCoord{0.0, 0.0}
	assertBBox(t, &geofence, &expected, &inside, &outside)
}

func Test_aroundZeroZero(t *testing.T) {
	verts := []GeoCoord{{0.4, -0.4}, {0.4, 0.4}, {-0.4, 0.4}, {-0.4, -0.4}}
	geofence := Geofence{4, verts}
	expected := BBox{0.4, -0.4, 0.4, -0.4}
	inside := GeoCoord{-0.1, -0.1}
	outside := GeoCoord{1.0, -1.0}
	assertBBox(t, &geofence, &expected, &inside, &outside)
}

func Test_transmeridian(t *testing.T) {
	verts := []GeoCoord{{0.4, M_PI - 0.1},
		{0.4, -M_PI + 0.1},
		{-0.4, -M_PI + 0.1},
		{-0.4, M_PI - 0.1}}
	geofence := Geofence{4, verts}
	expected := BBox{0.4, -0.4, -M_PI + 0.1, M_PI - 0.1}
	insideOnMeridian := GeoCoord{-0.1, M_PI}
	outside := GeoCoord{1.0, M_PI - 0.5}
	assertBBox(t, &geofence, &expected, &insideOnMeridian, &outside)
	westInside := GeoCoord{0.1, M_PI - 0.05}
	require.True(t, bboxContains(&expected, &westInside),
		"Contains expected west inside point")
	eastInside := GeoCoord{0.1, -M_PI + 0.05}
	require.True(t, bboxContains(&expected, &eastInside),
		"Contains expected east outside point")
	westOutside := GeoCoord{0.1, M_PI - 0.5}
	require.True(t, !bboxContains(&expected, &westOutside),
		"Does not contain expected west outside point")
	eastOutside := GeoCoord{0.1, -M_PI + 0.5}
	require.True(t, !bboxContains(&expected, &eastOutside),
		"Does not contain expected east outside point")
}

func Test_edgeOnNorthPole(t *testing.T) {
	verts := []GeoCoord{{M_PI_2 - 0.1, 0.1},
		{M_PI_2 - 0.1, 0.8},
		{M_PI_2, 0.8},
		{M_PI_2, 0.1}}
	geofence := Geofence{4, verts}
	expected := BBox{M_PI_2, M_PI_2 - 0.1, 0.8, 0.1}
	inside := GeoCoord{M_PI_2 - 0.01, 0.4}
	outside := GeoCoord{M_PI_2, 0.9}
	assertBBox(t, &geofence, &expected, &inside, &outside)
}

func Test_edgeOnSouthPole(t *testing.T) {
	verts := []GeoCoord{{-M_PI_2 + 0.1, 0.1},
		{-M_PI_2 + 0.1, 0.8},
		{-M_PI_2, 0.8},
		{-M_PI_2, 0.1}}
	geofence := Geofence{4, verts}
	expected := BBox{-M_PI_2 + 0.1, -M_PI_2, 0.8, 0.1}
	inside := GeoCoord{-M_PI_2 + 0.01, 0.4}
	outside := GeoCoord{-M_PI_2, 0.9}
	assertBBox(t, &geofence, &expected, &inside, &outside)
}

func Test_containsEdges(t *testing.T) {
	bbox := BBox{0.1, -0.1, 0.2, -0.2}
	points := []GeoCoord{
		{0.1, 0.2}, {0.1, 0.0}, {0.1, -0.2}, {0.0, 0.2},
		{-0.1, 0.2}, {-0.1, 0.0}, {-0.1, -0.2}, {0.0, -0.2},
	}
	numPoints := 8
	for i := 0; i < numPoints; i++ {
		require.True(t, bboxContains(&bbox, &points[i]), "Contains edge point")
	}
}

func Test_containsEdgesTransmeridian(t *testing.T) {
	bbox := BBox{0.1, -0.1, -M_PI + 0.2, M_PI - 0.2}
	points := []GeoCoord{
		{0.1, -M_PI + 0.2}, {0.1, M_PI}, {0.1, M_PI - 0.2},
		{0.0, -M_PI + 0.2}, {-0.1, -M_PI + 0.2}, {-0.1, M_PI},
		{-0.1, M_PI - 0.2}, {0.0, M_PI - 0.2},
	}
	numPoints := 8
	for i := 0; i < numPoints; i++ {
		require.True(t, bboxContains(&bbox, &points[i]), "Contains transmeridian edge point")
	}
}

func Test_bboxCenterBasicQuandrants(t *testing.T) {
	var center GeoCoord
	bbox1 := BBox{1.0, 0.8, 1.0, 0.8}
	expected1 := GeoCoord{0.9, 0.9}
	bboxCenter(&bbox1, &center)
	require.True(t, geoAlmostEqual(&center, &expected1), "pos/pos as expected")
	bbox2 := BBox{-0.8, -1.0, 1.0, 0.8}
	expected2 := GeoCoord{-0.9, 0.9}
	bboxCenter(&bbox2, &center)
	require.True(t, geoAlmostEqual(&center, &expected2), "neg/pos as expected")
	bbox3 := BBox{1.0, 0.8, -0.8, -1.0}
	expected3 := GeoCoord{0.9, -0.9}
	bboxCenter(&bbox3, &center)
	require.True(t, geoAlmostEqual(&center, &expected3), "pos/neg as expected")
	bbox4 := BBox{-0.8, -1.0, -0.8, -1.0}
	expected4 := GeoCoord{-0.9, -0.9}
	bboxCenter(&bbox4, &center)
	require.True(t, geoAlmostEqual(&center, &expected4), "neg/neg as expected")
	bbox5 := BBox{0.8, -0.8, 1.0, -1.0}
	expected5 := GeoCoord{0.0, 0.0}
	bboxCenter(&bbox5, &center)
	require.True(t, geoAlmostEqual(&center, &expected5),
		"around origin as expected")
}

func Test_bboxCenterTransmeridian(t *testing.T) {
	var center GeoCoord

	bbox1 := BBox{1.0, 0.8, -M_PI + 0.3, M_PI - 0.1}
	expected1 := GeoCoord{0.9, -M_PI + 0.1}
	bboxCenter(&bbox1, &center)
	require.True(t, geoAlmostEqual(&center, &expected1), "skew east as expected")
	bbox2 := BBox{1.0, 0.8, -M_PI + 0.1, M_PI - 0.3}
	expected2 := GeoCoord{0.9, M_PI - 0.1}
	bboxCenter(&bbox2, &center)
	require.True(t, geoAlmostEqual(&center, &expected2), "skew west as expected")
	bbox3 := BBox{1.0, 0.8, -M_PI + 0.1, M_PI - 0.1}
	expected3 := GeoCoord{0.9, M_PI}
	bboxCenter(&bbox3, &center)
	require.True(t, geoAlmostEqual(&center, &expected3),
		"on antimeridian as expected")
}

func Test_bboxIsTransmeridian(t *testing.T) {
	bboxNormal := BBox{1.0, 0.8, 1.0, 0.8}
	require.True(t, !bboxIsTransmeridian(&bboxNormal),
		"Normal bbox not transmeridian")
	bboxTransmeridian := BBox{1.0, 0.8, -M_PI + 0.3, M_PI - 0.1}
	require.True(t, bboxIsTransmeridian(&bboxTransmeridian), "Transmeridian bbox is transmeridian")
}

func Test_bboxEquals(t *testing.T) {
	bbox := BBox{1.0, 0.0, 1.0, 0.0}
	north := bbox
	north.north += 0.1
	south := bbox
	south.south += 0.1
	east := bbox
	east.east += 0.1
	west := bbox
	west.west += 0.1
	require.True(t, bboxEquals(&bbox, &bbox), "Equals self")
	require.True(t, !bboxEquals(&bbox, &north), "Not equals different north")
	require.True(t, !bboxEquals(&bbox, &south), "Not equals different south")
	require.True(t, !bboxEquals(&bbox, &east), "Not equals different east")
	require.True(t, !bboxEquals(&bbox, &west), "Not equals different west")
}
