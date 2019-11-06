package h3

import (
	"fmt"
	"math"
)

/**
 * Origin leading digit . index leading digit . rotations 60 cw
 * Either being 1 (K axis) is invalid.
 * No good default at 0.
 */
var PENTAGON_ROTATIONS = [7][7]int{
	{0, -1, 0, 0, 0, 0, 0},       // 0
	{-1, -1, -1, -1, -1, -1, -1}, // 1
	{0, -1, 0, 0, 0, 1, 0},       // 2
	{0, -1, 0, 0, 1, 1, 0},       // 3
	{0, -1, 0, 5, 0, 0, 0},       // 4
	{0, -1, 5, 5, 0, 0, 0},       // 5
	{0, -1, 0, 0, 0, 0, 0},       // 6
}

/**
 * Reverse base cell direction . leading index digit . rotations 60 ccw.
 * For reversing the rotation introduced in PENTAGON_ROTATIONS when
 * the origin is on a pentagon (regardless of the base cell of the index.)
 */
var PENTAGON_ROTATIONS_REVERSE = [7][7]int{
	{0, 0, 0, 0, 0, 0, 0},        // 0
	{-1, -1, -1, -1, -1, -1, -1}, // 1
	{0, 1, 0, 0, 0, 0, 0},        // 2
	{0, 1, 0, 0, 0, 1, 0},        // 3
	{0, 5, 0, 0, 0, 0, 0},        // 4
	{0, 5, 0, 5, 0, 0, 0},        // 5
	{0, 0, 0, 0, 0, 0, 0},        // 6
}

/**
 * Reverse base cell direction . leading index digit . rotations 60 ccw.
 * For reversing the rotation introduced in PENTAGON_ROTATIONS when the index is
 * on a pentagon and the origin is not.
 */
var PENTAGON_ROTATIONS_REVERSE_NONPOLAR = [7][7]int{
	{0, 0, 0, 0, 0, 0, 0},        // 0
	{-1, -1, -1, -1, -1, -1, -1}, // 1
	{0, 1, 0, 0, 0, 0, 0},        // 2
	{0, 1, 0, 0, 0, 1, 0},        // 3
	{0, 5, 0, 0, 0, 0, 0},        // 4
	{0, 1, 0, 5, 1, 1, 0},        // 5
	{0, 0, 0, 0, 0, 0, 0},        // 6
}

/**
 * Reverse base cell direction . leading index digit . rotations 60 ccw.
 * For reversing the rotation introduced in PENTAGON_ROTATIONS when the index is
 * on a polar pentagon and the origin is not.
 */
var PENTAGON_ROTATIONS_REVERSE_POLAR = [7][7]int{
	{0, 0, 0, 0, 0, 0, 0},        // 0
	{-1, -1, -1, -1, -1, -1, -1}, // 1
	{0, 1, 1, 1, 1, 1, 1},        // 2
	{0, 1, 0, 0, 0, 1, 0},        // 3
	{0, 1, 0, 0, 1, 1, 1},        // 4
	{0, 1, 0, 5, 1, 1, 0},        // 5
	{0, 1, 1, 0, 1, 1, 1},        // 6
}

/**
 * Prohibited directions when unfolding a pentagon.
 *
 * Indexes by two directions, both relative to the pentagon base cell. The first
 * is the direction of the origin index and the second is the direction of the
 * index to unfold. Direction refers to the direction from base cell to base
 * cell if the indexes are on different base cells, or the leading digit if
 * within the pentagon base cell.
 *
 * This previously included a Class II/Class III check but these were removed
 * due to failure cases. It's possible this could be restricted to a narrower
 * set of a failure cases. Currently, the logic is any unfolding across more
 * than one icosahedron face is not permitted.
 */
var FAILED_DIRECTIONS = [7][7]bool{
	{false, false, false, false, false, false, false}, // 0
	{false, false, false, false, false, false, false}, // 1
	{false, false, false, false, true, true, false},   // 2
	{false, false, false, false, true, false, true},   // 3
	{false, false, true, true, false, false, false},   // 4
	{false, false, true, false, false, false, true},   // 5
	{false, false, false, true, false, true, false},   // 6
}

/**
* Produces ijk+ coordinates for an index anchored by an origin.
*
* The coordinate space used by this function may have deleted
* regions or warping due to pentagonal distortion.
*
* Coordinates are only comparable if they come from the same
* origin index.
*
* Failure may occur if the index is too far away from the origin
* or if the index is on the other side of a pentagon.
*
* @param origin An anchoring index for the ijk+ coordinate system.
* @param index Index to find the coordinates of
* @param out ijk+ coordinates of the index will be placed here on success
* @return 0 on success, or another value on failure.
 */
