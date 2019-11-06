package h3

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

const deg2rad = math.Pi / 180.0
const rad2deg = 180.0 / math.Pi

func GeoFromWGS84(lat float64, lon float64) *GeoCoord {
	return &GeoCoord{lat * deg2rad, lon * deg2rad}
}

func Test_geoToH3(t *testing.T) {

	t.Run("geoToH3", func(t *testing.T) {
		require.Equal(t, H3Index(614553222213795839), geoToH3(GeoFromWGS84(0, 0), 8))
		require.Equal(t, H3Index(613287273236004863), geoToH3(GeoFromWGS84(45, 45), 8))
		require.Equal(t, H3Index(612544946678792191), geoToH3(GeoFromWGS84(90, 90), 8))
	})

	t.Run("geoToH3 2", func(t *testing.T) {
		g := GeoCoord{0.659966917655, 2*3.14159 - 2.1364398519396}

		expects := []uint64{
			577199624117288959,
			581672437419081727,
			586175487290638335,
			590678605881671679,
			595182179739238399,
			599685771850416127,
			604189370672480255,
			608692970266296319,
			613196569891569663,
			617700169517629439,
			622203769144967167,
			626707368772308991,
			631210968399678463,
			635714568027048703,
			640218167654419175,
			644721767281789666, // max 15
			0,
			0,
			0,
		}

		for i := range expects {
			require.Equal(t, H3Index(expects[i]), geoToH3(&g, i))
		}
	})

	t.Run("_h3ToFaceIjk & _faceIjkToH3", func(t *testing.T) {
		h := H3Index(613219835716829183)
		f := FaceIJK{}
		_h3ToFaceIjk(h, &f)
		require.Equal(t, h, _faceIjkToH3(&f, 8))
	})

	t.Run("_geoToFaceIjk & _faceIjkToGeo", func(t *testing.T) {
		g := GeoCoord{0, 0}

		f := FaceIJK{}
		_geoToFaceIjk(&g, 15, &f)
		g2 := GeoCoord{}
		_faceIjkToGeo(&f, 15, &g2)
	})
}

func Test_getPentagonIndexes(t *testing.T) {
	expectedCount := pentagonIndexCount()

	for res := 0; res <= 15; res++ {
		h3Indexes := make([]H3Index, 0)
		getPentagonIndexes(res, &h3Indexes)

		for _, h3Index := range h3Indexes {
			if h3Index != 0 {
				require.True(t, h3IsValid(h3Index), "index should be valid")
				require.True(t, h3IsPentagon(h3Index), "index should be pentagon")
				require.True(t, h3GetResolution(h3Index) == res, "index should have correct resolution")
			}
		}

		require.True(t, len(h3Indexes) == expectedCount, "there should be exactly 12 pentagons")
	}
}

