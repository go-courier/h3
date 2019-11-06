package h3

const HEX_RANGE_SUCCESS = 0
const HEX_RANGE_PENTAGON = 1
const HEX_RANGE_K_SUBSEQUENCE = 2
const MAX_ONE_RING_SIZE = 7

/**
 * Directions used for traversing a hexagonal ring counterclockwise around
 * {1, 0, 0}
 *
 * <pre>
 *      _
 *    _/ \\_
 *   / \\5/ \\
 *   \\0/ \\4/
 *   / \\_/ \\
 *   \\1/ \\3/
 *     \\2/
 * </pre>
 */
var DIRECTIONS = [6]Direction{J_AXES_DIGIT, JK_AXES_DIGIT, K_AXES_DIGIT, IK_AXES_DIGIT, I_AXES_DIGIT, IJ_AXES_DIGIT}

/**
 * Direction used for traversing to the next outward hexagonal ring.
 */
const NEXT_RING_DIRECTION = I_AXES_DIGIT

/**
 * New digit when traversing along class II grids.
 *
 * Current digit . direction . new digit.
 */
var NEW_DIGIT_II = [7][7]Direction{
	{CENTER_DIGIT, K_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT, IJ_AXES_DIGIT},
	{K_AXES_DIGIT, I_AXES_DIGIT, JK_AXES_DIGIT, IJ_AXES_DIGIT, IK_AXES_DIGIT, J_AXES_DIGIT, CENTER_DIGIT},
	{J_AXES_DIGIT, JK_AXES_DIGIT, K_AXES_DIGIT, I_AXES_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT, IK_AXES_DIGIT},
	{JK_AXES_DIGIT, IJ_AXES_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT, J_AXES_DIGIT},
	{I_AXES_DIGIT, IK_AXES_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT, K_AXES_DIGIT},
	{IK_AXES_DIGIT, J_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT, JK_AXES_DIGIT, IJ_AXES_DIGIT, I_AXES_DIGIT},
	{IJ_AXES_DIGIT, CENTER_DIGIT, IK_AXES_DIGIT, J_AXES_DIGIT, K_AXES_DIGIT, I_AXES_DIGIT, JK_AXES_DIGIT},
}

/**
 * New traversal direction when traversing along class II grids.
 *
 * Current digit . direction . new ap7 move (at coarser level).
 */
var NEW_ADJUSTMENT_II = [7][7]Direction{
	{CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, K_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT, CENTER_DIGIT, IK_AXES_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, CENTER_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT, CENTER_DIGIT, CENTER_DIGIT, J_AXES_DIGIT},
	{CENTER_DIGIT, K_AXES_DIGIT, JK_AXES_DIGIT, JK_AXES_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, I_AXES_DIGIT, I_AXES_DIGIT, IJ_AXES_DIGIT},
	{CENTER_DIGIT, IK_AXES_DIGIT, CENTER_DIGIT, CENTER_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, CENTER_DIGIT, J_AXES_DIGIT, CENTER_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT, IJ_AXES_DIGIT},
}

/**
 * New traversal direction when traversing along class III grids.
 *
 * Current digit . direction . new ap7 move (at coarser level).
 */
var NEW_DIGIT_III = [7][7]Direction{
	{CENTER_DIGIT, K_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT, IJ_AXES_DIGIT},
	{K_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT},
	{J_AXES_DIGIT, JK_AXES_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT},
	{JK_AXES_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT, J_AXES_DIGIT},
	{I_AXES_DIGIT, IK_AXES_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT},
	{IK_AXES_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT, I_AXES_DIGIT},
	{IJ_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT},
}

/**
 * New traversal direction when traversing along class III grids.
 *
 * Current digit . direction . new ap7 move (at coarser level).
 */
var NEW_ADJUSTMENT_III = [7][7]Direction{
	{CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, K_AXES_DIGIT, CENTER_DIGIT, JK_AXES_DIGIT, CENTER_DIGIT, K_AXES_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, CENTER_DIGIT, J_AXES_DIGIT, J_AXES_DIGIT, CENTER_DIGIT, CENTER_DIGIT, IJ_AXES_DIGIT},
	{CENTER_DIGIT, JK_AXES_DIGIT, J_AXES_DIGIT, JK_AXES_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, CENTER_DIGIT, I_AXES_DIGIT, IK_AXES_DIGIT, I_AXES_DIGIT},
	{CENTER_DIGIT, K_AXES_DIGIT, CENTER_DIGIT, CENTER_DIGIT, IK_AXES_DIGIT, IK_AXES_DIGIT, CENTER_DIGIT},
	{CENTER_DIGIT, CENTER_DIGIT, IJ_AXES_DIGIT, CENTER_DIGIT, I_AXES_DIGIT, CENTER_DIGIT, IJ_AXES_DIGIT},
}