func h3ToLocalIjk(origin H3Index, h3 H3Index, out *CoordIJK) int {
	res := H3_GET_RESOLUTION(origin)
	if res != H3_GET_RESOLUTION(h3) {
		return 1
	}

	originBaseCell := H3_GET_BASE_CELL(origin)
	baseCell := H3_GET_BASE_CELL(h3)

	// Direction from origin base cell to index base cell
	dir := CENTER_DIGIT
	revDir := CENTER_DIGIT
	if originBaseCell != baseCell {
		dir = _getBaseCellDirection(originBaseCell, baseCell)
		if dir == INVALID_DIGIT {
			// Base cells are not neighbors, can't unfold.
			return 2
		}
		revDir = _getBaseCellDirection(baseCell, originBaseCell)
	}

	originOnPent := _isBaseCellPentagon(originBaseCell)
	indexOnPent := _isBaseCellPentagon(baseCell)
	indexFijk := FaceIJK{}
	if dir != CENTER_DIGIT {
		// Rotate index into the orientation of the origin base cell.
		// cw because we are undoing the rotation into that base cell.
		baseCellRotations := baseCellNeighbor60CCWRots[originBaseCell][dir]
		if indexOnPent {
			for i := 0; i < baseCellRotations; i++ {
				h3 = _h3RotatePent60cw(h3)
				revDir = _rotate60cw(revDir)
				if revDir == K_AXES_DIGIT {
					revDir = _rotate60cw(revDir)
				}
			}
		} else {
			for i := 0; i < baseCellRotations; i++ {
				h3 = _h3Rotate60cw(h3)
				revDir = _rotate60cw(revDir)
			}
		}
	}

	// Face is unused. This produces coordinates in base cell coordinate space.
	_h3ToFaceIjkWithInitializedFijk(h3, &indexFijk)
	if dir != CENTER_DIGIT {
		if baseCell == originBaseCell {
			panic(fmt.Errorf("baseCell should not equal originBaseCell"))
		}

		pentagonRotations := 0
		directionRotations := 0
		if originOnPent {
			originLeadingDigit := _h3LeadingNonZeroDigit(origin)
			if FAILED_DIRECTIONS[originLeadingDigit][dir] {
				// TODO: We may be unfolding the pentagon incorrectly in this
				// case; return an error code until this is guaranteed to be
				// correct.
				return 3
			}

			directionRotations = PENTAGON_ROTATIONS[originLeadingDigit][dir]
			pentagonRotations = directionRotations
		} else if indexOnPent {
			indexLeadingDigit := _h3LeadingNonZeroDigit(h3)
			if FAILED_DIRECTIONS[indexLeadingDigit][revDir] {
				// TODO: We may be unfolding the pentagon incorrectly in this
				// case; return an error code until this is guaranteed to be
				// correct.
				return 4
			}

			pentagonRotations = PENTAGON_ROTATIONS[revDir][indexLeadingDigit]
		}

		if !(pentagonRotations >= 0) {
			panic(fmt.Errorf("pentagonRotations should be large than 0"))
		}

		if !(directionRotations >= 0) {
			panic(fmt.Errorf("directionRotations should be large than 0"))
		}

		for i := 0; i < pentagonRotations; i++ {
			_ijkRotate60cw(&indexFijk.coord)
		}

		offset := CoordIJK{}
		_neighbor(&offset, dir)

		// Scale offset based on resolution
		for r := res - 1; r >= 0; r-- {
			if isResClassIII(r + 1) {
				// rotate ccw
				_downAp7(&offset)
			} else {
				// rotate cw
				_downAp7r(&offset)
			}
		}

		for i := 0; i < directionRotations; i++ {
			_ijkRotate60cw(&offset)
		}

		// Perform necessary translation
		_ijkAdd(&indexFijk.coord, &offset, &indexFijk.coord)
		_ijkNormalize(&indexFijk.coord)
	} else if originOnPent && indexOnPent {
		// If the origin and index are on pentagon, and we checked that the base
		// cells are the same or neighboring, then they must be the same base
		// cell.
		if !(baseCell == originBaseCell) {
			panic(fmt.Sprintf("must be same base cell"))
		}

		originLeadingDigit := _h3LeadingNonZeroDigit(origin)
		indexLeadingDigit := _h3LeadingNonZeroDigit(h3)
		if FAILED_DIRECTIONS[originLeadingDigit][indexLeadingDigit] {
			// TODO: We may be unfolding the pentagon incorrectly in this case;
			// return an error code until this is guaranteed to be correct.
			return 5
		}

		withinPentagonRotations := PENTAGON_ROTATIONS[originLeadingDigit][indexLeadingDigit]
		for i := 0; i < withinPentagonRotations; i++ {
			_ijkRotate60cw(&indexFijk.coord)
		}
	}

	*out = indexFijk.coord
	return 0
}

