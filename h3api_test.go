package h3

import (
	"bytes"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
	"text/scanner"

	"github.com/stretchr/testify/require"
)

func TestH3Api(t *testing.T) {
	t.Run("geoToH3_res", func(t *testing.T) {
		anywhere := GeoCoord{0, 0}
		require.True(t, geoToH3(&anywhere, -1) == 0, "resolution below 0 is invalid")
		require.True(t, geoToH3(&anywhere, 16) == 0, "resolution above 15 is invalid")
	})

	t.Run("geoToH3_coord", func(t *testing.T) {
		invalidLat := GeoCoord{math.NaN(), 0}
		invalidLon := GeoCoord{0, math.NaN()}
		invalidLatLon := GeoCoord{math.Inf(1), math.Inf(-1)}
		require.True(t, geoToH3(&invalidLat, 1) == 0, "invalid latitude is rejected")
		require.True(t, geoToH3(&invalidLon, 1) == 0, "invalid longitude is rejected")
		require.True(t, geoToH3(&invalidLatLon, 1) == 0, "coordinates with infinity are rejected")
	})

	t.Run("h3ToGeoBoundary_classIIIEdgeVertex", func(t *testing.T) {
		// Bug test for https://github.com/uber/h3/issues/45
		hexes := []string{
			"894cc5349b7ffff", "894cc534d97ffff", "894cc53682bffff",
			"894cc536b17ffff", "894cc53688bffff", "894cead92cbffff",
			"894cc536537ffff", "894cc5acbabffff", "894cc536597ffff",
		}
		numHexes := len(hexes)
		var h3 H3Index
		for i := 0; i < numHexes; i++ {
			h3 = stringToH3(hexes[i])
			var b GeoBoundary
			h3ToGeoBoundary(h3, &b)
			require.True(t, b.numVerts == 7, "got expected vertex count")
		}
	})

	t.Run("h3ToGeoBoundary_classIIIEdgeVertex_exact", func(t *testing.T) {
		// Bug test for https://github.com/uber/h3/issues/45
		h3 := stringToH3("894cc536537ffff")
		var boundary GeoBoundary
		boundary.numVerts = 7
		boundary.Verts = make([]GeoCoord, 7)

		setGeoDegs(&boundary.Verts[0], 18.043333154, -66.27836523500002)
		setGeoDegs(&boundary.Verts[1], 18.042238363, -66.27929062800001)
		setGeoDegs(&boundary.Verts[2], 18.040818259, -66.27854193899998)
		setGeoDegs(&boundary.Verts[3], 18.040492975, -66.27686786700002)
		setGeoDegs(&boundary.Verts[4], 18.041040385, -66.27640518300001)
		setGeoDegs(&boundary.Verts[5], 18.041757122, -66.27596711500001)
		setGeoDegs(&boundary.Verts[6], 18.043007860, -66.27669118199998)
		assertBoundary(t, h3, &boundary)
	})

	t.Run("h3ToGeoBoundary_coslonConstrain", func(t *testing.T) {
		// Bug test for https://github.com/uber/h3/issues/212
		h3 := H3Index(0x87dc6d364ffffff)
		var boundary GeoBoundary
		boundary.numVerts = 6
		boundary.Verts = make([]GeoCoord, 6)

		setGeoDegs(&boundary.Verts[0], -52.0130533678236091, -34.6232931343713091)
		setGeoDegs(&boundary.Verts[1], -52.0041156384652012, -34.6096733160584549)
		setGeoDegs(&boundary.Verts[2], -51.9929610229502472, -34.6165157145896387)
		setGeoDegs(&boundary.Verts[3], -51.9907410568096608, -34.6369680004259877)
		setGeoDegs(&boundary.Verts[4], -51.9996738734672377, -34.6505896528323660)
		setGeoDegs(&boundary.Verts[5], -52.0108315681413629, -34.6437571897165668)
		assertBoundary(t, h3, &boundary)
	})

}