/**
 * Maximum number of indices that result from the kRing algorithm with the given
 * k. Formula source and proof: https://oeis.org/A003215
 *
 * @param k k value, k >= 0.
 */
func maxKringSize(k int) int { return 3*k*(k+1) + 1 }

/**
 * k-rings produces indices within k distance of the origin index.
 *
 * k-ring 0 is defined as the origin index, k-ring 1 is defined as k-ring 0 and
 * all neighboring indices, and so on.
 *
 * Output is placed in the provided array in no particular order. Elements of
 * the output array may be left zero, as can happen when crossing a pentagon.
 *
 * @param origin Origin location.
 * @param k k >= 0
 * @param out Zero-filled array which must be of size maxKringSize(k).
 */
func kRing(origin H3Index, k int, out []H3Index) {
	maxIdx := maxKringSize(k)
	distances := make([]int, maxIdx)
	kRingDistances(origin, k, out, distances)
}

/**
 * k-rings produces indices within k distance of the origin index.
 *
 * k-ring 0 is defined as the origin index, k-ring 1 is defined as k-ring 0 and
 * all neighboring indices, and so on.
 *
 * Output is placed in the provided array in no particular order. Elements of
 * the output array may be left zero, as can happen when crossing a pentagon.
 *
 * @param origin Origin location.
 * @param k k >= 0
 * @param out Zero-filled array which must be of size maxKringSize(k).
 * @param distances Zero-filled array which must be of size maxKringSize(k).
 */
func kRingDistances(origin H3Index, k int, out []H3Index, distances []int) {
	maxIdx := maxKringSize(k)
	// Optimistically try the faster hexRange algorithm first
	failed := hexRangeDistances(origin, k, out, distances)
	if failed != 0 {
		// Fast algo failed, fall back to slower, correct algo
		// and also wipe out array because contents untrustworthy
		out = make([]H3Index, maxIdx) // todo

		_kRingInternal(origin, k, out, distances, maxIdx, 0)
	}
}

/**
 * Internal helper function called recursively for kRingDistances.
 *
 * Adds the origin index to the output set (treating it as a hash set)
 * and recurses to its neighbors, if needed.
 *
 * @param origin
 * @param k Maximum distance to move from the origin.
 * @param out Array treated as a hash set, elements being either or H3Index 0.
 * @param distances Scratch area, with elements paralleling the out array.
 * Elements indicate ijk distance from the origin index to the output index.
 * @param maxIdx Size of out and scratch arrays (must be maxKringSize(k))
 * @param curK Current distance from the origin.
 */
func _kRingInternal(origin H3Index, k int, out []H3Index, distances []int, maxIdx int, curK int) {
	if origin == 0 {
		return
	}

	// Put origin in the output array. out is used as a hash set.
	off := uint64(origin) % uint64(maxIdx)
	for out[off] != 0 && out[off] != origin {
		off = (off + 1) % uint64(maxIdx)
	}

	// We either got a free slot in the hash set or hit a duplicate
	// We might need to process the duplicate anyways because we got
	// here on a longer path before.
	if out[off] == origin && distances[off] <= curK {
		return
	}
	out[off] = origin
	distances[off] = curK

	// Base case: reached an index k away from the origin.
	if curK >= k {
		return
	}

	// Recurse to all neighbors in no particular order.
	for i := 0; i < 6; i++ {
		rotations := 0
		_kRingInternal(h3NeighborRotations(origin, DIRECTIONS[i], &rotations), k, out, distances, maxIdx, curK+1)
	}
}

/**
 * Returns the hexagon index neighboring the origin, in the direction dir.
 *
 * Implementation note: The only reachable case where this returns 0 is if the
 * origin is a pentagon and the translation is in the k direction. Thus,
 * 0 can only be returned if origin is a pentagon.
 *
 * @param origin Origin index
 * @param dir Direction to move in
 * @param rotations Number of ccw rotations to perform to reorient the
 *                  translation vector. Will be modified to the new number of
 *                  rotations to perform (such as when crossing a face edge.)
 * @return of H3Index the specified neighbor or 0 if deleted k-subsequence
 *         distortion is encountered.
 */