/**
* Produces an index for ijk+ coordinates anchored by an origin.
*
* The coordinate space used by this function may have deleted
* regions or warping due to pentagonal distortion.
*
* Failure may occur if the coordinates are too far away from the origin
* or if the index is on the other side of a pentagon.
*
* @param origin An anchoring index for the ijk+ coordinate system.
* @param ijk IJK+ Coordinates to find the index of
* @param out The index will be placed here on success
* @return 0 on success, or another value on failure.
 */
func localIjkToH3(origin H3Index, ijk *CoordIJK, out *H3Index) int {
	res := H3_GET_RESOLUTION(origin)
	originBaseCell := H3_GET_BASE_CELL(origin)
	originOnPent := _isBaseCellPentagon(originBaseCell)

	// This logic is very similar to faceIjkToH3
	// initialize the index
	*out = H3_INIT
	H3_SET_MODE(out, H3_HEXAGON_MODE)
	H3_SET_RESOLUTION(out, res)

	// check for res 0/base cell
	if res == 0 {
		if ijk.i > 1 || ijk.j > 1 || ijk.k > 1 {
			// out of range input
			return 1
		}

		dir := _unitIjkToDigit(ijk)
		newBaseCell := _getBaseCellNeighbor(originBaseCell, dir)
		if newBaseCell == INVALID_BASE_CELL {
			// Moving in an invalid direction off a pentagon.
			return 1
		}
		H3_SET_BASE_CELL(out, newBaseCell)
		return 0
	}

	// we need to find the correct base cell offset (if any) for this H3 index;
	// start with the passed in base cell and resolution res ijk coordinates
	// in that base cell's coordinate system
	ijkCopy := *ijk

	// build the from H3Index finest res up
	// adjust r for the fact that the res 0 base cell offsets the indexing
	// digits
	for r := res - 1; r >= 0; r-- {
		lastIJK := &CoordIJK{}
		var lastCenter CoordIJK
		if isResClassIII(r + 1) {
			// rotate ccw
			_upAp7(&ijkCopy)
			lastCenter = ijkCopy
			_downAp7(&lastCenter)
		} else {
			// rotate cw
			_upAp7r(&ijkCopy)
			lastCenter = ijkCopy
			_downAp7r(&lastCenter)
		}

		var diff CoordIJK
		_ijkSub(lastIJK, &lastCenter, &diff)
		_ijkNormalize(&diff)
		H3_SET_INDEX_DIGIT(out, r+1, _unitIjkToDigit(&diff))
	}

	// ijkCopy should now hold the IJK of the base cell in the
	// coordinate system of the current base cell

	if ijkCopy.i > 1 || ijkCopy.j > 1 || ijkCopy.k > 1 {
		// out of range input
		return 2
	}

	// lookup the correct base cell
	dir := _unitIjkToDigit(&ijkCopy)
	baseCell := _getBaseCellNeighbor(originBaseCell, dir)
	// If baseCell is invalid, it must be because the origin base cell is a
	// pentagon, and because pentagon base cells do not border each other,
	// baseCell must not be a pentagon.
	indexOnPent := false
	if baseCell != INVALID_BASE_CELL {
		indexOnPent = _isBaseCellPentagon(baseCell)
	}

	if dir != CENTER_DIGIT {
		// If the index is in a warped direction, we need to unwarp the base
		// cell direction. There may be further need to rotate the index digits.
		pentagonRotations := 0
		if originOnPent {
			originLeadingDigit := _h3LeadingNonZeroDigit(origin)
			pentagonRotations = PENTAGON_ROTATIONS_REVERSE[originLeadingDigit][dir]
			for i := 0; i < pentagonRotations; i++ {
				dir = _rotate60ccw(dir)
			}
			// The pentagon rotations are being chosen so that dir is not the
			// deleted direction. If it still happens, it means we're moving
			// into a deleted subsequence, so there is no index here.
			if dir == K_AXES_DIGIT {
				return 3
			}
			baseCell = _getBaseCellNeighbor(originBaseCell, dir)

			// indexOnPent does not need to be checked again since no pentagon
			// base cells border each other.
			//assert(baseCell != INVALID_BASE_CELL);
			//assert(!_isBaseCellPentagon(baseCell));
		}

		// Now we can determine the relation between the origin and target base
		// cell.
		baseCellRotations := baseCellNeighbor60CCWRots[originBaseCell][dir]
		//assert(baseCellRotations >= 0);

		// Adjust for pentagon warping within the base cell. The base cell
		// should be in the right location, so now we need to rotate the index
		// back. We might not need to check for errors since we would just be
		// double mapping.
		if indexOnPent {
			revDir := _getBaseCellDirection(baseCell, originBaseCell)
			//assert(revDir != INVALID_DIGIT);

			// Adjust for the different coordinate space in the two base cells.
			// This is done first because we need to do the pentagon rotations
			// based on the leading digit in the pentagon's coordinate system.
			for i := 0; i < baseCellRotations; i++ {
				*out = _h3Rotate60ccw(*out)
			}

			indexLeadingDigit := _h3LeadingNonZeroDigit(*out)

			if _isBaseCellPolarPentagon(baseCell) {
				pentagonRotations = PENTAGON_ROTATIONS_REVERSE_POLAR[revDir][indexLeadingDigit]
			} else {
				pentagonRotations = PENTAGON_ROTATIONS_REVERSE_NONPOLAR[revDir][indexLeadingDigit]
			}

			//assert(pentagonRotations >= 0);
			for i := 0; i < pentagonRotations; i++ {
				*out = _h3RotatePent60ccw(*out)
			}
		} else {
			//assert(pentagonRotations >= 0);
			for i := 0; i < pentagonRotations; i++ {
				*out = _h3Rotate60ccw(*out)
			}

			// Adjust for the different coordinate space in the two base cells.
			for i := 0; i < baseCellRotations; i++ {
				*out = _h3Rotate60ccw(*out)
			}
		}
	} else if originOnPent && indexOnPent {
		originLeadingDigit := _h3LeadingNonZeroDigit(origin)
		indexLeadingDigit := _h3LeadingNonZeroDigit(*out)
		withinPentagonRotations := PENTAGON_ROTATIONS_REVERSE[originLeadingDigit][indexLeadingDigit]
		//assert(withinPentagonRotations >= 0);

		for i := 0; i < withinPentagonRotations; i++ {
			*out = _h3Rotate60ccw(*out)
		}
	}

	if indexOnPent {
		// TODO: There are cases in h3ToLocalIjk which are failed but not
		// accounted for here - instead just fail if the recovered index is
		// invalid.
		if _h3LeadingNonZeroDigit(*out) == K_AXES_DIGIT {
			return 4
		}
	}

	H3_SET_BASE_CELL(out, baseCell)
	return 0
}

