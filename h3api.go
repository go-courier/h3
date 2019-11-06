package h3

import (
	"bytes"
	"fmt"
)

// Maximum number of cell boundary vertices; worst case is pentagon:
// 5 original Verts + 5 edge crossings
const MAX_CELL_BNDRY_VERTS = 10

/**
  @brief latitude/longitude in radians
*/
type GeoCoord struct {
	Lat float64 ///< latitude in radians
	Lon float64 ///< longitude in radians
}

func (g GeoCoord) String() string {
	return fmt.Sprintf("%f,%f", g.Lat, g.Lon)
}

func (g GeoCoord) AsDegrees() *GeoCoord {
	return &GeoCoord{
		Lon: normalizeDegree(radsToDegs(g.Lon), -180.0, 180),
		Lat: normalizeDegree(radsToDegs(g.Lat), -90, 90),
	}
}

func (g GeoCoord) AsRadians() *GeoCoord {
	return &GeoCoord{
		Lon: radsToDegs(g.Lon),
		Lat: radsToDegs(g.Lat),
	}
}

/**
  @brief cell boundary in latitude/longitude
*/
type GeoBoundary struct {
	numVerts int        ///< number of vertices
	Verts    []GeoCoord ///< vertices in ccw order
}

func (gb GeoBoundary) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteRune('[')
	for i := range gb.Verts {
		if i != 0 {
			buf.WriteRune(' ')
		}
		buf.WriteString(gb.Verts[i].String())
	}
	buf.WriteRune(']')
	return buf.String()
}

func (gb GeoBoundary) AsDegrees() *GeoBoundary {
	list := make([]GeoCoord, len(gb.Verts))

	for i := range gb.Verts {
		list[i] = *(gb.Verts[i].AsDegrees())
	}

	gb.Verts = list

	return &gb
}

func (gb GeoBoundary) AsRadians() *GeoBoundary {
	list := make([]GeoCoord, len(gb.Verts))

	for i := range gb.Verts {
		list[i] = *(gb.Verts[i].AsDegrees())
	}

	gb.Verts = list

	return &gb
}

/**
 *  @brief similar to GeoBoundary, but requires more alloc work
 */
type Geofence struct {
	numVerts int
	verts    []GeoCoord
}

func (g *Geofence) IsZero() bool {
	return g == nil || g.numVerts == 0
}

func (g *Geofence) NewIterate() func(vertexA *GeoCoord, vertexB *GeoCoord) bool {
	loopIndex := -1

	return func(vertexA *GeoCoord, vertexB *GeoCoord) bool {
		loopIndex++

		if loopIndex >= g.numVerts {
			return false
		}

		*vertexA = g.verts[loopIndex]
		*vertexB = g.verts[(loopIndex+1)%g.numVerts]

		return true
	}
}

/**
 *  @brief Simplified core of GeoJSON Polygon coordinates definition
 */
type GeoPolygon struct {
	geofence Geofence   ///< exterior boundary of the polygon
	numHoles int        ///< number of elements in the array pointed to by holes
	holes    []Geofence ///< interior boundaries (holes) in the polygon
}

/**
 *  @brief Simplified core of GeoJSON MultiPolygon coordinates definition
 */
type GeoMultiPolygon struct {
	numPolygons int
	polygons    []GeoPolygon
}

/**
 *  @brief A coordinate node in a linked geo structure, part of a linked list
 */
type LinkedGeoCoord struct {
	vertex GeoCoord
	next   *LinkedGeoCoord
}

/**
 *  @brief A loop node in a linked geo structure, part of a linked list
 */
type LinkedGeoLoop struct {
	first *LinkedGeoCoord
	last  *LinkedGeoCoord
	next  *LinkedGeoLoop
}

func (l *LinkedGeoLoop) IsZero() bool {
	return l == nil || l.first == nil
}

func (l *LinkedGeoLoop) NewIterate() func(vertexA *GeoCoord, vertexB *GeoCoord) bool {
	var currentCoord, nextCoord *LinkedGeoCoord

	return func(vertexA *GeoCoord, vertexB *GeoCoord) bool {
		var getNextCoord = func(coordToCheck *LinkedGeoCoord) *LinkedGeoCoord {
			if coordToCheck == nil {
				return l.first
			}
			return currentCoord.next

		}

		currentCoord = getNextCoord(currentCoord)

		if currentCoord == nil {
			return false
		}

		*vertexA = currentCoord.vertex
		nextCoord = getNextCoord(currentCoord.next)
		*vertexB = nextCoord.vertex

		return true
	}
}

/**
 *  @brief A polygon node in a linked geo structure, part of a linked list.
 */
type LinkedGeoPolygon struct {
	first *LinkedGeoLoop
	last  *LinkedGeoLoop
	next  *LinkedGeoPolygon
}

/**
 * @brief IJ hexagon coordinates
 *
 * Each axis is spaced 120 degrees apart.
 */
type CoordIJ struct {
	i int ///< i component
	j int ///< j component
}