func h3NeighborRotations(origin H3Index, dir Direction, rotations *int) H3Index {
	out := origin
	for i := 0; i < *rotations; i++ {
		dir = _rotate60ccw(dir)
	}

	newRotations := 0
	oldBaseCell := H3_GET_BASE_CELL(out)
	oldLeadingDigit := _h3LeadingNonZeroDigit(out)

	// Adjust the indexing digits and, if needed, the base cell.
	r := H3_GET_RESOLUTION(out) - 1
	for {
		if r == -1 {
			H3_SET_BASE_CELL(&out, baseCellNeighbors[oldBaseCell][dir])
			newRotations = baseCellNeighbor60CCWRots[oldBaseCell][dir]
			if H3_GET_BASE_CELL(out) == INVALID_BASE_CELL {
				// Adjust for the deleted k vertex at the base cell level.
				// This edge actually borders a different neighbor.
				H3_SET_BASE_CELL(&out, baseCellNeighbors[oldBaseCell][IK_AXES_DIGIT])
				newRotations =
					baseCellNeighbor60CCWRots[oldBaseCell][IK_AXES_DIGIT]

				// perform the adjustment for the k-subsequence we're skipping
				// over.
				out = _h3Rotate60ccw(out)
				*rotations = *rotations + 1
			}

			break
		} else {
			oldDigit := H3_GET_INDEX_DIGIT(out, r+1)
			var nextDir Direction
			if isResClassIII(r + 1) {
				H3_SET_INDEX_DIGIT(&out, r+1, NEW_DIGIT_II[oldDigit][dir])
				nextDir = NEW_ADJUSTMENT_II[oldDigit][dir]
			} else {
				H3_SET_INDEX_DIGIT(&out, r+1, NEW_DIGIT_III[oldDigit][dir])
				nextDir = NEW_ADJUSTMENT_III[oldDigit][dir]
			}

			if nextDir != CENTER_DIGIT {
				dir = nextDir
				r--
			} else {
				// No more adjustment to perform
				break
			}
		}
	}

	newBaseCell := H3_GET_BASE_CELL(out)
	if _isBaseCellPentagon(newBaseCell) {
		alreadyAdjustedKSubsequence := 0

		// force rotation out of missing k-axes sub-sequence
		if _h3LeadingNonZeroDigit(out) == K_AXES_DIGIT {
			if oldBaseCell != newBaseCell {
				// in this case, we traversed into the deleted
				// k subsequence of a pentagon base cell.
				// We need to rotate out of that case depending
				// on how we got here.
				// check for a cw/ccw var face offset; default is ccw

				if _baseCellIsCwOffset(
					newBaseCell, baseCellData[oldBaseCell].homeFijk.face) {
					out = _h3Rotate60cw(out)
				} else {
					// See cwOffsetPent in testKRing.c for why this is
					// unreachable.
					out = _h3Rotate60ccw(out) // LCOV_EXCL_LINE
				}
				alreadyAdjustedKSubsequence = 1
			} else {
				// In this case, we traversed into the deleted
				// k subsequence from within the same pentagon
				// base cell.
				if oldLeadingDigit == CENTER_DIGIT {
					// Undefined: the k direction is deleted from here
					return H3_INVALID_INDEX
				} else if oldLeadingDigit == JK_AXES_DIGIT {
					// Rotate out of the deleted k subsequence
					// We also need an additional change to the direction we're
					// moving in
					out = _h3Rotate60ccw(out)
					*rotations = *rotations + 1
				} else if oldLeadingDigit == IK_AXES_DIGIT {
					// Rotate out of the deleted k subsequence
					// We also need an additional change to the direction we're
					// moving in
					out = _h3Rotate60cw(out)
					*rotations = *rotations + 5
				} else {
					// Should never occur
					return H3_INVALID_INDEX // LCOV_EXCL_LINE
				}
			}
		}

		for i := 0; i < newRotations; i++ {
			out = _h3RotatePent60ccw(out)
		}

		// Account for differing orientation of the base cells (this edge
		// might not follow properties of some other edges.)
		if oldBaseCell != newBaseCell {
			if _isBaseCellPolarPentagon(newBaseCell) {
				// 'polar' base cells behave differently because they have all
				// i neighbors.
				if oldBaseCell != 118 && oldBaseCell != 8 &&
					_h3LeadingNonZeroDigit(out) != JK_AXES_DIGIT {
					*rotations = *rotations + 1
				}
			} else if _h3LeadingNonZeroDigit(out) == IK_AXES_DIGIT && alreadyAdjustedKSubsequence == 0 {
				// account for distortion introduced to the 5 neighbor by the
				// deleted k subsequence.
				*rotations = *rotations + 1
			}
		}
	} else {
		for i := 0; i < newRotations; i++ {
			out = _h3Rotate60ccw(out)
		}
	}

	*rotations = (*rotations + newRotations) % 6
	return out
}

