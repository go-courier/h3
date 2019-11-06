package h3

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Fixtures
var center GeoCoord

var vertex1 GeoCoord

var vertex2 GeoCoord

var vertex3 GeoCoord

var vertex4 GeoCoord

var vertex5 GeoCoord

var vertex6 GeoCoord

func init() {
	setGeoDegs(&center, 37.77362016769341, -122.41673772517154)
	setGeoDegs(&vertex1, 87.372002166, 166.160981117)
	setGeoDegs(&vertex2, 87.370101364, 166.160184306)
	setGeoDegs(&vertex3, 87.369088356, 166.196239997)
	setGeoDegs(&vertex4, 87.369975080, 166.233115768)
	setGeoDegs(&vertex5, 0, 0)
	setGeoDegs(&vertex6, -10, -10)
}

func Test_makeVertexGraph(t *testing.T) {
	var graph VertexGraph
	initVertexGraph(&graph, 10, 9)

	require.True(t, graph.numBuckets == 10, "numBuckets set")
	require.True(t, graph.size == 0, "size set")

	destroyVertexGraph(&graph)
}

// TODO fixme
func _Test_vertexHash(t *testing.T) {
	//var centerIndex H3Index
	var outline GeoBoundary
	var hash1 uint32
	var hash2 uint32
	numBuckets := 1000
	for res := 0; res < 11; res++ {
		//centerIndex = geoToH3(&center, res)
		//h3ToGeoBoundary(centerIndex, &outline)
		for i := 0; i < outline.numVerts; i++ {
			hash1 = _hashVertex(&outline.Verts[i], res, numBuckets)
			hash2 = _hashVertex(&outline.Verts[(i+1)%outline.numVerts], res, numBuckets)
			require.True(t, hash1 != hash2, "Hashes must not be equal")
		}
	}
}

func Test_vertexHashNegative(t *testing.T) {
	numBuckets := 10
	require.True(t, _hashVertex(&vertex5, 5, numBuckets) < uint32(numBuckets), "zero vertex hashes correctly")
	require.True(t, _hashVertex(&vertex6, 5, numBuckets) < uint32(numBuckets), "negative coordinates vertex hashes correctly")
}

func Test_addVertexNode(t *testing.T) {
	var graph VertexGraph
	initVertexGraph(&graph, 10, 9)

	var node *VertexNode
	var addedNode *VertexNode

	t.Run("basic add", func(t *testing.T) {
		addedNode = addVertexNode(&graph, &vertex1, &vertex2)
		node = findNodeForEdge(&graph, &vertex1, &vertex2)

		require.True(t, node != nil, "Node found")
		require.True(t, node == addedNode, "Right node found")
		require.True(t, graph.size == 1, "Graph size incremented")

		t.Run("collision add", func(t *testing.T) {
			addedNode = addVertexNode(&graph, &vertex1, &vertex3)
			node = findNodeForEdge(&graph, &vertex1, &vertex3)

			require.True(t, node != nil, "Node found after hash collision")
			require.True(t, node == addedNode, "Right node found")
			require.True(t, graph.size == 2, "Graph size incremented")

			t.Run("collision add #2", func(t *testing.T) {
				addedNode = addVertexNode(&graph, &vertex1, &vertex4)
				node = findNodeForEdge(&graph, &vertex1, &vertex4)
				require.True(t, node != nil, "Node found after 2nd hash collision")
				require.True(t, node == addedNode, "Right node found")
				require.True(t, graph.size == 3, "Graph size incremented")

				t.Run("Exact match no-op", func(t *testing.T) {
					node = findNodeForEdge(&graph, &vertex1, &vertex2)
					addedNode = addVertexNode(&graph, &vertex1, &vertex2)
					require.True(t, node == findNodeForEdge(&graph, &vertex1, &vertex2), "Exact match did not change existing node")
					require.True(t, node == addedNode, "Old node returned")
					require.True(t, graph.size == 3, "Graph size was not changed")
				})
			})
		})
	})

	destroyVertexGraph(&graph)
}

func Test_addVertexNodeDupe(t *testing.T) {
	var graph VertexGraph
	initVertexGraph(&graph, 10, 9)
	var node *VertexNode
	var addedNode *VertexNode

	t.Run("basic add", func(t *testing.T) {
		addedNode = addVertexNode(&graph, &vertex1, &vertex2)
		node = findNodeForEdge(&graph, &vertex1, &vertex2)

		require.True(t, node != nil, "Node found")
		require.True(t, node == addedNode, "Right node found")
		require.True(t, graph.size == 1, "Graph size incremented")

		t.Run("dupe add", func(t *testing.T) {
			addedNode = addVertexNode(&graph, &vertex1, &vertex2)

			require.True(t, node == addedNode, "addVertexNode returned the original node")
			require.True(t, graph.size == 1, "Graph size not incremented")
		})
	})

	destroyVertexGraph(&graph)
}

func Test_findNodeForEdge(t *testing.T) {
	// Basic lookup tested in testAddVertexNode, only test failures here
	var graph VertexGraph
	initVertexGraph(&graph, 10, 9)
	var node *VertexNode

	t.Run("empty graph", func(t *testing.T) {
		node = findNodeForEdge(&graph, &vertex1, &vertex2)
		require.True(t, node == nil, "Node lookup failed correctly for empty graph")
		addVertexNode(&graph, &vertex1, &vertex2)

		t.Run("different hash", func(t *testing.T) {
			node = findNodeForEdge(&graph, &vertex3, &vertex2)
			require.True(t, node == nil, "Node lookup failed correctly for different hash")

			t.Run("hash collision", func(t *testing.T) {
				node = findNodeForEdge(&graph, &vertex1, &vertex3)
				require.True(t, node == nil, "Node lookup failed correctly for hash collision")
				addVertexNode(&graph, &vertex1, &vertex4)

				t.Run("list iteration", func(t *testing.T) {
					node = findNodeForEdge(&graph, &vertex1, &vertex3)
					require.True(t, node == nil, "Node lookup failed correctly for collision w/iteration")
				})
			})

		})
	})

	destroyVertexGraph(&graph)
}

