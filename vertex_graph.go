package h3

import (
	"math"
)

type VertexNode struct {
	from GeoCoord
	to   GeoCoord
	next *VertexNode
}

type VertexGraph struct {
	buckets    []*VertexNode
	numBuckets int
	size       int
	res        int
}

/**
 * Initialize a new VertexGraph
 * @param graph       Graph to initialize
 * @param  numBuckets Number of buckets to include in the graph
 * @param  res        Resolution of the hexagons whose vertices we're storing
 */
func initVertexGraph(graph *VertexGraph, numBuckets int, res int) {
	if numBuckets > 0 {
		graph.buckets = make([]*VertexNode, numBuckets)
	} else {
		graph.buckets = nil
	}

	graph.numBuckets = numBuckets
	graph.size = 0
	graph.res = res
}

/**
 * Destroy a VertexGraph's sub-objects, freeing their memory. The caller is
 * responsible for freeing memory allocated to the VertexGraph struct itself.
 * @param graph Graph to destroy
 */
func destroyVertexGraph(graph *VertexGraph) {
	node := firstVertexNode(graph);
	for node != nil {
		removeVertexNode(graph, node)
		node = firstVertexNode(graph)
	}
	graph.buckets = nil
}

/**
 * Get an integer hash for a Lat/Lon point, at a precision determined
 * by the current hexagon resolution.
 * TODO: Light testing suggests this might not be sufficient at resolutions
 * finer than 10. Design a better hash function if performance and collisions
 * seem to be an issue here.
 * @param  vertex     Lat/Lon vertex to hash
 * @param  res        Resolution of the hexagon the vertex belongs to
 * @param  numBuckets Number of buckets in the graph
 * @return            Integer hash
 */
func _hashVertex(vertex *GeoCoord, res int, numBuckets int) uint32 {
	// Simple hash: Take the sum of the Lat and Lon with a precision level
	// determined by the resolution, converted to int, modulo bucket count.
	return uint32(math.Mod(math.Abs((vertex.Lat+vertex.Lon)*math.Pow(10, float64(15-res))), float64(numBuckets)))
}

func _initVertexNode(node *VertexNode, fromVtx *GeoCoord, toVtx *GeoCoord) {
	node.from = *fromVtx
	node.to = *toVtx
	node.next = nil
}

/**
 * Add a edge to the graph
 * @param graph   Graph to add node to
 * @param fromVtx Start vertex
 * @param toVtx   End vertex
 * @return        Pointer to the new node
 */
func addVertexNode(graph *VertexGraph, fromVtx *GeoCoord, toVtx *GeoCoord) *VertexNode {
	// Make the new node
	node := &VertexNode{}

	_initVertexNode(node, fromVtx, toVtx)
	// Determine location
	index := _hashVertex(fromVtx, graph.res, graph.numBuckets)
	// Check whether there's an existing node in that spot
	currentNode := graph.buckets[index]

	if currentNode == nil {
		// Set bucket to the new node
		graph.buckets[index] = node
	} else {
		// Find the end of the list
		for {
			// Check the the edge we're adding doesn't already exist
			if geoAlmostEqual(&currentNode.from, fromVtx) && geoAlmostEqual(&currentNode.to, toVtx) {
				node = nil
				// already exists, bail
				return currentNode
			}

			if currentNode.next != nil {
				currentNode = currentNode.next
			}

			if currentNode.next == nil {
				break
			}
		}
		// Add the new node to the end of the list
		currentNode.next = node
	}

	graph.size++
	return node
}

/**
 * Remove a node from the graph. The input node will be freed, and should
 * not be used after removal.
 * @param graph Graph to mutate
 * @param node  Node to remove
 * @return      0 on success, 1 on failure (node not found)
 */
func removeVertexNode(graph *VertexGraph, node *VertexNode) int {
	// Determine location
	index := _hashVertex(&node.from, graph.res, graph.numBuckets)
	currentNode := graph.buckets[index]
	found := 0

	if currentNode != nil {
		if currentNode == node {
			graph.buckets[index] = node.next
			found = 1
		}
		// Look through the list
		for found != 1 && currentNode.next != nil {
			if currentNode.next == node {
				// splice the node out
				currentNode.next = node.next
				found = 1
			}
			currentNode = currentNode.next
		}
	}

	if found > 0 {
		node = nil
		graph.size--
		return 0
	}

	return 1
}

/**
 * Find the Vertex node for a given edge, if it exists
 * @param  graph   Graph to look in
 * @param  fromVtx Start vertex
 * @param  toVtx   End vertex, or nil if we don't care
 * @return         Pointer to the vertex node, if found
 */
func findNodeForEdge(graph *VertexGraph, fromVtx *GeoCoord, toVtx *GeoCoord) *VertexNode {
	// Determine location
	index := _hashVertex(fromVtx, graph.res, graph.numBuckets)
	// Check whether there's an existing node in that spot
	node := graph.buckets[index]

	// Look through the list and see if we find the edge
	for node != nil {
		if geoAlmostEqual(&node.from, fromVtx) && (toVtx == nil || geoAlmostEqual(&node.to, toVtx)) {
			return node
		}
		node = node.next
	}
	// Iteration lookup fail
	return nil
}

/**
 * Find a Vertex node starting at the given vertex
 * @param  graph   Graph to look in
 * @param  fromVtx Start vertex
 * @return         Pointer to the vertex node, if found
 */
func findNodeForVertex(graph *VertexGraph, fromVtx *GeoCoord) *VertexNode {
	return findNodeForEdge(graph, fromVtx, nil)
}

/**
 * Get the next vertex node in the graph.
 * @param  graph Graph to iterate
 * @return       Vertex node, or nil if at the end
 */
func firstVertexNode(graph *VertexGraph) *VertexNode {
	var node *VertexNode

	currentIndex := 0

	for node == nil {
		if currentIndex < graph.numBuckets {
			// find the first node in the next bucket
			node = graph.buckets[currentIndex]
		} else {
			// end of iteration
			return nil
		}
		currentIndex++
	}

	return node
}