/**
 * hexRange produces indexes within k distance of the origin index.
 * Output behavior is undefined when one of the indexes returned by this
 * function is a pentagon or is in the pentagon distortion area.
 *
 * k-ring 0 is defined as the origin index, k-ring 1 is defined as k-ring 0 and
 * all neighboring indexes, and so on.
 *
 * Output is placed in the provided array in order of increasing distance from
 * the origin.
 *
 * @param origin Origin location.
 * @param k k >= 0
 * @param out Array which must be of size maxKringSize(k).
 * @return 0 if no pentagon or pentagonal distortion area was encountered.
 */
func hexRange(origin H3Index, k int, out []H3Index) int {
	return hexRangeDistances(origin, k, out, []int{})
}

/**
 * hexRange produces indexes within k distance of the origin index.
 * Output behavior is undefined when one of the indexes returned by this
 * function is a pentagon or is in the pentagon distortion area.
 *
 * k-ring 0 is defined as the origin index, k-ring 1 is defined as k-ring 0 and
 * all neighboring indexes, and so on.
 *
 * Output is placed in the provided array in order of increasing distance from
 * the origin. The distances in hexagons is placed in the distances array at
 * the same offset.
 *
 * @param origin Origin location.
 * @param k k >= 0
 * @param out Array which must be of size maxKringSize(k).
 * @param distances Null or array which must be of size maxKringSize(k).
 * @return 0 if no pentagon or pentagonal distortion area was encountered.
 */
func hexRangeDistances(origin H3Index, k int, out []H3Index, distances []int) int {
	// Return codes:
	// 1 Pentagon was encountered
	// 2 Pentagon distortion (deleted k subsequence) was encountered
	// Pentagon being encountered is not itself var problem a; really the deleted
	// k-subsequence is the problem, but for compatibility reasons we fail on
	// the pentagon.

	// k must be >= 0, so origin is always needed
	idx := 0
	out[idx] = origin
	if len(distances) > idx {
		distances[idx] = 0
	}
	idx++
	if h3IsPentagon(origin) {
		// Pentagon var encountered was; bail out as user doesn't want this.
		return HEX_RANGE_PENTAGON
	}

	// 0 < ring <= k, current ring
	ring := 1
	// 0 <= direction < 6, current side of the ring
	direction := 0
	// 0 <= i < ring, current position on the side of the ring
	i := 0
	// Number of 60 degree ccw rotations to perform on the direction (based on
	// which faces have been crossed.)
	rotations := 0
	for ring <= k {
		if direction == 0 && i == 0 {
			// Not putting in the output set as it will be done later, at
			// the end of this ring.
			origin =
				h3NeighborRotations(origin, NEXT_RING_DIRECTION, &rotations)
			if origin == 0 { // LCOV_EXCL_BR_LINE
				// Should not be possible because `origin` would have to be a
				// pentagon
				return HEX_RANGE_K_SUBSEQUENCE // LCOV_EXCL_LINE
			}

			if h3IsPentagon(origin) {
				// Pentagon var encountered was; bail out as user doesn't want this.
				return HEX_RANGE_PENTAGON
			}
		}

		origin = h3NeighborRotations(origin, DIRECTIONS[direction], &rotations)
		if origin == 0 { // LCOV_EXCL_BR_LINE
			// Should not be possible because `origin` would have to be a
			// pentagon
			return HEX_RANGE_K_SUBSEQUENCE // LCOV_EXCL_LINE
		}
		out[idx] = origin
		if len(distances) > idx {
			distances[idx] = ring
		}
		idx++
		i++
		// Check if end of this side of the k-ring
		if i == ring {
			i = 0
			direction++
			// Check if end of this ring.
			if direction == 6 {
				direction = 0
				ring++
			}
		}

		if h3IsPentagon(origin) {
			// Pentagon var encountered was; bail out as user doesn't want this.
			return HEX_RANGE_PENTAGON
		}
	}
	return HEX_RANGE_SUCCESS
}

/**
 * hexRanges takes an array of input hex IDs and a max k-ring and returns an
 * array of hexagon IDs sorted first by the original hex IDs and then by the
 * k-ring (0 to max), with no guaranteed sorting within each k-ring group.
 *
 * @param h3Set A pointer to an array of H3Indexes
 * @param length The total number of H3Indexes in h3Set
 * @param k The number of rings to generate
 * @param out A pointer to the output memory to dump the new set of H3Indexes to
 *            The memory block should be equal to maxKringSize(k) * length
 * @return 0 if no pentagon is encountered. Cannot trust output otherwise
 */
