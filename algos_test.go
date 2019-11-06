package h3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_h3SetToLinkedGeo(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var polygon LinkedGeoPolygon
		h3SetToLinkedGeo(nil, 0, &polygon)
		require.True(t, countLinkedLoops(&polygon) == 0, "No loops added to polygon")
		destroyLinkedPolygon(&polygon)
	})

	t.Run("singleHex", func(t *testing.T) {
		var polygon LinkedGeoPolygon
		set := []H3Index{0x890dab6220bffff}
		numHexes := len(set)

		h3SetToLinkedGeo(set, numHexes, &polygon)
		require.True(t, countLinkedLoops(&polygon) == 1, "1 loop added to polygon")
		require.True(t, countLinkedCoords(polygon.first) == 6,
			"6 coords added to loop")
		destroyLinkedPolygon(&polygon)
	})

	t.Run("contiguous2", func(t *testing.T) {
		var polygon LinkedGeoPolygon
		set := []H3Index{0x8928308291bffff, 0x89283082957ffff}
		numHexes := len(set)

		h3SetToLinkedGeo(set, numHexes, &polygon)
		require.True(t, countLinkedLoops(&polygon) == 1, "1 loop added to polygon")
		require.True(t, countLinkedCoords(polygon.first) == 10,
			"All coords added to loop except 2 shared")
		destroyLinkedPolygon(&polygon)
	})

	// TODO: This test asserts incorrect behavior - we should be creating
	// multiple polygons, each with their own single loop. Update when the
	// algorithm is corrected.
	t.Run("nonContiguous2", func(t *testing.T) {
		var polygon LinkedGeoPolygon
		set := []H3Index{0x8928308291bffff, 0x89283082943ffff}
		numHexes := len(set)

		h3SetToLinkedGeo(set, numHexes, &polygon)
		require.True(t, countLinkedPolygons(&polygon) == 2, "2 polygons added")
		require.True(t, countLinkedLoops(&polygon) == 1, "1 loop on the first polygon")
		require.True(t, countLinkedCoords(polygon.first) == 6, "All coords for one hex added to first loop")
		require.True(t, countLinkedLoops(polygon.next) == 1, "Loop count on second polygon correct")
		require.True(t, countLinkedCoords(polygon.next.first) == 6, "All coords for one hex added to second polygon")

		destroyLinkedPolygon(&polygon)
	})

	t.Run("contiguous3", func(t *testing.T) {
		var polygon LinkedGeoPolygon
		set := []H3Index{0x8928308288bffff, 0x892830828d7ffff, 0x8928308289bffff}
		numHexes := len(set)

		h3SetToLinkedGeo(set, numHexes, &polygon)
		require.True(t, countLinkedLoops(&polygon) == 1, "1 loop added to polygon")
		require.True(t, countLinkedCoords(polygon.first) == 12, "All coords added to loop except 6 shared")
		destroyLinkedPolygon(&polygon)
	})

	t.Run("hole", func(t *testing.T) {
		var polygon LinkedGeoPolygon
		set := []H3Index{0x892830828c7ffff, 0x892830828d7ffff, 0x8928308289bffff, 0x89283082813ffff, 0x8928308288fffff, 0x89283082883ffff}
		numHexes := len(set)

		h3SetToLinkedGeo(set, numHexes, &polygon)
		require.True(t, countLinkedLoops(&polygon) == 2, "2 loops added to polygon")
		require.True(t, countLinkedCoords(polygon.first) == 6*3, "All outer coords added to first loop")
		require.True(t, countLinkedCoords(polygon.first.next) == 6, "All inner coords added to second loop")
		destroyLinkedPolygon(&polygon)
	})

	t.Run("pentagon", func(t *testing.T) {
		var polygon LinkedGeoPolygon
		set := []H3Index{0x851c0003fffffff}
		numHexes := len(set)

		h3SetToLinkedGeo(set, numHexes, &polygon)
		require.Equal(t, 1, countLinkedLoops(&polygon), "1 loop added to polygon")
		require.Equal(t, 10, countLinkedCoords(polygon.first), "10 coords (distorted pentagon) added to loop")
		destroyLinkedPolygon(&polygon)
	})

	t.Run("2Ring", func(t *testing.T) {
		var polygon LinkedGeoPolygon
		// 2-ring, in order returned by k-ring algo
		set := []H3Index{
			0x8930062838bffff, 0x8930062838fffff, 0x89300628383ffff,
			0x8930062839bffff, 0x893006283d7ffff, 0x893006283c7ffff,
			0x89300628313ffff, 0x89300628317ffff, 0x893006283bbffff,
			0x89300628387ffff, 0x89300628397ffff, 0x89300628393ffff,
			0x89300628067ffff, 0x8930062806fffff, 0x893006283d3ffff,
			0x893006283c3ffff, 0x893006283cfffff, 0x8930062831bffff,
			0x89300628303ffff}
		numHexes := len(set)

		h3SetToLinkedGeo(set, numHexes, &polygon)
		require.True(t, countLinkedLoops(&polygon) == 1, "1 loop added to polygon")
		require.True(t, countLinkedCoords(polygon.first) == (6*(2*2+1)), "Expected number of coords added to loop")
		destroyLinkedPolygon(&polygon)
	})

	t.Run("2RingUnordered", func(t *testing.T) {
		var polygon LinkedGeoPolygon
		// 2-ring in random order
		set := []H3Index{
			0x89300628393ffff, 0x89300628383ffff, 0x89300628397ffff,
			0x89300628067ffff, 0x89300628387ffff, 0x893006283bbffff,
			0x89300628313ffff, 0x893006283cfffff, 0x89300628303ffff,
			0x89300628317ffff, 0x8930062839bffff, 0x8930062838bffff,
			0x8930062806fffff, 0x8930062838fffff, 0x893006283d3ffff,
			0x893006283c3ffff, 0x8930062831bffff, 0x893006283d7ffff,
			0x893006283c7ffff}
		numHexes := len(set)

		h3SetToLinkedGeo(set, numHexes, &polygon)
		require.True(t, countLinkedLoops(&polygon) == 1, "1 loop added to polygon")
		require.True(t, countLinkedCoords(polygon.first) == (6*(2*2+1)),
			"Expected number of coords added to loop")
		destroyLinkedPolygon(&polygon)
	})

	t.Run("nestedDonut", func(t *testing.T) {
		var polygon LinkedGeoPolygon
		// hollow 1-ring + hollow 3-ring around the same hex
		set := []H3Index{
			0x89283082813ffff, 0x8928308281bffff, 0x8928308280bffff,
			0x8928308280fffff, 0x89283082807ffff, 0x89283082817ffff,
			0x8928308289bffff, 0x892830828d7ffff, 0x892830828c3ffff,
			0x892830828cbffff, 0x89283082853ffff, 0x89283082843ffff,
			0x8928308284fffff, 0x8928308287bffff, 0x89283082863ffff,
			0x89283082867ffff, 0x8928308282bffff, 0x89283082823ffff,
			0x89283082837ffff, 0x892830828afffff, 0x892830828a3ffff,
			0x892830828b3ffff, 0x89283082887ffff, 0x89283082883ffff}
		numHexes := len(set)

		h3SetToLinkedGeo(set, numHexes, &polygon)

		// Note that the polygon order here is arbitrary, making this test
		// somewhat brittle, but it's difficult to assert correctness otherwise
		require.True(t, countLinkedPolygons(&polygon) == 2, "Polygon count correct")
		require.True(t, countLinkedLoops(&polygon) == 2, "Loop count on first polygon correct")
		require.True(t, countLinkedCoords(polygon.first) == 42, "Got expected big outer loop")
		require.True(t, countLinkedCoords(polygon.first.next) == 30, "Got expected big inner loop")
		require.True(t, countLinkedLoops(polygon.next) == 2, "Loop count on second polygon correct")
		require.True(t, countLinkedCoords(polygon.next.first) == 18, "Got expected outer loop")
		require.True(t, countLinkedCoords(polygon.next.first.next) == 6, "Got expected inner loop")
		destroyLinkedPolygon(&polygon)
	})

	t.Run("nestedDonutTransmeridian", func(t *testing.T) {
		var polygon LinkedGeoPolygon
		// hollow 1-ring + hollow 3-ring around the hex at (0, -180)
		set := []H3Index{
			0x897eb5722c7ffff, 0x897eb5722cfffff, 0x897eb572257ffff,
			0x897eb57220bffff, 0x897eb572203ffff, 0x897eb572213ffff,
			0x897eb57266fffff, 0x897eb5722d3ffff, 0x897eb5722dbffff,
			0x897eb573537ffff, 0x897eb573527ffff, 0x897eb57225bffff,
			0x897eb57224bffff, 0x897eb57224fffff, 0x897eb57227bffff,
			0x897eb572263ffff, 0x897eb572277ffff, 0x897eb57223bffff,
			0x897eb572233ffff, 0x897eb5722abffff, 0x897eb5722bbffff,
			0x897eb572287ffff, 0x897eb572283ffff, 0x897eb57229bffff}
		numHexes := len(set)

		h3SetToLinkedGeo(set, numHexes, &polygon)

		// Note that the polygon order here is arbitrary, making this test
		// somewhat brittle, but it's difficult to assert correctness otherwise
		require.True(t, countLinkedPolygons(&polygon) == 2, "Polygon count correct")
		require.True(t, countLinkedLoops(&polygon) == 2,
			"Loop count on first polygon correct")
		require.True(t, countLinkedCoords(polygon.first) == 18,
			"Got expected outer loop")
		require.True(t, countLinkedCoords(polygon.first.next) == 6,
			"Got expected inner loop")
		require.True(t, countLinkedLoops(polygon.next) == 2,
			"Loop count on second polygon correct")
		require.True(t, countLinkedCoords(polygon.next.first) == 42,
			"Got expected big outer loop")
		require.True(t, countLinkedCoords(polygon.next.first.next) == 30,
			"Got expected big inner loop")
		destroyLinkedPolygon(&polygon)
	})

	t.Run("contiguous2distorted", func(t *testing.T) {
		var polygon LinkedGeoPolygon
		set := []H3Index{0x894cc5365afffff, 0x894cc536537ffff}
		numHexes := len(set)

		h3SetToLinkedGeo(set, numHexes, &polygon)
		require.True(t, countLinkedLoops(&polygon) == 1, "1 loop added to polygon")
		require.True(t, countLinkedCoords(polygon.first) == 12,
			"All coords added to loop except 2 shared")
		destroyLinkedPolygon(&polygon)
	})

	t.Run("negativeHashedCoordinates", func(t *testing.T) {
		var polygon LinkedGeoPolygon
		set := []H3Index{0x88ad36c547fffff, 0x88ad36c467fffff}
		numHexes := len(set)
		h3SetToLinkedGeo(set, numHexes, &polygon)
		require.True(t, countLinkedPolygons(&polygon) == 2, "2 polygons added")
		require.True(t, countLinkedLoops(&polygon) == 1,
			"1 loop on the first polygon")
		require.True(t, countLinkedCoords(polygon.first) == 6,
			"All coords for one hex added to first loop")
		require.True(t, countLinkedLoops(polygon.next) == 1,
			"Loop count on second polygon correct")
		require.True(t, countLinkedCoords(polygon.next.first) == 6,
			"All coords for one hex added to second polygon")
		destroyLinkedPolygon(&polygon)
	})
}