func Test_findNodeForVertex(t *testing.T) {
	var graph VertexGraph
	initVertexGraph(&graph, 10, 9)
	var node *VertexNode

	// Empty graph
	node = findNodeForVertex(&graph, &vertex1)
	require.True(t, node == nil, "Node lookup failed correctly for empty graph")
	addVertexNode(&graph, &vertex1, &vertex2)
	node = findNodeForVertex(&graph, &vertex1)
	require.True(t, node != nil, "Node lookup succeeded for correct node")
	node = findNodeForVertex(&graph, &vertex3)
	require.True(t, node == nil,
		"Node lookup failed correctly for different node")
	destroyVertexGraph(&graph)
}

func Test_removeVertexNode(t *testing.T) {
	var graph VertexGraph
	initVertexGraph(&graph, 10, 9)
	var node *VertexNode
	var success bool

	// Straight removal
	node = addVertexNode(&graph, &vertex1, &vertex2)
	success = removeVertexNode(&graph, node) == 0
	require.True(t, success, "Removal successful")
	require.True(t, findNodeForVertex(&graph, &vertex1) == nil,
		"Node lookup cannot find node")
	require.True(t, graph.size == 0, "Graph size decremented")

	// Remove end of list
	addVertexNode(&graph, &vertex1, &vertex2)
	node = addVertexNode(&graph, &vertex1, &vertex3)
	success = removeVertexNode(&graph, node) == 0
	require.True(t, success, "Removal successful")
	require.True(t, findNodeForEdge(&graph, &vertex1, &vertex3) == nil,
		"Node lookup cannot find node")
	require.True(t, findNodeForEdge(&graph, &vertex1, &vertex2).next == nil,
		"Base bucket node not pointing to node")
	require.True(t, graph.size == 1, "Graph size decremented")

	// This removal is just cleanup
	node = findNodeForVertex(&graph, &vertex1)
	require.True(t, removeVertexNode(&graph, node) == 0)

	// Remove beginning of list
	node = addVertexNode(&graph, &vertex1, &vertex2)
	addVertexNode(&graph, &vertex1, &vertex3)
	success = removeVertexNode(&graph, node) == 0
	require.True(t, success, "Removal successful")
	require.True(t, findNodeForEdge(&graph, &vertex1, &vertex2) == nil,
		"Node lookup cannot find node")
	require.True(t, findNodeForEdge(&graph, &vertex1, &vertex3) != nil,
		"Node lookup can find previous end of list")
	require.True(t, findNodeForEdge(&graph, &vertex1, &vertex3).next == nil,
		"Base bucket node not pointing to node")
	require.True(t, graph.size == 1, "Graph size decremented")

	// This removal is just cleanup
	node = findNodeForVertex(&graph, &vertex1)
	require.True(t, removeVertexNode(&graph, node) == 0)

	// Remove middle of list
	addVertexNode(&graph, &vertex1, &vertex2)
	node = addVertexNode(&graph, &vertex1, &vertex3)
	addVertexNode(&graph, &vertex1, &vertex4)
	success = removeVertexNode(&graph, node) == 0
	require.True(t, success, "Removal successful")
	require.True(t, findNodeForEdge(&graph, &vertex1, &vertex3) == nil,
		"Node lookup cannot find node")
	require.True(t, findNodeForEdge(&graph, &vertex1, &vertex4) != nil,
		"Node lookup can find previous end of list")
	require.True(t, graph.size == 2, "Graph size decremented")

	// Remove non-existent node
	node = &VertexNode{}
	success = removeVertexNode(&graph, node) == 0
	require.True(t, !success, "Removal of non-existent node fails")
	require.True(t, graph.size == 2, "Graph size unchanged")
	destroyVertexGraph(&graph)
}

func Test_firstVertexNode(t *testing.T) {
	var graph VertexGraph
	initVertexGraph(&graph, 10, 9)
	var node *VertexNode
	var addedNode *VertexNode
	node = firstVertexNode(&graph)
	require.True(t, node == nil, "No node found for empty graph")
	addedNode = addVertexNode(&graph, &vertex1, &vertex2)
	node = firstVertexNode(&graph)
	require.True(t, node == addedNode, "Node found")
	destroyVertexGraph(&graph)
}

func Test_destroyEmptyVertexGraph(t *testing.T) {
	var graph VertexGraph
	initVertexGraph(&graph, 10, 9)
	destroyVertexGraph(&graph)
}

func Test_singleBucketVertexGraph(t *testing.T) {
	var graph VertexGraph
	initVertexGraph(&graph, 1, 9)
	var node *VertexNode
	require.True(t, graph.numBuckets == 1, "1 bucket created")

	node = firstVertexNode(&graph)
	require.True(t, node == nil, "No node found for empty graph")

	node = addVertexNode(&graph, &vertex1, &vertex2)
	require.True(t, node != nil, "Node added")
	require.True(t, firstVertexNode(&graph) == node, "First node is node")

	addVertexNode(&graph, &vertex2, &vertex3)
	addVertexNode(&graph, &vertex3, &vertex4)
	require.True(t, firstVertexNode(&graph) == node, "First node is still node")
	require.True(t, graph.size == 3, "Graph size updated")

	destroyVertexGraph(&graph)
}