func Test_GeoToH3(t *testing.T) {
	t.Run("geoToH3ExtremeCoordinates", func(t *testing.T) {
		// Check that none of these cause crashes.
		g := GeoCoord{0, 1e45}
		t.Log(geoToH3(&g, 14)) // 641677981140798679

		g2 := GeoCoord{1e46, 1e45}
		t.Log(geoToH3(&g2, 15))

		var g4 GeoCoord
		setGeoDegs(&g4, 2, -3e39)

		t.Log(geoToH3(&g4, 0))
	})

	t.Run("faceIjkToH3ExtremeCoordinates", func(t *testing.T) {
		fijk0I := FaceIJK{0, CoordIJK{3, 0, 0}}
		require.True(t, _faceIjkToH3(&fijk0I, 0) == 0, "i out of bounds at res 0")

		fijk0J := FaceIJK{1, CoordIJK{0, 4, 0}}
		require.True(t, _faceIjkToH3(&fijk0J, 0) == 0, "j out of bounds at res 0")

		fijk0K := FaceIJK{2, CoordIJK{2, 0, 5}}
		require.True(t, _faceIjkToH3(&fijk0K, 0) == 0, "k out of bounds at res 0")

		fijk1I := FaceIJK{3, CoordIJK{6, 0, 0}}
		require.True(t, _faceIjkToH3(&fijk1I, 1) == 0, "i out of bounds at res 1")

		fijk1J := FaceIJK{4, CoordIJK{0, 7, 1}}
		require.True(t, _faceIjkToH3(&fijk1J, 1) == 0, "j out of bounds at res 1")

		fijk1K := FaceIJK{5, CoordIJK{2, 0, 8}}
		require.True(t, _faceIjkToH3(&fijk1K, 1) == 0, "k out of bounds at res 1")

		fijk2I := FaceIJK{6, CoordIJK{18, 0, 0}}
		require.True(t, _faceIjkToH3(&fijk2I, 2) == 0, "i out of bounds at res 2")

		fijk2J := FaceIJK{7, CoordIJK{0, 19, 1}}
		require.True(t, _faceIjkToH3(&fijk2J, 2) == 0, "j out of bounds at res 2")

		fijk2K := FaceIJK{8, CoordIJK{2, 0, 20}}
		require.True(t, _faceIjkToH3(&fijk2K, 2) == 0, "k out of bounds at res 2")
	})

	t.Run("h3IsValidAtResolution", func(t *testing.T) {
		for i := 0; i <= MAX_H3_RES; i++ {
			geoCoord := GeoCoord{0, 0}

			h3 := geoToH3(&geoCoord, i)

			require.True(t, h3IsValid(h3), "h3IsValid failed on resolution %d", i)
		}
	})

	t.Run("h3IsValidDigits", func(t *testing.T) {
		geoCoord := GeoCoord{0, 0}
		h3 := geoToH3(&geoCoord, 1)
		// Set a bit for an unused digit to something else.
		h3 ^= 1
		require.True(t, !h3IsValid(h3), "h3IsValid failed on invalid unused digits")
	})

	t.Run("h3IsValidBaseCell", func(t *testing.T) {
		for i := 0; i < NUM_BASE_CELLS; i++ {
			h := H3_INIT
			H3_SET_MODE(&h, H3_HEXAGON_MODE)
			H3_SET_BASE_CELL(&h, i)
			require.True(t, h3IsValid(h), "h3IsValid failed on base cell %d", i)
			require.True(t, h3GetBaseCell(h) == i,
				"failed to recover base cell")
		}
	})

	t.Run("h3IsValidBaseCellInvalid", func(t *testing.T) {
		hWrongBaseCell := H3_INIT
		H3_SET_MODE(&hWrongBaseCell, H3_HEXAGON_MODE)
		H3_SET_BASE_CELL(&hWrongBaseCell, NUM_BASE_CELLS)
		require.True(t, !h3IsValid(hWrongBaseCell), "h3IsValid failed on invalid base cell")
	})

	t.Run("h3IsValidWithMode", func(t *testing.T) {
		for i := 0; i <= 0xf; i++ {
			h := H3_INIT
			H3_SET_MODE(&h, H3Mode(i))
			require.True(t, !h3IsValid(h) || i == 1, "h3IsValid failed on mode %d", i)
		}
	})

	t.Run("h3BadDigitInvalid", func(t *testing.T) {
		h := H3_INIT
		// By default the first index digit is out of range.
		H3_SET_MODE(&h, H3_HEXAGON_MODE)
		H3_SET_RESOLUTION(&h, 1)
		require.True(t, !h3IsValid(h), "h3IsValid failed on too large digit")
	})

	t.Run("h3DeletedSubsequenceInvalid", func(t *testing.T) {
		var h H3Index
		// Create an index located in a deleted subsequence of a pentagon.
		setH3Index(&h, 1, 4, K_AXES_DIGIT)

		require.True(t, !h3IsValid(h), "h3IsValid failed on deleted subsequence")
	})

	t.Run("h3ToString", func(t *testing.T) {
		require.Equal(t, "880a000001fffff", h3ToString(612665471184928767))
	})

	t.Run("stringToH3", func(t *testing.T) {
		require.Equal(t, H3Index(612665471184928767), stringToH3("880a000001fffff"))
		require.True(t, stringToH3("") == 0, "got an index from nothing")
		require.True(t, stringToH3("**") == 0, "got an index from junk")
		require.True(t, stringToH3("ffffffffffffffff") == 0xffffffffffffffff, "failed on large input")
	})

	t.Run("setH3Index", func(t *testing.T) {
		var h H3Index
		setH3Index(&h, 5, 12, 1)

		require.True(t, H3_GET_RESOLUTION(h) == 5, "resolution as expected")
		require.True(t, H3_GET_BASE_CELL(h) == 12, "base cell as expected")
		require.True(t, H3_GET_MODE(h) == H3_HEXAGON_MODE, "mode as expected")

		for i := 1; i <= 5; i++ {
			require.True(t, H3_GET_INDEX_DIGIT(h, i) == 1, "digit as expected")
		}

		for i := 6; i <= MAX_H3_RES; i++ {
			require.True(t, H3_GET_INDEX_DIGIT(h, i) == 7, "blanked digit as expected")
		}

		require.True(t, h == 599405990164561919, "index matches expected")
		require.True(t, h == 0x85184927fffffff, "index matches expected")
	})

	t.Run("h3IsResClassIII", func(t *testing.T) {
		coord := GeoCoord{0, 0}
		for i := 0; i <= MAX_H3_RES; i++ {
			h := geoToH3(&coord, i)
			require.True(t, h3IsResClassIII(h) == isResClassIII(i),
				"matches existing definition")
		}
	})
}