func hexRanges(h3Set []H3Index, length int, k int, out []H3Index) int {
	success := 0
	var segment []H3Index
	segmentSize := maxKringSize(k)
	for i := 0; i < length; i++ {
		// Determine the appropriate segment of the output array to operate on

		//segment = out + i*segmentSize;
		segment = append(out, make([]H3Index, i*segmentSize)...)

		success = hexRange(h3Set[i], k, segment)
		if success != 0 {
			return success
		}
	}
	return 0
}

/**
 * Returns the "hollow" ring of hexagons at exactly grid distance k from
 * the origin hexagon. In particular, k=0 returns just the origin hexagon.
 *
 * A nonzero failure code may be returned in some cases, for example,
 * if a pentagon is encountered.
 * Failure cases may be fixed in future versions.
 *
 * @param origin Origin location.
 * @param k k >= 0
 * @param out Array which must be of size 6 * k (or 1 if k == 0)
 * @return 0 var successful if; nonzero otherwise.
 */
func hexRing(origin H3Index, k int, out []H3Index) int {
	// Short-circuit on 'identity' ring
	if k == 0 {
		out[0] = origin
		return 0
	}
	idx := 0
	// Number of 60 degree ccw rotations to perform on the direction (based on
	// which faces have been crossed.)
	rotations := 0
	// Scratch structure for checking for pentagons
	if h3IsPentagon(origin) {
		// Pentagon var encountered was; bail out as user doesn't want this.
		return HEX_RANGE_PENTAGON
	}

	for ring := 0; ring < k; ring++ {
		origin = h3NeighborRotations(origin, NEXT_RING_DIRECTION, &rotations)
		if origin == 0 { // LCOV_EXCL_BR_LINE
			// Should not be possible because `origin` would have to be a
			// pentagon
			return HEX_RANGE_K_SUBSEQUENCE // LCOV_EXCL_LINE
		}

		if h3IsPentagon(origin) {
			return HEX_RANGE_PENTAGON
		}
	}

	lastIndex := origin
	out[idx] = origin
	idx++
	for direction := 0; direction < 6; direction++ {
		for pos := 0; pos < k; pos++ {
			origin = h3NeighborRotations(origin, DIRECTIONS[direction], &rotations)
			if origin == 0 { // LCOV_EXCL_BR_LINE
				// Should not be possible because `origin` would have to be a
				// pentagon
				return HEX_RANGE_K_SUBSEQUENCE // LCOV_EXCL_LINE
			}

			// Skip the very last index, it was already added. We do
			// however need to traverse to it because of the pentagonal
			// distortion check, below.
			if pos != k-1 || direction != 5 {
				out[idx] = origin
				idx++
				if h3IsPentagon(origin) {
					return HEX_RANGE_PENTAGON
				}
			}
		}
	}

	// Check that this matches the expected lastIndex, if it doesn't,
	// it indicates pentagonal distortion occurred and we should report
	// failure.
	if lastIndex != origin {
		return HEX_RANGE_PENTAGON
	}
	return HEX_RANGE_SUCCESS
}

/**
 * maxPolyfillSize returns the number of hexagons to allocate space for when
 * performing a polyfill on the given GeoJSON-like data structure.
 *
 * The size is the maximum of either the number of points in the geofence or the
 * number of hexagons in the bounding box of the geofence.
 *
 * @param geoPolygon A GeoJSON-like data structure indicating the poly to fill
 * @param res Hexagon resolution (0-15)
 * @return number of hexagons to allocate for
 */
func maxPolyfillSize(geoPolygon *GeoPolygon, res int) int {
	// Get the bounding box for the GeoJSON-like struct
	var bbox BBox
	geofence := geoPolygon.geofence
	bboxFrom(&geofence, &bbox)
	numHexagons := bboxHexEstimate(&bbox, res)
	// This algorithm assumes that the number of vertices is usually less than
	// the number of hexagons, but when it's wrong, this will keep it from
	// failing
	totalVerts := geofence.numVerts
	for i := 0; i < geoPolygon.numHoles; i++ {
		totalVerts += geoPolygon.holes[i].numVerts
	}
	if numHexagons < totalVerts {
		numHexagons = totalVerts
	}
	return numHexagons
}

/**
 * polyfill takes a given GeoJSON-like data structure and preallocated,
 * zeroed memory, and fills it with the hexagons that are contained by
 * the GeoJSON-like data structure.
 *
 * This implementation traces the GeoJSON geofence(s) in cartesian space with
 * hexagons, tests them and their neighbors to be contained by the geofence(s),
 * and then any newly found hexagons are used to test again until no new
 * hexagons are found.
 *
 * @param geoPolygon The geofence and holes defining the relevant area
 * @param res The Hexagon resolution (0-15)
 * @param out The slab of zeroed memory to write to. Assumed to be big enough.
 */
