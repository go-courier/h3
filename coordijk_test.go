package h3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CoordIJK(t *testing.T) {
	t.Run("_unitIjkToDigit", func(t *testing.T) {
		zero := CoordIJK{}
		i := CoordIJK{1, 0, 0}
		outOfRange := CoordIJK{2, 0, 0}
		unnormalizedZero := CoordIJK{2, 2, 2}
		require.True(t, _unitIjkToDigit(&zero) == CENTER_DIGIT, "Unit IJK to zero")
		require.True(t, _unitIjkToDigit(&i) == I_AXES_DIGIT, "Unit IJK to I axis")
		require.True(t, _unitIjkToDigit(&outOfRange) == INVALID_DIGIT, "Unit IJK out of range")
		require.True(t, _unitIjkToDigit(&unnormalizedZero) == CENTER_DIGIT, "Unnormalized unit IJK to zero")
	})

	t.Run("_neighbor", func(t *testing.T) {
		ijk := CoordIJK{}
		zero := CoordIJK{}
		i := CoordIJK{1, 0, 0}
		_neighbor(&ijk, CENTER_DIGIT)
		require.True(t, _ijkMatches(&ijk, &zero), "Center neighbor is self")
		_neighbor(&ijk, I_AXES_DIGIT)
		require.True(t, _ijkMatches(&ijk, &i), "I neighbor as expected")
		_neighbor(&ijk, INVALID_DIGIT)
		require.True(t, _ijkMatches(&ijk, &i), "Invalid neighbor is self")
	})
}

func Test_CoordIJ(t *testing.T) {
	t.Run("ijkToIj_zero", func(t *testing.T) {
		ijk := CoordIJK{}
		ij := CoordIJ{}
		ijkToIj(&ijk, &ij)
		require.True(t, ij.i == 0, "ij.i zero")
		require.True(t, ij.j == 0, "ij.j zero")
		ijToIjk(&ij, &ijk)
		require.True(t, ijk.i == 0, "ijk.i zero")
		require.True(t, ijk.j == 0, "ijk.j zero")
		require.True(t, ijk.k == 0, "ijk.k zero")
	})

	t.Run("ijkToIj_roundtrip", func(t *testing.T) {
		for dir := CENTER_DIGIT; dir < NUM_DIGITS; dir++ {
			ijk := CoordIJK{}
			_neighbor(&ijk, dir)
			ij := CoordIJ{}
			ijkToIj(&ijk, &ij)
			recovered := CoordIJK{}
			ijToIjk(&ij, &recovered)
			require.True(t, _ijkMatches(&ijk, &recovered), "got same ijk coordinates back")
		}
	})

	t.Run("ijkToCube_roundtrip", func(t *testing.T) {
		for dir := CENTER_DIGIT; dir < NUM_DIGITS; dir++ {
			ijk := CoordIJK{}
			_neighbor(&ijk, dir)
			original := CoordIJK{ijk.i, ijk.j, ijk.k}

			ijkToCube(&ijk)
			cubeToIjk(&ijk)
			require.True(t, _ijkMatches(&ijk, &original), "got same ijk coordinates back")
		}
	})
}
