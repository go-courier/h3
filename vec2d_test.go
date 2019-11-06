package h3

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_v2dMag(t *testing.T) {
	v := Vec2d{3.0, 4.0}

	require.True(t, math.Abs(_v2dMag(&v)-5.0) < math.SmallestNonzeroFloat64, "magnitude as expected")
}

func Test_v2dIntersect(t *testing.T) {
	p0 := Vec2d{2.0, 2.0}
	p1 := Vec2d{6.0, 6.0}
	p2 := Vec2d{0.0, 4.0}
	p3 := Vec2d{10.0, 4.0}
	intersection := Vec2d{0.0, 0.0}
	_v2dIntersect(&p0, &p1, &p2, &p3, &intersection)

	require.True(t, math.Abs(intersection.x-4.0) < math.SmallestNonzeroFloat64, "X coord as expected")
	require.True(t, math.Abs(intersection.y-4.0) < math.SmallestNonzeroFloat64, "Y coord as expected")
}

func Test_v2dEquals(t *testing.T) {
	v1 := Vec2d{3.0, 4.0}
	v2 := Vec2d{3.0, 4.0}
	v3 := Vec2d{3.5, 4.0}
	v4 := Vec2d{3.0, 4.5}

	require.True(t, _v2dEquals(&v1, &v2), "true for equal vectors")
	require.True(t, !_v2dEquals(&v1, &v3), "false for different x")
	require.True(t, !_v2dEquals(&v1, &v4), "false for different y")
}
