package h3

/**
 * Returns whether or not the provided H3Indexes are neighbors.
 * @param origin The origin H3 index.
 * @param destination The destination H3 index.
 * @return 1 if the indexes are neighbors, 0 otherwise;
 */
func h3IndexesAreNeighbors(origin H3Index, destination H3Index) int {
	// Make sure they're hexagon indexes
	if H3_GET_MODE(origin) != H3_HEXAGON_MODE ||
		H3_GET_MODE(destination) != H3_HEXAGON_MODE {
		return 0
	}

	// Hexagons cannot be neighbors with themselves
	if origin == destination {
		return 0
	}

	// Only hexagons in the same resolution can be neighbors
	if H3_GET_RESOLUTION(origin) != H3_GET_RESOLUTION(destination) {
		return 0
	}

	// H3 Indexes that share the same parent are very likely to be neighbors
	// Child 0 is neighbor with all of its parent's 'offspring', the other
	// children are neighbors with 3 of the 7 children. So a simple comparison
	// of origin and destination parents and then a lookup table of the children
	// is a super-cheap way to possibly determine they are neighbors.
	parentRes := H3_GET_RESOLUTION(origin) - 1
	if parentRes > 0 && (h3ToParent(origin, parentRes) ==
		h3ToParent(destination, parentRes)) {
		originResDigit := H3_GET_INDEX_DIGIT(origin, parentRes+1)
		destinationResDigit := H3_GET_INDEX_DIGIT(destination, parentRes+1)
		if originResDigit == CENTER_DIGIT ||
			destinationResDigit == CENTER_DIGIT {
			return 1
		}
		// These sets are the relevant neighbors in the clockwise
		// and counter-clockwise
		neighborSetClockwise := []Direction{
			CENTER_DIGIT,
			JK_AXES_DIGIT,
			IJ_AXES_DIGIT,
			J_AXES_DIGIT,
			IK_AXES_DIGIT,
			K_AXES_DIGIT,
			I_AXES_DIGIT,
		}

		neighborSetCounterclockwise := []Direction{
			CENTER_DIGIT, IK_AXES_DIGIT, JK_AXES_DIGIT, K_AXES_DIGIT,
			IJ_AXES_DIGIT, I_AXES_DIGIT, J_AXES_DIGIT}

		if neighborSetClockwise[originResDigit] == destinationResDigit || neighborSetCounterclockwise[originResDigit] ==
			destinationResDigit {
			return 1
		}
	}

	// Otherwise, we have to determine the neighbor relationship the "hard" way.
	neighborRing := make([]H3Index, 7)
	kRing(origin, 1, neighborRing)
	for i := 0; i < 7; i++ {
		if neighborRing[i] == destination {
			return 1
		}
	}

	// Made it here, they definitely aren't neighbors
	return 0
}

/**
 * Returns a unidirectional edge H3 index based on the provided origin and
 * destination
 * @param origin The origin H3 hexagon index
 * @param destination The destination H3 hexagon index
 * @return The unidirectional edge H3Index, or 0 on failure.
 */
func getH3UnidirectionalEdge(origin H3Index, destination H3Index) H3Index {
	// Short-circuit and return an invalid index value if they are not neighbors
	if h3IndexesAreNeighbors(origin, destination) == 0 {
		return H3_INVALID_INDEX
	}

	// Otherwise, determine the IJK direction from the origin to the destination
	output := origin
	H3_SET_MODE(&output, H3_UNIEDGE_MODE)

	// Checks each neighbor, in order, to determine which direction the
	// destination neighbor is located. Skips CENTER_DIGIT since that
	// would be this index.
	var neighbor H3Index
	for direction := K_AXES_DIGIT; direction < NUM_DIGITS; direction++ {
		rotations := 0
		neighbor = h3NeighborRotations(origin, direction, &rotations)
		if neighbor == destination {
			H3_SET_RESERVED_BITS(&output, int(direction))
			return output
		}
	}

	// This should be impossible, return an invalid H3Index in this case;
	return H3_INVALID_INDEX // LCOV_EXCL_LINE
}

/**
 * Returns the origin hexagon from the unidirectional edge H3Index
 * @param edge The edge H3 index
 * @return The origin H3 hexagon index
 */
func getOriginH3IndexFromUnidirectionalEdge(edge H3Index) H3Index {
	if H3_GET_MODE(edge) != H3_UNIEDGE_MODE {
		return H3_INVALID_INDEX
	}
	origin := edge
	H3_SET_MODE(&origin, H3_HEXAGON_MODE)
	H3_SET_RESERVED_BITS(&origin, 0)
	return origin
}

