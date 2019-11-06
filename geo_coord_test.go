package h3

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_radsToDegs(t *testing.T) {
	originalRads := 0.1
	degs := radsToDegs(originalRads)
	rads := degsToRads(degs)
	require.True(t, math.Abs(rads-originalRads) < EPSILON_RAD, "radsToDegs/degsToRads invertible")
}

func Test__geoDistRads(t *testing.T) {
	var p1 GeoCoord
	setGeoDegs(&p1, 10, 10)
	var p2 GeoCoord
	setGeoDegs(&p2, 0, 10)

	// TODO: Epsilon is relatively large
	require.True(t, _geoDistRads(&p1, &p1) < EPSILON_RAD*1000, "0 distance as expected")
	require.True(t, math.Abs(_geoDistRads(&p1, &p2)-degsToRads(10)) < EPSILON_RAD*1000, "distance along longitude as expected")
}

func Test_constrainLatLng(t *testing.T) {
	require.True(t, constrainLat(0) == 0, "Lat 0")
	require.True(t, constrainLat(1) == 1, "Lat 1")
	require.True(t, constrainLat(M_PI_2) == M_PI_2, "Lat pi/2")
	require.True(t, constrainLat(M_PI) == 0, "Lat pi")
	require.True(t, constrainLat(M_PI+1) == 1, "Lat pi+1")
	require.True(t, constrainLat(2*M_PI+1) == 1, "Lat 2pi+1")
	require.True(t, constrainLng(0) == 0, "lng 0")
	require.True(t, constrainLng(1) == 1, "lng 1")
	require.True(t, constrainLng(M_PI) == M_PI, "lng pi")
	require.True(t, constrainLng(2*M_PI) == 0, "lng 2pi")
	require.True(t, constrainLng(3*M_PI) == M_PI, "lng 2pi")
	require.True(t, constrainLng(4*M_PI) == 0, "lng 4pi")
}

func Test__geoAzDistanceRads_noop(t *testing.T) {
	start := GeoCoord{15, 10}
	var out GeoCoord
	expected := GeoCoord{15, 10}
	_geoAzDistanceRads(&start, 0, 0, &out)
	require.True(t, geoAlmostEqual(&expected, &out),
		"0 distance produces same point")
}

func Test__geoAzDistanceRads_dueNorthSouth(t *testing.T) {
	var start GeoCoord
	var out GeoCoord
	var expected GeoCoord

	// Due north to north pole
	setGeoDegs(&start, 45, 1)
	setGeoDegs(&expected, 90, 0)
	_geoAzDistanceRads(&start, 0, degsToRads(45), &out)
	require.True(t, geoAlmostEqual(&expected, &out),
		"due north to north pole produces north pole")

	// Due north to south pole, which doesn't get wrapped correctly
	setGeoDegs(&start, 45, 1)
	setGeoDegs(&expected, 270, 1)
	_geoAzDistanceRads(&start, 0, degsToRads(45+180), &out)
	require.True(t, geoAlmostEqual(&expected, &out),
		"due north to south pole produces south pole")

	// Due south to south pole
	setGeoDegs(&start, -45, 2)
	setGeoDegs(&expected, -90, 0)
	_geoAzDistanceRads(&start, degsToRads(180), degsToRads(45), &out)
	require.True(t, geoAlmostEqual(&expected, &out),
		"due south to south pole produces south pole")

	// Due north to non-pole
	setGeoDegs(&start, -45, 10)
	setGeoDegs(&expected, -10, 10)
	_geoAzDistanceRads(&start, 0, degsToRads(35), &out)
	require.True(t, geoAlmostEqual(&expected, &out),
		"due north produces expected result")
}

func Test__geoAzDistanceRads_poleToPole(t *testing.T) {
	var start GeoCoord
	var out GeoCoord
	var expected GeoCoord

	// Azimuth doesn't really matter in this case. Any azimuth from the
	// north pole is south, any azimuth from the south pole is north.

	setGeoDegs(&start, 90, 0)
	setGeoDegs(&expected, -90, 0)
	_geoAzDistanceRads(&start, degsToRads(12), degsToRads(180), &out)
	require.True(t, geoAlmostEqual(&expected, &out),
		"some direction to south pole produces south pole")
	setGeoDegs(&start, -90, 0)
	setGeoDegs(&expected, 90, 0)
	_geoAzDistanceRads(&start, degsToRads(34), degsToRads(180), &out)
	require.True(t, geoAlmostEqual(&expected, &out), "some direction to north pole produces north pole")
}

func Test__geoAzDistanceRads_invertible(t *testing.T) {
	var start GeoCoord
	setGeoDegs(&start, 15, 10)

	var out GeoCoord
	azimuth := degsToRads(20)
	degrees180 := degsToRads(180)
	distance := degsToRads(15)

	_geoAzDistanceRads(&start, azimuth, distance, &out)

	require.True(t, math.Abs(_geoDistRads(&start, &out)-distance) < EPSILON_RAD, "moved distance is as expected")
	var start2 GeoCoord = out
	_geoAzDistanceRads(&start2, azimuth+degrees180, distance, &out)
	// TODO: Epsilon is relatively large
	require.True(t, _geoDistRads(&start, &out) < 0.01, "moved back to origin")
}

func Test__geoDistRads_wrappedLongitude(t *testing.T) {
	negativeLongitude := GeoCoord{Lat: -(M_PI + M_PI_2)}
	zero := GeoCoord{}

	require.True(t, math.Abs(M_PI_2-_geoDistRads(&negativeLongitude, &zero)) < EPSILON_RAD, "Distance with wrapped longitude")
	require.True(t, math.Abs(M_PI_2-_geoDistRads(&zero, &negativeLongitude)) < EPSILON_RAD, "Distance with wrapped longitude and swapped arguments")
}

func Test_doubleConstants(t *testing.T) {
	// Simple checks for ordering of values
	testDecreasingFunction(t, hexAreaKm2, "hexAreaKm2 ordering")
	testDecreasingFunction(t, hexAreaM2, "hexAreaM2 ordering")
	testDecreasingFunction(t, edgeLengthKm, "edgeLengthKm ordering")
	testDecreasingFunction(t, edgeLengthM, "edgeLengthM ordering")
}

func testDecreasingFunction(t *testing.T, fn func(i int) float64, message string) {
	last := float64(0)
	next := float64(0)

	for i := MAX_H3_RES; i >= 0; i-- {
		next = fn(i)
		require.True(t, next > last, message)
		last = next
	}
}

func Test_intConstants(t *testing.T) {
	// Simple checks for ordering of values
	last := 0
	next := 0
	for i := 0; i <= MAX_H3_RES; i++ {
		next = int(numHexagons(i))
		require.True(t, next > last, "numHexagons ordering")
		last = next
	}
}