func polyfill(geoPolygon *GeoPolygon, res int, out []H3Index) {
	// TODO: Eliminate this wrapper with the H3 4.0.0 release
	failure := _polyfillInternal(geoPolygon, res, out)
	// The polyfill algorithm can theoretically fail if the allocated memory is
	// not large enough for the polygon, but this should be impossible given the
	// conservative overestimation of the number of hexagons possible.
	// LCOV_EXCL_START
	if failure != 0 {
		numHexagons := maxPolyfillSize(geoPolygon, res)
		for i := 0; i < numHexagons; i++ {
			out[i] = H3_INVALID_INDEX
		}
	}
	// LCOV_EXCL_STOP
}

/**
 * _getEdgeHexagons takes a given geofence ring (either the main geofence or
 * one of the holes) and traces it with hexagons and updates the search and
 * found memory blocks. This is used for determining the initial hexagon set
 * for the polyfill algorithm to execute on.
 *
 * @param geofence The geofence (or hole) to be traced
 * @param numHexagons The maximum number of hexagons possible for the geofence
 *                    (also the bounds of the search and found arrays)
 * @param res The hexagon resolution (0-15)
 * @param numSearchHexes The number of hexagons found so far to be searched
 * @param search The block of memory containing the hexagons to search from
 * @param found The block of memory containing the hexagons found from the
 * search
 *
 * @return An error code if the hash function cannot insert a found hexagon
 *         into the found array.
 */
func _getEdgeHexagons(geofence *Geofence, numHexagons int, res int, numSearchHexes *int, search []H3Index, found []H3Index) int {
	for i := 0; i < geofence.numVerts; i++ {
		origin := geofence.verts[i]
		var destination GeoCoord

		if i == geofence.numVerts-1 {
			destination = geofence.verts[0]
		} else {
			destination = geofence.verts[i+1]
		}

		numHexesEstimate := lineHexEstimate(&origin, &destination, res)
		for j := 0; j < numHexesEstimate; j++ {
			var interpolate GeoCoord
			interpolate.Lat = (origin.Lat * float64(numHexesEstimate-j) / float64(numHexesEstimate)) + (destination.Lat * float64(j) / float64(numHexesEstimate))
			interpolate.Lon = (origin.Lon * float64(numHexesEstimate-j) / float64(numHexesEstimate)) + (destination.Lon * float64(j) / float64(numHexesEstimate))
			pointHex := geoToH3(&interpolate, res)
			// A simple hash to store the hexagon, or move to another place if
			// needed
			loc := (int)(uint64(pointHex) % uint64(numHexagons))
			loopCount := 0
			for found[loc] != 0 {
				// If this conditional is reached, the `found` memory block is
				// too small for the given polygon. This should not happen.
				if loopCount > numHexagons {
					return -1
				} // LCOV_EXCL_LINE
				if found[loc] == pointHex {
					break // At least two points of the geofence index to the
				}
				// same cell
				loc = (loc + 1) % numHexagons
				loopCount++
			}
			if found[loc] == pointHex {
				continue // Skip this hex, already exists in the found hash
			}
			// Otherwise, set it in the found hash for now
			found[loc] = pointHex
			search[*numSearchHexes] = pointHex
			(*numSearchHexes)++
		}
	}
	return 0
}

/**
 * _polyfillInternal traces the provided geoPolygon data structure with hexagons
 * and then iteratively searches through these hexagons and their immediate
 * neighbors to see if they are contained within the polygon or not. Those that
 * are found are added to the out array as well as the found array. Once all
 * hexagons to search are checked, the found hexagons become the new search
 * array and the found array is wiped and the process repeats until no new
 * hexagons can be found.
 *
 * @param geoPolygon The geofence and holes defining the relevant area
 * @param res The Hexagon resolution (0-15)
 * @param out The slab of zeroed memory to write to. Assumed to be big enough.
 *
 * @return An error code if any of the hash operations fails to insert a hexagon
 *         into an array of memory.
 */