/**
* Produces ij coordinates for an index anchored by an origin.
*
* The coordinate space used by this function may have deleted
* regions or warping due to pentagonal distortion.
*
* Coordinates are only comparable if they come from the same
* origin index.
*
* Failure may occur if the index is too far away from the origin
* or if the index is on the other side of a pentagon.
*
* This function is experimental, and its output is not guaranteed
* to be compatible across different versions of H3.
*
* @param origin An anchoring index for the ij coordinate system.
* @param index Index to find the coordinates of
* @param out ij coordinates of the index will be placed here on success
* @return 0 on success, or another value on failure.
 */
func experimentalH3ToLocalIj(origin H3Index, h3 H3Index, out *CoordIJ) int {
	// This function is currently experimental. Once ready to be part of the
	// non-experimental API, this function (with the experimental prefix) will
	// be marked as deprecated and to be removed in the next major version. It
	// will be replaced with a non-prefixed function name.
	var ijk CoordIJK
	failed := h3ToLocalIjk(origin, h3, &ijk)

	if failed != 0 {
		return failed
	}

	ijkToIj(&ijk, out)

	return 0
}

/**
* Produces an index for ij coordinates anchored by an origin.
*
* The coordinate space used by this function may have deleted
* regions or warping due to pentagonal distortion.
*
* Failure may occur if the index is too far away from the origin
* or if the index is on the other side of a pentagon.
*
* This function is experimental, and its output is not guaranteed
* to be compatible across different versions of H3.
*
* @param origin An anchoring index for the ij coordinate system.
* @param out ij coordinates to index.
* @param index Index will be placed here on success.
* @return 0 on success, or another value on failure.
 */
func experimentalLocalIjToH3(origin H3Index, ij *CoordIJ, out *H3Index) int {
	// This function is currently experimental. Once ready to be part of the
	// non-experimental API, this function (with the experimental prefix) will
	// be marked as deprecated and to be removed in the next major version. It
	// will be replaced with a non-prefixed function name.
	var ijk CoordIJK
	ijToIjk(ij, &ijk)
	return localIjkToH3(origin, &ijk, out)
}

