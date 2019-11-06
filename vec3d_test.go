package h3

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_pointSquareDist(t *testing.T) {
	v1 := Vec3d{0, 0, 0}
	v2 := Vec3d{1, 0, 0}
	v3 := Vec3d{0, 1, 1}
	v4 := Vec3d{1, 1, 1}
	v5 := Vec3d{1, 1, 2}

	require.True(t, math.Abs(_pointSquareDist(&v1, &v1)) < math.SmallestNonzeroFloat64, "distance to self is 0")
	require.True(t, math.Abs(_pointSquareDist(&v1, &v2)-1) < math.SmallestNonzeroFloat64, "distance to <1,0,0> is 1")
	require.True(t, math.Abs(_pointSquareDist(&v1, &v3)-2) < math.SmallestNonzeroFloat64, "distance to <0,1,1> is 2")
	require.True(t, math.Abs(_pointSquareDist(&v1, &v4)-3) < math.SmallestNonzeroFloat64, "distance to <1,1,1> is 3")
	require.True(t, math.Abs(_pointSquareDist(&v1, &v5)-6) < math.SmallestNonzeroFloat64, "distance to <1,1,2> is 6")
}

func Test_geoToVec3d(t *testing.T) {
	origin := Vec3d{}
	c1 := GeoCoord{0, 0}
	var p1 Vec3d
	_geoToVec3d(&c1, &p1)
	require.True(t, math.Abs(_pointSquareDist(&origin, &p1)-1) < EPSILON_RAD, "Geo point is on the unit sphere")

	c2 := GeoCoord{M_PI_2, 0}
	var p2 Vec3d
	_geoToVec3d(&c2, &p2)
	require.True(t, math.Abs(_pointSquareDist(&p1, &p2)-2) < EPSILON_RAD, "Geo point is on another axis")

	c3 := GeoCoord{M_PI, 0}
	var p3 Vec3d
	_geoToVec3d(&c3, &p3)
	require.True(t, math.Abs(_pointSquareDist(&p1, &p3)-4) < EPSILON_RAD, "Geo point is the other side of the sphere")
}