func _polyfillInternal(geoPolygon *GeoPolygon, res int, out []H3Index) int {
	// One of the goals of the polyfill algorithm is that two adjacent polygons
	// with zero overlap have zero overlapping hexagons. That the hexagons are
	// uniquely assigned. There are a few approaches to take here, such as
	// deciding based on which polygon has the greatest overlapping area of the
	// hexagon, or the most number of contained points on the hexagon (using the
	// center poas int a tiebreaker).
	//
	// But if the polygons are convex, both of these more complex algorithms can
	// be reduced down to checking whether or not the center of the hexagon is
	// contained in the polygon, and so this is the approach that this polyfill
	// algorithm will follow, as it's simpler, faster, and the error for concave
	// polygons is still minimal (only affecting concave shapes on the order of
	// magnitude of the hexagon size or smaller, not impacting larger concave
	// shapes)
	//
	// This first part is identical to the maxPolyfillSize above.

	// Get the bounding boxes for the polygon and any holes
	bboxes := make([]BBox, geoPolygon.numHoles+1)
	bboxesFromGeoPolygon(geoPolygon, bboxes)

	// Get the estimated number of hexagons and allocate some temporary memory
	// for the hexagons
	numHexagons := maxPolyfillSize(geoPolygon, res)
	search := make([]H3Index, numHexagons)
	found := make([]H3Index, numHexagons)

	// Some metadata for tracking the state of the search and found memory
	// blocks
	numSearchHexes := 0
	numFoundHexes := 0

	// 1. Trace the hexagons along the polygon defining the outer geofence and
	// add them to the search hash. The hexagon containing the geofence point
	// may or may not be contained by the geofence (as the hexagon's center
	// pomay int be outside of the boundary.)
	geofence := geoPolygon.geofence
	failure := _getEdgeHexagons(&geofence, numHexagons, res, &numSearchHexes,
		search, found)
	// If this branch is reached, we have exceeded the maximum number of
	// hexagons possible and need to clean up the allocated memory.
	// LCOV_EXCL_START
	if failure != 0 {
		search = nil
		found = nil
		bboxes = nil

		return failure
	}
	// LCOV_EXCL_STOP

	// 2. Iterate over all holes, trace the polygons defining the holes with
	// hexagons and add to only the search hash. We're going to temporarily use
	// the `found` hash to use for dedupe purposes and then re-zero it once
	// we're done here, otherwise we'd have to scan the whole set on each insert
	// to make sure there's no duplicates, which is very inefficient.
	for i := 0; i < geoPolygon.numHoles; i++ {
		hole := &(geoPolygon.holes[i])
		failure = _getEdgeHexagons(hole, numHexagons, res, &numSearchHexes,
			search, found)
		// If this branch is reached, we have exceeded the maximum number of
		// hexagons possible and need to clean up the allocated memory.
		// LCOV_EXCL_START
		if failure != 0 {
			search = nil
			found = nil
			bboxes = nil

			return failure
		}
		// LCOV_EXCL_STOP
	}

	// 3. Re-zero the found hash so it can be used in the main loop below
	for i := 0; i < numHexagons; i++ {
		found[i] = 0
	}

	// 4. Begin main loop. While the search hash is not empty do the following
	for numSearchHexes > 0 {
		// Iterate through all hexagons in the current search hash, then loop
		// through all neighbors and test Point-in-Poly, if point-in-poly
		// succeeds, add to out and found hashes if not already there.
		currentSearchNum := 0
		i := 0
		for currentSearchNum < numSearchHexes {
			ring := make([]H3Index, MAX_ONE_RING_SIZE)
			searchHex := search[i]
			kRing(searchHex, 1, ring)
			for j := 0; j < MAX_ONE_RING_SIZE; j++ {
				if ring[j] == H3_INVALID_INDEX {
					continue // Skip if this was a pentagon and only had 5
					// neighbors
				}

				hex := ring[j]

				// A simple hash to store the hexagon, or move to another place
				// if needed. This MUST be done before the point-in-poly check
				// since that's far more expensive
				loc := (int)(uint64(hex) % uint64(numHexagons))
				loopCount := 0
				for out[loc] != 0 {
					// If this branch is reached, we have exceeded the maximum
					// number of hexagons possible and need to clean up the
					// allocated memory.
					// LCOV_EXCL_START
					if loopCount > numHexagons {
						search = nil
						found = nil
						bboxes = nil

						return -1
					}
					// LCOV_EXCL_STOP
					if out[loc] == hex {
						break // Skip duplicates found
					}
					loc = (loc + 1) % numHexagons
					loopCount++
				}
				if out[loc] == hex {
					continue // Skip this hex, already exists in the out hash
				}

				// Check if the hexagon is in the polygon or not
				var hexCenter GeoCoord
				h3ToGeo(hex, &hexCenter)

				// If not, skip
				if !pointInsidePolygon(geoPolygon, bboxes, &hexCenter) {
					continue
				}

				// Otherwise set it in the output array
				out[loc] = hex

				// Set the hexagon in the found hash
				found[numFoundHexes] = hex
				numFoundHexes++
			}
			currentSearchNum++
			i++
		}

		// Swap the search and found pointers, copy the found hex count to the
		// search hex count, and zero everything related to the found memory.
		temp := search
		search = found
		found = temp
		for j := 0; j < numSearchHexes; j++ {
			found[j] = 0
		}
		numSearchHexes = numFoundHexes
		numFoundHexes = 0
		// Repeat until no new hexagons are found
	}
	// The out memory structure should be complete, end it here
	search = nil
	found = nil
	bboxes = nil
	return 0
}