/**
* Produces the grid distance between the two indexes.
*
* This function may fail to find the distance between two indexes, for
* example if they are very far apart. It may also fail when finding
* distances for indexes on opposite sides of a pentagon.
*
* @param origin Index to find the distance from.
* @param index Index to find the distance to.
* @return The distance, or a negative number if the library could not
* compute the distance.
 */
func h3Distance(origin H3Index, h3 H3Index) int {
	var originIjk, h3Ijk CoordIJK
	if h3ToLocalIjk(origin, origin, &originIjk) != 0 {
		// Currently there are no tests that would cause getting the coordinates
		// for an index the same as the origin to fail.
		return -1 // LCOV_EXCL_LINE
	}

	if h3ToLocalIjk(origin, h3, &h3Ijk) != 0 {
		return -1
	}

	return int(ijkDistance(&originIjk, &h3Ijk))
}

/**
* Number of indexes in a line from the start index to the end index,
* to be used for allocating memory. Returns a negative number if the
* line cannot be computed.
*
* @param start Start index of the line
* @param end End index of the line
* @return Size of the line, or a negative number if the line cannot
* be computed.
 */
func h3LineSize(start H3Index, end H3Index) int {
	distance := h3Distance(start, end)
	if distance >= 0 {
		return distance + 1
	}
	return distance
}

/**
* Given cube coords as doubles, round to valid integer coordinates. Algorithm
* from https://www.redblobgames.com/grids/hexagons/#rounding
* @param i   Floating-point I coord
* @param j   Floating-point J coord
* @param k   Floating-point K coord
* @param ijk IJK coord struct, modified in place
 */
func cubeRound(i, j, k float64, ijk *CoordIJK) {
	ri := math.Round(i)
	rj := math.Round(j)
	rk := math.Round(k)
	iDiff := math.Abs(ri - i)
	jDiff := math.Abs(rj - j)
	kDiff := math.Abs(rk - k)

	// Round, maintaining valid cube coords
	if iDiff > jDiff && iDiff > kDiff {
		ri = -rj - rk
	} else if jDiff > kDiff {
		rj = -ri - rk
	} else {
		rk = -ri - rj
	}

	ijk.i = int(ri)
	ijk.j = int(rj)
	ijk.k = int(rk)
}

/**
* Given two H3 indexes, return the line of indexes between them (inclusive).
*
* This function may fail to find the line between two indexes, for
* example if they are very far apart. It may also fail when finding
* distances for indexes on opposite sides of a pentagon.
*
* Notes:
*
*  - The specific output of this function should not be considered stable
*    across library versions. The only guarantees the library provides are
*    that the line length will be `h3Distance(start, end) + 1` and that
*    every index in the line will be a neighbor of the preceding index.
*  - Lines are drawn in grid space, and may not correspond exactly to either
*    Cartesian lines or great arcs.
*
* @param start Start index of the line
* @param end End index of the line
* @param out Output array, which must be of size h3LineSize(start, end)
* @return 0 on success, or another value on failure.
 */
func h3Line(start H3Index, end H3Index, out []H3Index) int {
	distance := h3Distance(start, end)
	// Early exit if we can't calculate the line
	if distance < 0 {
		return distance
	}

	// Get IJK coords for the start and end. We've already confirmed
	// that these can be calculated with the distance check above.
	startIjk := &CoordIJK{}
	endIjk := &CoordIJK{}

	// Convert H3 addresses to IJK coords
	h3ToLocalIjk(start, start, startIjk)
	h3ToLocalIjk(start, end, endIjk)

	// Convert IJK to cube coordinates suitable for linear interpolation
	ijkToCube(startIjk)
	ijkToCube(endIjk)
	iStep := func() float64 {
		if distance != 0 {
			return float64(endIjk.i-startIjk.i) / float64(distance)
		}
		return 0
	}()
	jStep := func() float64 {
		if distance != 0 {
			return float64(endIjk.j-startIjk.j) / float64(distance)
		}
		return 0
	}()

	kStep := func() float64 {
		if distance != 0 {
			return float64(endIjk.k-startIjk.k) / float64(distance)
		}
		return 0
	}()

	currentIjk := &CoordIJK{startIjk.i, startIjk.j, startIjk.k}

	for n := 0; n <= distance; n++ {
		cubeRound(float64(startIjk.i)+iStep*float64(n), float64(startIjk.j)+jStep*float64(n), float64(startIjk.k)+kStep*float64(n), currentIjk)
		// Convert cube . ijk . h3 index
		cubeToIjk(currentIjk)
		localIjkToH3(start, currentIjk, &out[n])
	}

	return 0
}