var centerFiles []string
var cellFiles []string

func init() {
	fileList, err := ioutil.ReadDir("./testdata")
	if err != nil {
		panic(err)
	}

	for _, f := range fileList {
		name := f.Name()

		if strings.Contains(name, "centers") {
			centerFiles = append(centerFiles, "./testdata/"+f.Name())
		}

		if strings.Contains(name, "cells") {
			cellFiles = append(cellFiles, "./testdata/"+f.Name())
		}
	}
}

func _Test_testdata(t *testing.T) {
	t.Run("centers", func(t *testing.T) {
		for _, f := range centerFiles {
			data := loadData(f)

			for i, coords := range data {
				h := i

				if len(coords) == 1 {
					expectG := &coords[0]

					g := GeoCoord{}
					h3ToGeo(h, &g)

					require.True(t, geoDegreeEqual(&g, expectG), "h3ToGeo %x %d %v %v", h, h, expectG.AsDegrees(), g.AsDegrees())
					require.Equal(t, h, geoToH3(expectG, h3GetResolution(h)), "h3ToGeo %x %d", h, h)
				}
			}
		}
	})

	t.Run("cells", func(t *testing.T) {
		for _, f := range cellFiles {
			data := loadData(f)

			for i, coords := range data {
				h := i

				if len(coords) > 0 {
					expectGB := &GeoBoundary{numVerts: len(coords), Verts: coords}

					assertBoundary(t, h, expectGB)
				}
			}
		}
	})
}

func assertBoundary(t *testing.T, h H3Index, expectGB *GeoBoundary) {
	gb := GeoBoundary{}
	h3ToGeoBoundary(h, &gb)

	require.Equal(t, expectGB.numVerts, gb.numVerts, "h3ToGeoBoundary %x %v %v", h, gb.AsDegrees(), expectGB.AsDegrees())

	for i := 0; i < expectGB.numVerts; i++ {
		require.True(t, geoDegreeEqual(&expectGB.Verts[i], &gb.Verts[i]), "h3ToGeoBoundary %x &v &v", h, expectGB.Verts[i].AsDegrees(), gb.Verts[i].AsDegrees())
	}
}

func geoDegreeEqual(p1 *GeoCoord, p2 *GeoCoord) bool {
	return geoAlmostEqualThreshold(p1.AsDegrees(), p2.AsDegrees(), 0.000001)
}

type expects map[H3Index][]GeoCoord

func loadData(filename string) expects {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	return parseInputs(f)
}

func parseInputs(reader io.Reader) expects {
	set := expects{}

	var s scanner.Scanner
	s.Init(reader)

	tmp := bytes.NewBuffer(nil)

	clearTmp := func() {
		tmp = bytes.NewBuffer(nil)
	}

	var h uint64
	var coord *GeoCoord
	var inCell bool

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Next() {
		switch tok {
		case '\n', ' ':
			if tmp.Len() == 0 {
				break
			}

			if h == 0 {
				h, _ = strconv.ParseUint(tmp.String(), 16, 64)
				clearTmp()
				set[H3Index(h)] = []GeoCoord{}
				break
			}

			if coord == nil {
				coord = &GeoCoord{}
			}

			if coord.Lat == 0 {
				f, _ := strconv.ParseFloat(tmp.String(), 64)
				clearTmp()
				coord.Lat = f * deg2rad
				break
			}

			if coord.Lon == 0 {
				f, _ := strconv.ParseFloat(tmp.String(), 64)
				clearTmp()
				coord.Lon = f * deg2rad

				set[H3Index(h)] = append(set[H3Index(h)], *coord)

				// clear
				coord = nil

				if !inCell {
					h = 0
				}
			}
		case '{':
			inCell = true
		case '}':
			inCell = false
			h = 0
		default:
			tmp.WriteRune(tok)
		}
	}

	return set
}