func Test_h3ToCenterChild(t *testing.T) {
	var baseHex H3Index
	var baseCentroid GeoCoord

	setH3Index(&baseHex, 8, 4, 2)
	h3ToGeo(baseHex, &baseCentroid)

	t.Run("propertyTests", func(t *testing.T) {
		for res := 0; res <= MAX_H3_RES-1; res++ {
			for childRes := res + 1; childRes <= MAX_H3_RES; childRes++ {
				var centroid GeoCoord
				h3Index := geoToH3(&baseCentroid, res)
				h3ToGeo(h3Index, &centroid)
				geoChild := geoToH3(&centroid, childRes)
				centerChild := h3ToCenterChild(h3Index, childRes)
				require.True(t, centerChild == geoChild, "center child should be same as indexed centroid at child resolution")
				require.True(t, h3GetResolution(centerChild) == childRes, "center child should have correct resolution")
				require.True(t, h3ToParent(centerChild, res) == h3Index, "parent at original resolution should be initial index")
			}
		}
	})

	t.Run("sameRes", func(t *testing.T) {
		res := h3GetResolution(baseHex)
		require.True(t, h3ToCenterChild(baseHex, res) == baseHex, "center child at same resolution should return self")
	})

	t.Run("invalidInputs", func(t *testing.T) {
		res := h3GetResolution(baseHex)
		require.True(t, h3ToCenterChild(baseHex, res-1) == 0,
			"should fail at coarser resolution")
		require.True(t, h3ToCenterChild(baseHex, -1) == 0,
			"should fail for negative resolution")
		require.True(t, h3ToCenterChild(baseHex, MAX_H3_RES+1) == 0, "should fail beyond finest resolution")
	})
}

func Test_h3ToGeoBoundary(t *testing.T) {
	t.Run("h3ToGeo", func(t *testing.T) {
		expectGeo := GeoFromWGS84(37.812291538780364, -122.41353593838753)

		g := GeoCoord{}
		h3ToGeo(613196569891569663, &g)

		require.True(t, geoAlmostEqual(expectGeo, &g))
	})

	t.Run("h3ToGeoBoundary", func(t *testing.T) {
		expect := GeoBoundary{
			numVerts: 6,
			Verts: []GeoCoord{
				*GeoFromWGS84(37.80760100422449, -122.41208776737979),
				*GeoFromWGS84(37.81114379658359, -122.40761222203226),
				*GeoFromWGS84(37.815834307032965, -122.4090602822424),
				*GeoFromWGS84(37.816981839321244, -122.41498422661992),
				*GeoFromWGS84(37.81343893006517, -122.41945955867847),
				*GeoFromWGS84(37.808748605427716, -122.41801115969699),
			},
		}

		gb := GeoBoundary{}
		h3ToGeoBoundary(H3Index(613196569891569663), &gb)

		require.Equal(t, expect.numVerts, gb.numVerts)

		for i := 0; i < expect.numVerts; i++ {
			require.True(t, geoAlmostEqual(&expect.Verts[i], &gb.Verts[i]))
		}
	})

	t.Run("h3ToGeoBoundary pentagon", func(t *testing.T) {
		expect := GeoBoundary{
			numVerts: 10,
			Verts: []GeoCoord{
				*GeoFromWGS84(50.104450101, -143.478843877),
				*GeoFromWGS84(50.103795870, -143.480089732),
				*GeoFromWGS84(50.103371455, -143.480450779),
				*GeoFromWGS84(50.102409316, -143.479865681),
				*GeoFromWGS84(50.102057919, -143.479347956),
				*GeoFromWGS84(50.102117500, -143.477740557),
				*GeoFromWGS84(50.102324725, -143.477059533),
				*GeoFromWGS84(50.103323690, -143.476651121),
				*GeoFromWGS84(50.103803169, -143.476747929),
				*GeoFromWGS84(50.104360999, -143.478102984),
			},
		}

		gb := GeoBoundary{}
		h3ToGeoBoundary(H3Index(0x891c0000003ffff), &gb)

		require.Equal(t, expect.numVerts, gb.numVerts)

		for i := 0; i < expect.numVerts; i++ {
			require.True(t, geoAlmostEqual(&expect.Verts[i], &gb.Verts[i]))
		}
	})
}