func Test_h3SetToVertexGraph(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var graph VertexGraph
		h3SetToVertexGraph(nil, 0, &graph)
		require.True(t, graph.size == 0, "No edges added to graph")
		destroyVertexGraph(&graph)
	})

	t.Run("singleHex", func(t *testing.T) {
		var graph VertexGraph
		set := []H3Index{0x890dab6220bffff}
		numHexes := len(set)
		h3SetToVertexGraph(set, numHexes, &graph)
		require.True(t, graph.size == 6, "All edges of one hex added to graph")
		destroyVertexGraph(&graph)
	})

	t.Run("nonContiguous2", func(t *testing.T) {
		var graph VertexGraph
		set := []H3Index{0x8928308291bffff, 0x89283082943ffff}
		numHexes := len(set)
		h3SetToVertexGraph(set, numHexes, &graph)

		require.Equal(t, 12, graph.size, "All edges of two non-contiguous hexes added to graph")
		destroyVertexGraph(&graph)
	})

	t.Run("contiguous2", func(t *testing.T) {
		var graph VertexGraph
		set := []H3Index{0x8928308291bffff, 0x89283082957ffff}
		numHexes := len(set)
		h3SetToVertexGraph(set, numHexes, &graph)

		require.Equal(t, 10, graph.size, "All edges except 2 shared added to graph")
		destroyVertexGraph(&graph)
	})

	t.Run("contiguous2distorted", func(t *testing.T) {
		var graph VertexGraph
		set := []H3Index{0x894cc5365afffff, 0x894cc536537ffff}
		numHexes := len(set)

		h3SetToVertexGraph(set, numHexes, &graph)
		require.Equal(t, 12, graph.size, "All edges except 2 shared added to graph")
		destroyVertexGraph(&graph)
	})

	t.Run("contiguous3", func(t *testing.T) {
		var graph VertexGraph
		set := []H3Index{0x8928308288bffff, 0x892830828d7ffff,
			0x8928308289bffff}
		numHexes := len(set)

		h3SetToVertexGraph(set, numHexes, &graph)
		require.True(t, graph.size == 3*4, "All edges except 6 shared added to graph")
		destroyVertexGraph(&graph)
	})

	t.Run("hole", func(t *testing.T) {
		var graph VertexGraph
		set := []H3Index{
			0x892830828c7ffff,
			0x892830828d7ffff,
			0x8928308289bffff,
			0x89283082813ffff,
			0x8928308288fffff,
			0x89283082883ffff,
		}
		numHexes := len(set)
		h3SetToVertexGraph(set, numHexes, &graph)
		require.True(t, graph.size == (6*3)+6, "All outer edges and inner hole edges added to graph")
		destroyVertexGraph(&graph)
	})
}