/**
 * Internal: Create a vertex graph from a set of hexagons. It is the
 * responsibility of the caller to call destroyon VertexGraph the populated
 * graph, otherwise the memory in the graph nodes will not be freed.
 * @private
 * @param h3Set    Set of hexagons
 * @param numHexes Number of hexagons in the set
 * @param graph    Output graph
 */
func h3SetToVertexGraph(h3Set []H3Index, numHexes int, graph *VertexGraph) {
	if numHexes < 1 {
		// We still need to init the graph, or calls to destroywill VertexGraph
		// fail
		initVertexGraph(graph, 0, 0)
		return
	}

	res := H3_GET_RESOLUTION(h3Set[0])
	minBuckets := 6

	// TODO: Better way to calculate/guess?
	numBuckets := 0
	if numHexes > minBuckets {
		numBuckets = numHexes
	} else {
		numBuckets = minBuckets
	}

	initVertexGraph(graph, numBuckets, res)
	// Iterate through every hexagon
	for i := 0; i < numHexes; i++ {
		var vertices GeoBoundary
		var fromVtx *GeoCoord
		var toVtx *GeoCoord
		var edge *VertexNode

		h3ToGeoBoundary(h3Set[i], &vertices)

		// iterate through every edge
		for j := 0; j < vertices.numVerts; j++ {

			fromVtx = &vertices.Verts[j]
			toVtx = &vertices.Verts[(j+1)%vertices.numVerts]
			// If we've seen this edge already, it will be reversed
			edge = findNodeForEdge(graph, toVtx, fromVtx)
			if edge != nil {
				// If we've seen it, drop it. No edge is shared by more than 2
				// hexagons, so we'll never see it again.
				removeVertexNode(graph, edge)
			} else {
				// Add a new node for this edge
				addVertexNode(graph, fromVtx, toVtx)
			}
		}
	}
}

/**
 * Internal: Create a from Linkeda GeoPolygon vertex graph. It is the
 * responsibility of the caller to call destroyLinkedPolygon on the populated
 * linked geo structure, or the memory for that structure will not be freed.
 * @private
 * @param graph Input graph
 * @param out   Output polygon
 */
func _vertexGraphToLinkedGeo(graph *VertexGraph, out *LinkedGeoPolygon) {
	*out = LinkedGeoPolygon{}
	var loop *LinkedGeoLoop
	var nextVtx GeoCoord

	// Find the next unused entry point
	for edge := firstVertexNode(graph); edge != nil; edge = firstVertexNode(graph) {
		loop = addNewLinkedLoop(out)
		// Walk the graph to get the outline
		for {
			addLinkedCoord(loop, &edge.from)
			nextVtx = edge.to
			// Remove frees the node, so we can't use edge after this
			removeVertexNode(graph, edge)
			edge = findNodeForVertex(graph, &nextVtx)
			if edge == nil {
				break
			}
		}
	}
}

/**
 * Create a describing Linkedthe GeoPolygon outline(s) of a set of  hexagons.
 * Polygon outlines will follow GeoJSON MultiPolygon order: Each polygon will
 * have one outer loop, which is first in the list, followed by any holes.
 *
 * It is the responsibility of the caller to call destroyLinkedPolygon on the
 * populated linked geo structure, or the memory for that structure will
 * not be freed.
 *
 * It is expected that all hexagons in the set have the same resolution and
 * that the set contains no duplicates. Behavior is undefined if duplicates
 * or multiple resolutions are present, and the algorithm may produce
 * unexpected or invalid output.
 *
 * @param h3Set    Set of hexagons
 * @param numHexes Number of hexagons in set
 * @param out      Output polygon
 */
func h3SetToLinkedGeo(h3Set []H3Index, numHexes int, out *LinkedGeoPolygon) {
	var graph VertexGraph
	h3SetToVertexGraph(h3Set, numHexes, &graph)
	_vertexGraphToLinkedGeo(&graph, out)
	// TODO: The return value, possibly indicating an error, is discarded here -
	// we should use this when we update the API to return a value
	normalizeMultiPolygon(out)
	destroyVertexGraph(&graph)
}