func Test_h3ToChildren(t *testing.T) {
	sf := GeoCoord{0.659966917655, 2*3.14159 - 2.1364398519396}
	sfHex8 := geoToH3(&sf, 8) // 613196569891569663

	var verifyCountAndUniqueness = func(t *testing.T, children []H3Index, paddedCount int, expectedCount int) {
		numFound := 0
		for i := 0; i < paddedCount; i++ {
			if len(children) == i {
				break
			}

			currIndex := children[i]
			if currIndex == 0 {
				continue
			}
			numFound++

			// verify uniqueness
			indexSeen := 0
			for j := i + 1; j < paddedCount; j++ {
				if len(children) == j {
					break
				}

				if children[j] == currIndex {
					indexSeen++
				}
			}
			require.True(t, indexSeen == 0, "index was seen only once")
		}

		require.True(t, numFound == expectedCount, "got expected number of children")
	}

	t.Run("oneResStep", func(t *testing.T) {
		sfHex9s := make([]H3Index, 0)
		h3ToChildren(sfHex8, 9, &sfHex9s)

		var center GeoCoord
		h3ToGeo(sfHex8, &center)
		sfHex9_0 := geoToH3(&center, 9)

		numFound := 0
		for i := range sfHex9s {
			if sfHex9_0 == sfHex9s[i] {
				numFound++
			}
		}

		require.True(t, numFound == 1, "found the center hex")

		// Get the neighbor hexagons by averaging the center point and outer
		// points then query for those independently

		var outside GeoBoundary
		h3ToGeoBoundary(sfHex8, &outside)

		for i := 0; i < outside.numVerts; i++ {
			avg := GeoCoord{}
			avg.Lat = (outside.Verts[i].Lat + center.Lat) / 2
			avg.Lon = (outside.Verts[i].Lon + center.Lon) / 2
			avgHex9 := geoToH3(&avg, 9)
			for j := range sfHex9s {
				if avgHex9 == sfHex9s[j] {
					numFound++
				}
			}
		}

		require.True(t, numFound == 7, "found all expected children")
	})

	t.Run("multipleResSteps", func(t *testing.T) {
		// Lots of children. Will just confirm number and uniqueness
		children := make([]H3Index, 0)
		h3ToChildren(sfHex8, 10, &children)

		verifyCountAndUniqueness(t, children, 60, 49)
	})

	t.Run("sameRes", func(t *testing.T) {
		children := make([]H3Index, 0)
		h3ToChildren(sfHex8, 8, &children)
		verifyCountAndUniqueness(t, children, 7, 1)
	})

	t.Run("childResTooCoarse", func(t *testing.T) {
		children := make([]H3Index, 0)
		h3ToChildren(sfHex8, 7, &children)
		verifyCountAndUniqueness(t, children, 7, 0)
	})

	t.Run("childResTooFine", func(t *testing.T) {
		sfHexMax := geoToH3(&sf, MAX_H3_RES)
		children := make([]H3Index, 0)
		h3ToChildren(sfHexMax, MAX_H3_RES+1, &children)
		verifyCountAndUniqueness(t, children, 7, 0)
	})

	t.Run("pentagonChildren", func(t *testing.T) {
		var pentagon H3Index
		setH3Index(&pentagon, 1, 4, 0)

		expectedCount := (5 * 7) + 6
		paddedCount := maxH3ToChildrenSize(pentagon, 3)

		children := make([]H3Index, 0)

		// todo don't why need this
		// h3ToChildren(sfHex8, 10, &children)

		h3ToChildren(pentagon, 3, &children)

		verifyCountAndUniqueness(t, children, paddedCount, expectedCount)
	})
}

func Test_maxH3ToChildrenSize(t *testing.T) {
	sf := GeoCoord{0.659966917655, 2*3.14159 - 2.1364398519396}

	parent := geoToH3(&sf, 7)

	require.True(t, maxH3ToChildrenSize(parent, 3) == 0, "got expected size for coarser res")
	require.True(t, maxH3ToChildrenSize(parent, 7) == 1, "got expected size for same res")
	require.True(t, maxH3ToChildrenSize(parent, 8) == 7, "got expected size for child res")
	require.True(t, maxH3ToChildrenSize(parent, 9) == 7*7, "got expected size for grandchild res")
}