/**
 * Returns the destination hexagon from the unidirectional edge H3Index
 * @param edge The edge H3 index
 * @return The destination H3 hexagon index
 */
func getDestinationH3IndexFromUnidirectionalEdge(edge H3Index) H3Index {
	if H3_GET_MODE(edge) != H3_UNIEDGE_MODE {
		return H3_INVALID_INDEX
	}
	direction := H3_GET_RESERVED_BITS(edge)
	rotations := 0
	destination := h3NeighborRotations(getOriginH3IndexFromUnidirectionalEdge(edge), Direction(direction), &rotations)
	return destination
}

/**
 * Determines if the provided H3Index is a valid unidirectional edge index
 * @param edge The unidirectional edge H3Index
 * @return 1 if it is a unidirectional edge H3Index, otherwise 0.
 */
func h3UnidirectionalEdgeIsValid(edge H3Index) bool {
	if H3_GET_MODE(edge) != H3_UNIEDGE_MODE {
		return false
	}

	neighborDirection := H3_GET_RESERVED_BITS(edge)
	if Direction(neighborDirection) <= CENTER_DIGIT || Direction(neighborDirection) >= NUM_DIGITS {
		return false
	}

	origin := getOriginH3IndexFromUnidirectionalEdge(edge)
	if h3IsPentagon(origin) && Direction(neighborDirection) == K_AXES_DIGIT {
		return false
	}

	return h3IsValid(origin)
}

/**
 * Returns the origin, destination pair of hexagon IDs for the given edge ID
 * @param edge The unidirectional edge H3Index
 * @param originDestination Pointer to memory to store origin and destination
 * IDs
 */
func getH3IndexesFromUnidirectionalEdge(edge H3Index, originDestination []H3Index) {
	originDestination[0] = getOriginH3IndexFromUnidirectionalEdge(edge)
	originDestination[1] = getDestinationH3IndexFromUnidirectionalEdge(edge)
}

/**
 * Provides all of the unidirectional edges from the current H3Index.
 * @param origin The origin hexagon H3Index to find edges for.
 * @param edges The memory to store all of the edges inside.
 */
func getH3UnidirectionalEdgesFromHexagon(origin H3Index, edges []H3Index) {
	// Determine if the origin is a pentagon and special treatment needed.
	isPentagon := h3IsPentagon(origin)

	// This is actually quite simple. Just modify the bits of the origin
	// slightly for each direction, except the 'k' direction in pentagons,
	// which is zeroed.
	for i := 0; i < 6; i++ {
		if isPentagon && i == 0 {
			edges[i] = H3_INVALID_INDEX
		} else {
			edges[i] = origin
			H3_SET_MODE(&edges[i], H3_UNIEDGE_MODE)
			H3_SET_RESERVED_BITS(&edges[i], i+1)
		}
	}
}

/**
 * Whether the given coordinate has a matching vertex in the given geo boundary.
 * @param  vertex   Coordinate to check
 * @param  boundary Geo boundary to look in
 * @return          Whether a match was found
 */
func _hasMatchingVertex(vertex *GeoCoord, boundary *GeoBoundary) bool {
	for i := 0; i < boundary.numVerts; i++ {
		if geoAlmostEqualThreshold(vertex, &boundary.Verts[i], 0.000001) {
			return true
		}
	}
	return false
}

/**
 * Provides the coordinates defining the unidirectional edge.
 * @param edge The unidirectional edge H3Index
 * @param gb The geoboundary object to store the edge coordinates.
 */
func getH3UnidirectionalEdgeBoundary(edge H3Index, gb *GeoBoundary) {
	// TODO: More efficient solution :)
	origin := GeoBoundary{}
	destination := GeoBoundary{}
	postponedVertex := GeoCoord{}
	hasPostponedVertex := false

	h3ToGeoBoundary(getOriginH3IndexFromUnidirectionalEdge(edge), &origin)
	h3ToGeoBoundary(getDestinationH3IndexFromUnidirectionalEdge(edge), &destination)

	k := 0
	for i := 0; i < origin.numVerts; i++ {
		if _hasMatchingVertex(&origin.Verts[i], &destination) {
			// If we are on vertex 0, we need to handle the case where it's the
			// end of the edge, not the beginning.
			if i == 0 && !_hasMatchingVertex(&origin.Verts[i+1], &destination) {
				postponedVertex = origin.Verts[i]
				hasPostponedVertex = true
			} else {
				gb.Verts[k] = origin.Verts[i]
				k++
			}
		}
	}
	// If we postponed adding the last vertex, add it now
	if hasPostponedVertex {
		gb.Verts[k] = postponedVertex
		k++
	}
	gb.numVerts = k
}
