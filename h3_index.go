package h3

import (
	"math"
	"strconv"
)

type H3Index uint64

// define's of constants and macros for bitwise manipulation of H3Index's.

/** The number of bits in an H3 index. */
const H3_NUM_BITS = 64

/** The bit offset of the max resolution digit in an H3 index. */
const H3_MAX_OFFSET = 63

/** The bit offset of the mode in an H3 index. */
const H3_MODE_OFFSET = 59

/** The bit offset of the base cell in an H3 index. */
const H3_BC_OFFSET = 45

/** The bit offset of the resolution in an H3 index. */
const H3_RES_OFFSET = 52

/** The bit offset of the reserved bits in an H3 index. */
const H3_RESERVED_OFFSET = 56

/** The number of bits in a single H3 resolution digit. */
const H3_PER_DIGIT_OFFSET = 3

/** 1's in the 4 mode bits, 0's everywhere else. */
const H3_MODE_MASK = (uint64)(15) << H3_MODE_OFFSET

/** 0's in the 4 mode bits, 1's everywhere else. */
const H3_MODE_MASK_NEGATIVE = ^H3_MODE_MASK

/** 1's in the 7 base cell bits, 0's everywhere else. */
const H3_BC_MASK = (uint64)(127) << H3_BC_OFFSET

/** 0's in the 7 base cell bits, 1's everywhere else. */
const H3_BC_MASK_NEGATIVE = ^H3_BC_MASK

/** 1's in the 4 resolution bits, 0's everywhere else. */
const H3_RES_MASK = uint64(15) << H3_RES_OFFSET

/** 0's in the 4 resolution bits, 1's everywhere else. */
const H3_RES_MASK_NEGATIVE = ^H3_RES_MASK

/** 1's in the 3 reserved bits, 0's everywhere else. */
const H3_RESERVED_MASK = (uint64)(7) << H3_RESERVED_OFFSET

/** 0's in the 3 reserved bits, 1's everywhere else. */
const H3_RESERVED_MASK_NEGATIVE = ^H3_RESERVED_MASK

/** 1's in the 3 bits of res 15 digit bits, 0's everywhere else. */
const H3_DIGIT_MASK = Direction(7)

/** 0's in the 7 base cell bits, 1's everywhere else. */
//const H3_DIGIT_MASK_NEGATIVE = ^H3_DIGIT_MASK_NEGATIVE

/** H3 index with mode 0, res 0, base cell 0, and 7 for all index digits. */
const H3_INIT = H3Index(35184372088831)

/**
 * Gets the integer mode of h3.
 */
func H3_GET_MODE(h3 H3Index) H3Mode {
	return H3Mode(uint64(h3)&H3_MODE_MASK) >> H3_MODE_OFFSET
}

/**
 * Sets the integer mode of h3 to v.
 */
func H3_SET_MODE(h3 *H3Index, v H3Mode) {
	*h3 = H3Index((uint64(*h3) & H3_MODE_MASK_NEGATIVE) | uint64(v)<<H3_MODE_OFFSET)
}

/**
 * Gets the integer base cell of h3.
 */
func H3_GET_BASE_CELL(h3 H3Index) int {
	return int(uint64(h3)&H3_BC_MASK) >> H3_BC_OFFSET
}

/**
 * Sets the integer base cell of h3 to bc.
 */
func H3_SET_BASE_CELL(h3 *H3Index, bc int) {
	*h3 = H3Index((uint64(*h3) & H3_BC_MASK_NEGATIVE) | (uint64(bc) << H3_BC_OFFSET))
}

/**
 * Gets the integer resolution of h3.
 */
func H3_GET_RESOLUTION(h3 H3Index) int {
	return int(uint64(h3)&H3_RES_MASK) >> H3_RES_OFFSET
}

/**
 * Sets the integer resolution of h3.
 */
func H3_SET_RESOLUTION(h3 *H3Index, res int) {
	*h3 = H3Index((uint64(*h3) & H3_RES_MASK_NEGATIVE) | (uint64(res))<<H3_RES_OFFSET)
}

/**
 * Gets the resolution res integer digit (0-7) of h3.
 */
func H3_GET_INDEX_DIGIT(h3 H3Index, res int) Direction {
	return Direction((uint64(h3) >> ((MAX_H3_RES - (res)) * H3_PER_DIGIT_OFFSET)) & uint64(H3_DIGIT_MASK))
}

/**
 * Sets a value in the reserved space. Setting to non-zero may produce invalid
 * indexes.
 */
func H3_SET_RESERVED_BITS(h3 *H3Index, v int) {
	*h3 = H3Index((uint64(*h3) & H3_RESERVED_MASK_NEGATIVE) | (((uint64)(v)) << H3_RESERVED_OFFSET))
}

/**
 * Gets a value in the reserved space. Should always be zero for valid indexes.
 */
func H3_GET_RESERVED_BITS(h3 H3Index) int {
	return int((uint64(h3) & H3_RESERVED_MASK) >> H3_RESERVED_OFFSET)
}

/**
 * Sets the resolution res digit of h3 to the integer digit (0-7)
 */
func H3_SET_INDEX_DIGIT(h3 *H3Index, res int, digit Direction) {
	*h3 = H3Index((uint64(*h3) & ^(uint64(H3_DIGIT_MASK) << ((MAX_H3_RES - (res)) * H3_PER_DIGIT_OFFSET))) | (uint64(digit) << ((MAX_H3_RES - (res)) * H3_PER_DIGIT_OFFSET)))
}

/**
 * Invalid index used to indicate an error from geoToH3 and related functions.
 */
const H3_INVALID_INDEX = H3Index(0)

/**
 * Returns the H3 resolution of an H3 index.
 * @param h The H3 index.
 * @return The resolution of the H3 index argument.
 */
func h3GetResolution(h H3Index) int { return H3_GET_RESOLUTION(h) }

/**
 * Returns the H3 base cell number of an H3 index.
 * @param h The H3 index.
 * @return The base cell of the H3 index argument.
 */
func h3GetBaseCell(h H3Index) int { return H3_GET_BASE_CELL(h) }

/**
 * Converts a string representation of an H3 index into an H3 index.
 * @param str The string representation of an H3 index.
 * @return The H3 index corresponding to the string argument, or 0 if invalid.
 */
func stringToH3(str string) H3Index {
	h := H3_INVALID_INDEX
	// If failed, h will be unmodified and we should return 0 anyways.

	i, err := strconv.ParseUint(str, 16, 64)
	if err == nil {
		h = H3Index(i)
	}

	return h
}

/**
 * Converts an H3 index into a string representation.
 * @param h The H3 index to convert.
 */
func h3ToString(h H3Index) string {
	return strconv.FormatUint(uint64(h), 16)
}

/**
 * Returns whether or not an H3 index is valid.
 * @param h The H3 index to validate.
 */
func h3IsValid(h H3Index) bool {
	if H3_GET_MODE(h) != H3_HEXAGON_MODE {
		return false
	}
	baseCell := H3_GET_BASE_CELL(h)
	if baseCell < 0 || baseCell >= NUM_BASE_CELLS {
		return false
	}
	res := H3_GET_RESOLUTION(h)
	if res < 0 || res > MAX_H3_RES {
		return false
	}
	foundFirstNonZeroDigit := false
	for r := 1; r <= res; r++ {
		digit := H3_GET_INDEX_DIGIT(h, r)
		if !foundFirstNonZeroDigit && digit != CENTER_DIGIT {
			foundFirstNonZeroDigit = true
			if _isBaseCellPentagon(baseCell) && digit == K_AXES_DIGIT {
				return false
			}
		}

		if digit < CENTER_DIGIT || digit >= NUM_DIGITS {
			return false
		}
	}

	for r := res + 1; r <= MAX_H3_RES; r++ {
		digit := H3_GET_INDEX_DIGIT(h, r)
		if digit != INVALID_DIGIT {
			return false
		}
	}

	return true
}

/**
 * Initializes an H3 index.
 * @param hp The H3 index to initialize.
 * @param res The H3 resolution to initialize the index to.
 * @param baseCell The H3 base cell to initialize the index to.
 * @param initDigit The H3 digit (0-7) to initialize all of the index digits to.
 */
func setH3Index(hp *H3Index, res int, baseCell int, initDigit Direction) {
	h := H3_INIT
	H3_SET_MODE(&h, H3_HEXAGON_MODE)
	H3_SET_RESOLUTION(&h, res)
	H3_SET_BASE_CELL(&h, baseCell)
	for r := 1; r <= res; r++ {
		H3_SET_INDEX_DIGIT(&h, r, initDigit)
	}
	*hp = h
}

/**
 * h3ToParent produces the parent index for a given H3 index
 *
 * @param h to H3Index find parent of
 * @param parentRes The resolution to switch to (parent, grandparent, etc)
 *
 * @return of H3Index the parent, or 0 if you actually asked for a child
 */
func h3ToParent(h H3Index, parentRes int) H3Index {
	childRes := H3_GET_RESOLUTION(h)
	if parentRes > childRes {
		return H3_INVALID_INDEX
	} else if parentRes == childRes {
		return h
	} else if parentRes < 0 || parentRes > MAX_H3_RES {
		return H3_INVALID_INDEX
	}
	parentH := h
	H3_SET_RESOLUTION(&parentH, parentRes)
	for i := parentRes + 1; i <= childRes; i++ {
		H3_SET_INDEX_DIGIT(&parentH, i, H3_DIGIT_MASK)
	}
	return parentH
}

/**
 * Determines whether one resolution is a valid child resolution of another.
 * Each resolution is considered a valid child resolution of itself.
 *
 * @param parentRes resolution int of the parent
 * @param childRes resolution int of the child
 *
 * @return The validity of the child resolution
 */
func _isValidChildRes(parentRes int, childRes int) bool {
	if childRes < parentRes || childRes > MAX_H3_RES {
		return false
	}
	return true
}

/**
 * maxH3ToChildrenSize returns the maximum number of children possible for a
 * given child level.
 *
 * @param h to H3Index find the number of children of
 * @param childRes The resolution of the child level you're interested in
 *
 * @return count int of maximum number of children (equal for hexagons, less for
 * pentagons
 */
func maxH3ToChildrenSize(h H3Index, childRes int) int {
	parentRes := H3_GET_RESOLUTION(h)
	if !_isValidChildRes(parentRes, childRes) {
		return 0
	}
	return _ipow(7, childRes-parentRes)
}

/**
 * makeDirectChild takes an index and immediately returns the immediate child
 * index based on the specified cell number. Bit operations only, could generate
 * invalid indexes if not careful (deleted cell under a pentagon).
 *
 * @param h to H3Index find the direct child of
 * @param cellNumber id int of the direct child (0-6)
 *
 * @return The new for H3Index the child
 */
func makeDirectChild(h H3Index, cellNumber Direction) H3Index {
	childRes := H3_GET_RESOLUTION(h) + 1
	childH := h
	H3_SET_RESOLUTION(&childH, childRes)
	H3_SET_INDEX_DIGIT(&childH, childRes, cellNumber)
	return childH
}

/**
* h3ToChildren takes the given hexagon id and generates all of the children
* at the specified resolution storing them into the provided memory pointer.
* It's assumed that maxH3ToChildrenSize was used to determine the allocation.
*
* @param h to H3Index find the children of
* @param childRes the int child level to produce
* @param children H3Index* the memory to store the resulting addresses in
 */
func h3ToChildren(h H3Index, childRes int, children *[]H3Index) {
	parentRes := H3_GET_RESOLUTION(h)

	if !_isValidChildRes(parentRes, childRes) {
		return
	}

	if parentRes == childRes {
		*children = append(*children, h)
		return
	}

	bufferSize := maxH3ToChildrenSize(h, childRes)
	bufferChildStep := bufferSize / 7
	isAPentagon := h3IsPentagon(h)

	for i := Direction(0); i < 7; i++ {
		if isAPentagon && i == K_AXES_DIGIT {
			nextChild := make([]H3Index, bufferChildStep)
			for len(*children) < len(nextChild) {
				*children = append(*children, H3_INVALID_INDEX)
			}
		} else {
			h3ToChildren(makeDirectChild(h, i), childRes, children)
		}
	}
}

/**
* h3ToCenterChild produces the center child index for a given H3 index at
* the specified resolution
*
* @param h to H3Index find center child of
* @param childRes The resolution to switch to
*
* @return of H3Index the center child, or 0 if you actually asked for a parent
 */
func h3ToCenterChild(h H3Index, childRes int) H3Index {
	parentRes := H3_GET_RESOLUTION(h)
	if !_isValidChildRes(parentRes, childRes) {
		return H3_INVALID_INDEX
	}
	if childRes == parentRes {
		return h
	}
	child := h
	H3_SET_RESOLUTION(&child, childRes)
	for i := parentRes + 1; i <= childRes; i++ {
		H3_SET_INDEX_DIGIT(&child, i, 0)
	}
	return child
}

/**
* compact takes a set of hexagons all at the same resolution and compresses
* them by pruning full child branches to the parent level. This is also done
* for all parents recursively to get the minimum number of hex addresses that
* perfectly cover the defined space.
* @param h3Set Set of hexagons
* @param compactedSet The output array of compressed hexagons (preallocated)
* @param numHexes The size of the input and output arrays (possible that no
* contiguous regions exist in the set at all and no compression possible)
* @return an error code on bad input data
 */
func compact(h3Set []H3Index, compactedSet []H3Index, numHexes int) int {
	if numHexes == 0 {
		return 0
	}
	res := H3_GET_RESOLUTION(h3Set[0])
	if res == 0 {
		// No compaction possible, just copy the set to output
		for i := 0; i < numHexes; i++ {
			compactedSet[i] = h3Set[i]
		}
		return 0
	}
	remainingHexes := make([]H3Index, numHexes)
	copy(remainingHexes, h3Set)

	hashSetArray := make([]H3Index, numHexes)
	compactedSetOffset := compactedSet
	numRemainingHexes := numHexes
	for numRemainingHexes != 0 {
		res = H3_GET_RESOLUTION(remainingHexes[0])
		parentRes := res - 1
		// Put the parents of the hexagons into the temp array
		// via a hashing mechanism, and use the reserved bits
		// to track how many times a parent is duplicated
		for i := 0; i < numRemainingHexes; i++ {
			currIndex := remainingHexes[i]
			if currIndex != 0 {
				parent := h3ToParent(currIndex, parentRes)
				// Modulus hash the parent into the temp array
				loc := (int)(uint64(parent) % uint64(numRemainingHexes))
				loopCount := 0
				for hashSetArray[loc] != 0 {
					if loopCount > numRemainingHexes { // LCOV_EXCL_BR_LINE
						// LCOV_EXCL_START
						// This case should not be possible because at most one
						// index is placed into hashSetArray per
						// numRemainingHexes.
						remainingHexes = nil
						hashSetArray = nil
						return -1
						// LCOV_EXCL_STOP
					}
					tempIndex := H3Index(uint64(hashSetArray[loc]) & H3_RESERVED_MASK_NEGATIVE)
					if tempIndex == parent {
						count := H3_GET_RESERVED_BITS(hashSetArray[loc]) + 1
						if count > 7 {
							// Only possible on duplicate input
							remainingHexes = nil
							hashSetArray = nil
							return -2
						}
						H3_SET_RESERVED_BITS(&parent, count)
						hashSetArray[loc] = H3_INVALID_INDEX
					} else {
						loc = (loc + 1) % numRemainingHexes
					}
					loopCount++
				}
				hashSetArray[loc] = parent
			}
		}
		// Determine which parent hexagons have a complete set
		// of children and put them in the compactableHexes array
		compactableCount := 0
		maxCompactableCount := numRemainingHexes / 6 // Somehow all pentagons; conservative
		if maxCompactableCount == 0 {
			copy(compactedSetOffset, remainingHexes)
			break
		}

		compactableHexes := make([]H3Index, maxCompactableCount)
		for i := 0; i < numRemainingHexes; i++ {
			if hashSetArray[i] == 0 {
				continue
			}
			count := H3_GET_RESERVED_BITS(hashSetArray[i]) + 1
			// Include the deleted direction for pentagons as implicitly "there"
			if h3IsPentagon(H3Index(uint64(hashSetArray[i]) & H3_RESERVED_MASK_NEGATIVE)) {
				// We need this later on, no need to recalculate
				H3_SET_RESERVED_BITS(&hashSetArray[i], count)
				// Increment count after setting the reserved bits,
				// since count is already incremented above, so it
				// will be the expected value for a complete hexagon.
				count++
			}
			if count == 7 {
				// Bingo! Full set!
				compactableHexes[compactableCount] = H3Index(uint64(hashSetArray[i]) & H3_RESERVED_MASK_NEGATIVE)
				compactableCount++
			}
		}
		// Uncompactable hexes are immediately copied into the
		// output compactedSetOffset
		uncompactableCount := 0
		for i := 0; i < numRemainingHexes; i++ {
			currIndex := remainingHexes[i]
			if currIndex != H3_INVALID_INDEX {
				parent := h3ToParent(currIndex, parentRes)
				// Modulus hash the parent into the temp array
				// to determine if this index was included in
				// the compactableHexes array
				loc := (int)(uint64(parent) % uint64(numRemainingHexes))
				loopCount := 0
				isUncompactable := true
				for {
					if loopCount > numRemainingHexes { // LCOV_EXCL_BR_LINE
						// LCOV_EXCL_START
						// This case should not be possible because at most one
						// index is placed into hashSetArray per input hexagon.
						compactableHexes = nil
						remainingHexes = nil
						hashSetArray = nil
						return -1 // Only possible on duplicate input
						// LCOV_EXCL_STOP
					}
					tempIndex := H3Index(uint64(hashSetArray[loc]) & H3_RESERVED_MASK_NEGATIVE)
					if tempIndex == parent {
						count := H3_GET_RESERVED_BITS(hashSetArray[loc]) + 1
						if count == 7 {
							isUncompactable = false
						}
						break
					} else {
						loc = (loc + 1) % numRemainingHexes
					}
					loopCount++
					if hashSetArray[loc] != parent {
						break
					}
				}
				if isUncompactable {
					compactedSetOffset[uncompactableCount] = remainingHexes[i]
					uncompactableCount++
				}
			}
		}
		// Set up for the next loop
		//memset(hashSetArray, 0, numHexes * sizeof(H3Index));
		//compactedSetOffset += uncompactableCount;

		copy(remainingHexes, compactableHexes)
		numRemainingHexes = compactableCount
		compactableHexes = nil
	}

	remainingHexes = nil
	hashSetArray = nil

	return 0
}

/**
* uncompact takes a compressed set of hexagons and expands back to the
* original set of hexagons.
* @param compactedSet Set of hexagons
* @param numHexes The number of hexes in the input set
* @param h3Set Output array of decompressed hexagons (preallocated)
* @param maxHexes The size of the output array to bound check against
* @param res The hexagon resolution to decompress to
* @return An error code if output array is too small or any hexagon is
* smaller than the output resolution.
 */
func uncompact(compactedSet []H3Index, numHexes int, h3Set []H3Index, maxHexes int, res int) int {
	outOffset := 0
	for i := 0; i < numHexes; i++ {
		if compactedSet[i] == 0 {
			continue
		}
		if outOffset >= maxHexes {
			// We went too far, abort!
			return -1
		}
		currentRes := H3_GET_RESOLUTION(compactedSet[i])
		if !_isValidChildRes(currentRes, res) {
			// Nonsensical. Abort.
			return -2
		}
		if currentRes == res {
			// Just copy and move along
			h3Set[outOffset] = compactedSet[i]
			outOffset++
		} else {
			// Bigger hexagon to reduce in size
			numHexesToGen := maxH3ToChildrenSize(compactedSet[i], res)
			if outOffset+numHexesToGen > maxHexes {
				// We're about to go too far, abort!
				return -1
			}
			// todo fix
			//h3ToChildren(compactedSet[i], res, h3Set+outOffset)
			outOffset += numHexesToGen
		}
	}
	return 0
}

/**
* maxUncompactSize takes a compacted set of hexagons are provides an
* upper-bound estimate of the size of the uncompacted set of hexagons.
* @param compactedSet Set of hexagons
* @param numHexes The number of hexes in the input set
* @param res The hexagon resolution to decompress to
* @return The number of hexagons to allocate memory for, or a negative
* number if an error occurs.
 */
func maxUncompactSize(compactedSet []H3Index, numHexes int, res int) int {
	maxNumHexagons := 0
	for i := 0; i < numHexes; i++ {
		if compactedSet[i] == 0 {
			continue
		}

		currentRes := H3_GET_RESOLUTION(compactedSet[i])
		if !_isValidChildRes(currentRes, res) {
			// Nonsensical. Abort.
			return -1
		}
		if currentRes == res {
			maxNumHexagons++
		} else {
			// Bigger hexagon to reduce in size
			numHexesToGen := maxH3ToChildrenSize(compactedSet[i], res)
			maxNumHexagons += numHexesToGen
		}
	}
	return maxNumHexagons
}

/**
 * h3IsResClassIII takes a hexagon ID and determines if it is in a
 * Class III resolution (rotated versus the icosahedron and subject
 * to shape distortion adding extra points on icosahedron edges, making
 * them not true hexagons).
 * @param h The to H3Index check.
 * @return Returns 1 if the hexagon is class III, otherwise 0.
 */
func h3IsResClassIII(h H3Index) bool { return H3_GET_RESOLUTION(h)%2 != 0 }

//

/**
* h3IsPentagon takes an and H3Index determines if it is actually a
* pentagon.
* @param h The to H3Index check.
* @return Returns 1 if it is a pentagon, otherwise 0.
 */
func h3IsPentagon(h H3Index) bool {
	return _isBaseCellPentagon(H3_GET_BASE_CELL(h)) && _h3LeadingNonZeroDigit(h) == 0
}

/**
* Returns the highest resolution non-zero digit in an H3Index.
* @param h The H3Index.
* @return The highest resolution non-zero digit in the H3Index.
 */
func _h3LeadingNonZeroDigit(h H3Index) Direction {
	for r := 1; r <= H3_GET_RESOLUTION(h); r++ {
		d := H3_GET_INDEX_DIGIT(h, r)
		if d != 0 {
			return d
		}
	}
	// if we're here it's all 0's
	return CENTER_DIGIT
}

/**
* Rotate an 60 H3Index degrees counter-clockwise about a pentagonal center.
* @param h The H3Index.
 */
func _h3RotatePent60ccw(h H3Index) H3Index {
	// rotate in place; skips any leading 1 digits (k-axis)

	foundFirstNonZeroDigit := 0
	res := H3_GET_RESOLUTION(h)
	for r := 1; r <= res; r++ {
		// rotate this digit
		H3_SET_INDEX_DIGIT(&h, r, _rotate60ccw(H3_GET_INDEX_DIGIT(h, r)))

		// look for the first non-zero digit so we
		// can adjust for deleted k-axes sequence
		// if necessary
		if foundFirstNonZeroDigit == 0 && H3_GET_INDEX_DIGIT(h, r) != 0 {
			foundFirstNonZeroDigit = 1

			// adjust for deleted k-axes sequence
			if _h3LeadingNonZeroDigit(h) == K_AXES_DIGIT {
				h = _h3Rotate60ccw(h)
			}
		}
	}
	return h
}

/**
* Rotate an 60 H3Index degrees clockwise about a pentagonal center.
* @param h The H3Index.
 */
func _h3RotatePent60cw(h H3Index) H3Index {
	// rotate in place; skips any leading 1 digits (k-axis)
	foundFirstNonZeroDigit := 0
	res := H3_GET_RESOLUTION(h)
	for r := 1; r <= res; r++ {
		// rotate this digit
		H3_SET_INDEX_DIGIT(&h, r, _rotate60cw(H3_GET_INDEX_DIGIT(h, r)))

		// look for the first non-zero digit so we
		// can adjust for deleted k-axes sequence
		// if necessary
		if foundFirstNonZeroDigit == 0 && H3_GET_INDEX_DIGIT(h, r) != 0 {
			foundFirstNonZeroDigit = 1
			// adjust for deleted k-axes sequence
			if _h3LeadingNonZeroDigit(h) == K_AXES_DIGIT {
				h = _h3Rotate60cw(h)
			}
		}
	}
	return h
}

/**
 * Rotate an 60 H3Index degrees counter-clockwise.
 * @param h The H3Index.
 */
func _h3Rotate60ccw(h H3Index) H3Index {
	res := H3_GET_RESOLUTION(h)
	for r := 1; r <= res; r++ {
		oldDigit := H3_GET_INDEX_DIGIT(h, r)
		H3_SET_INDEX_DIGIT(&h, r, _rotate60ccw(oldDigit))
	}

	return h
}

/**
 * Rotate an 60 H3Index degrees clockwise.
 * @param h The H3Index.
 */
func _h3Rotate60cw(h H3Index) H3Index {
	res := H3_GET_RESOLUTION(h)

	for r := 1; r <= res; r++ {
		H3_SET_INDEX_DIGIT(&h, r, _rotate60cw(H3_GET_INDEX_DIGIT(h, r)))
	}

	return h
}

/**
* Convert an FaceIJK address to the corresponding H3Index.
* @param fijk The FaceIJK address.
* @param res The cell resolution.
* @return The encoded H3Index (or 0 on failure).
 */
func _faceIjkToH3(fijk *FaceIJK, res int) H3Index {
	// initialize the index
	h := H3_INIT
	H3_SET_MODE(&h, H3_HEXAGON_MODE)
	H3_SET_RESOLUTION(&h, res)

	// check for res 0/base cell
	if res == 0 {
		if fijk.coord.i > MAX_FACE_COORD || fijk.coord.j > MAX_FACE_COORD || fijk.coord.k > MAX_FACE_COORD {
			// out of range input
			return H3_INVALID_INDEX
		}

		H3_SET_BASE_CELL(&h, _faceIjkToBaseCell(fijk))
		return h
	}

	// we need to find the correct base cell FaceIJK for this H3 index;
	// start with the passed in face and resolution res ijk coordinates
	// in that face's coordinate system
	fijkBC := *fijk

	// build the from H3Index finest res up
	// adjust r for the fact that the res 0 base cell offsets the indexing
	// digits
	ijk := &fijkBC.coord
	for r := res - 1; r >= 0; r-- {
		lastIJK := *ijk
		var lastCenter CoordIJK
		if isResClassIII(r + 1) {
			// rotate ccw
			_upAp7(ijk)
			lastCenter = *ijk
			_downAp7(&lastCenter)
		} else {
			// rotate cw
			_upAp7r(ijk)
			lastCenter = *ijk
			_downAp7r(&lastCenter)
		}

		var diff CoordIJK
		_ijkSub(&lastIJK, &lastCenter, &diff)
		_ijkNormalize(&diff)
		H3_SET_INDEX_DIGIT(&h, r+1, _unitIjkToDigit(&diff))
	}

	// fijkBC should now hold the IJK of the base cell in the
	// coordinate system of the current face

	if fijkBC.coord.i > MAX_FACE_COORD || fijkBC.coord.j > MAX_FACE_COORD ||
		fijkBC.coord.k > MAX_FACE_COORD {
		// out of range input
		return H3_INVALID_INDEX
	}

	// lookup the correct base cell
	baseCell := _faceIjkToBaseCell(&fijkBC)
	H3_SET_BASE_CELL(&h, baseCell)

	// rotate if necessary to get canonical base cell orientation
	// for this base cell
	numRots := _faceIjkToBaseCellCCWrot60(&fijkBC)
	if _isBaseCellPentagon(baseCell) {
		// force rotation out of missing k-axes sub-sequence
		if _h3LeadingNonZeroDigit(h) == K_AXES_DIGIT {
			// check for a cw/ccw offset face; default is ccw
			if _baseCellIsCwOffset(baseCell, fijkBC.face) {
				h = _h3Rotate60cw(h)
			} else {
				h = _h3Rotate60ccw(h)
			}
		}

		for i := 0; i < numRots; i++ {
			h = _h3RotatePent60ccw(h)
		}
	} else {
		for i := 0; i < numRots; i++ {
			h = _h3Rotate60ccw(h)
		}
	}

	return h
}

/**
* Encodes a coordinate on the sphere to the H3 index of the containing cell at
* the specified resolution.
*
* Returns 0 on invalid input.
*
* @param g The spherical coordinates to encode.
* @param res The desired H3 resolution for the encoding.
* @return The encoded H3Index (or 0 on failure).
 */
func geoToH3(g *GeoCoord, res int) H3Index {
	if res < 0 || res > MAX_H3_RES {
		return H3_INVALID_INDEX
	}

	if math.IsInf(g.Lat, 0) || math.IsInf(g.Lon, 0) {
		return H3_INVALID_INDEX
	}

	fijk := FaceIJK{}
	_geoToFaceIjk(g, res, &fijk)
	return _faceIjkToH3(&fijk, res)
}

/**
* Convert an to H3Index the FaceIJK address on a specified icosahedral face.
* @param h The H3Index.
* @param fijk The FaceIJK address, initialized with the desired face
*        and normalized base cell coordinates.
* @return Returns 1 if the possibility of overage exists, otherwise 0.
 */
func _h3ToFaceIjkWithInitializedFijk(h H3Index, fijk *FaceIJK) int {
	ijk := &fijk.coord
	res := H3_GET_RESOLUTION(h)

	// center base cell hierarchy is entirely on this face
	possibleOverage := 1
	if !_isBaseCellPentagon(H3_GET_BASE_CELL(h)) && (res == 0 || (fijk.coord.i == 0 && fijk.coord.j == 0 && fijk.coord.k == 0)) {
		possibleOverage = 0
	}
	for r := 1; r <= res; r++ {
		if isResClassIII(r) {
			// Class III == rotate ccw
			_downAp7(ijk)
		} else {
			// Class II == rotate cw
			_downAp7r(ijk)
		}

		_neighbor(ijk, H3_GET_INDEX_DIGIT(h, r))
	}

	return possibleOverage
}

/**
* Convert an to H3Index a FaceIJK address.
* @param h The H3Index.
* @param fijk The corresponding FaceIJK address.
 */
func _h3ToFaceIjk(h H3Index, fijk *FaceIJK) {
	baseCell := H3_GET_BASE_CELL(h)
	// adjust for the pentagonal missing sequence; all of sub-sequence 5 needs
	// to be adjusted (and some of sub-sequence 4 below)
	if _isBaseCellPentagon(baseCell) && _h3LeadingNonZeroDigit(h) == 5 {
		h = _h3Rotate60cw(h)
	}

	// start with the "home" face and ijk+ coordinates for the base cell of c
	*fijk = baseCellData[baseCell].homeFijk

	if _h3ToFaceIjkWithInitializedFijk(h, fijk) == 0 {
		return // no overage is possible; h lies on this face
	}

	// if we're here we have the potential for an "overage"; i.e., it is
	// possible that c lies on an adjacent face

	origIJK := fijk.coord

	// if we're in Class III, drop into the next finer Class II grid
	res := H3_GET_RESOLUTION(h)
	if isResClassIII(res) {
		// Class III
		_downAp7r(&fijk.coord)
		res++
	}

	// adjust for overage if needed
	// a pentagon base cell with a leading 4 digit requires special handling
	pentLeading4 := 0

	if _isBaseCellPentagon(baseCell) && _h3LeadingNonZeroDigit(h) == 4 {
		pentLeading4 = 1
	}

	if _adjustOverageClassII(fijk, res, pentLeading4, 0) != NO_OVERAGE {
		// if the base cell is a pentagon we have the potential for secondary
		// overages
		if _isBaseCellPentagon(baseCell) {
			for _adjustOverageClassII(fijk, res, 0, 0) != NO_OVERAGE {

			}
		}

		if res != H3_GET_RESOLUTION(h) {
			_upAp7r(&fijk.coord)
		}

	} else if res != H3_GET_RESOLUTION(h) {
		fijk.coord = origIJK
	}
}

/**
* Determines the spherical coordinates of the center poof int an H3 index.
*
* @param h3 The H3 index.
* @param g The spherical coordinates of the H3 cell center.
 */
func h3ToGeo(h3 H3Index, g *GeoCoord) {
	var fijk FaceIJK
	_h3ToFaceIjk(h3, &fijk)
	_faceIjkToGeo(&fijk, H3_GET_RESOLUTION(h3), g)
}

/**
* Determines the cell boundary in spherical coordinates for an H3 index.
*
* @param h3 The H3 index.
* @param gb The boundary of the H3 cell in spherical coordinates.
 */
func h3ToGeoBoundary(h3 H3Index, gb *GeoBoundary) {
	var fijk FaceIJK
	_h3ToFaceIjk(h3, &fijk)
	_faceIjkToGeoBoundary(&fijk, H3_GET_RESOLUTION(h3), h3IsPentagon(h3), gb)
}

/**
* Returns the max number of possible icosahedron faces an H3 index
* may intersect.
*
* @return count int of faces
 */
func maxFaceCount(h3 H3Index) int {
	// a pentagon always intersects 5 faces, a hexagon never intersects more
	// than 2 (but may only intersect 1)
	if h3IsPentagon(h3) {
		return 5
	}
	return 2
}

/**
* Find all icosahedron faces intersected by a given H3 index, represented
* as integers from 0-19. The array is sparse; since 0 is a valid value,
* invalid array values are represented as -1. It is the responsibility of
* the caller to filter out invalid values.
*
* @param h3 The H3 index
* @param out Output array. Must be of size maxFaceCount(h3).
 */
func h3GetFaces(h3 H3Index, out []int) {
	res := H3_GET_RESOLUTION(h3)
	isPentagon := h3IsPentagon(h3)

	// We can't use the vertex-based approach here for class II pentagons,
	// because all their vertices are on the icosahedron edges. Their
	// direct child pentagons cross the same faces, so use those instead.
	if isPentagon && !isResClassIII(res) {
		// Note that this would not work for res 15, but this is only run on
		// Class II pentagons, it should never be invoked for a res 15 index.
		childPentagon := makeDirectChild(h3, 0)
		h3GetFaces(childPentagon, out)
		return
	}

	// convert to FaceIJK
	var fijk FaceIJK
	_h3ToFaceIjk(h3, &fijk)

	// Get all vertices as FaceIJK addresses. For simplicity, always
	// initialize the array with 6 Verts, ignoring the last one for pentagons
	var fijkVerts []FaceIJK
	var vertexCount int

	if isPentagon {
		vertexCount = NUM_PENT_VERTS
		_faceIjkPentToVerts(&fijk, &res, fijkVerts)
	} else {
		vertexCount = NUM_HEX_VERTS
		_faceIjkToVerts(&fijk, &res, fijkVerts)
	}

	// We may not use all of the slots in the output array,
	// so fill with invalid values to indicate unused slots
	faceCount := maxFaceCount(h3)
	for i := 0; i < faceCount; i++ {
		out[i] = INVALID_FACE
	}

	// add each vertex face, using the output array as a hash set
	for i := 0; i < vertexCount; i++ {
		vert := &fijkVerts[i]

		// Adjust overage, determining whether this vertex is
		// on another face
		if isPentagon {
			_adjustPentVertOverage(vert, res)
		} else {
			_adjustOverageClassII(vert, res, 0, 1)
		}

		// Save the face to the output array
		face := vert.face
		pos := 0
		// Find the first empty output position, or the first position
		// matching the current face
		for out[pos] != INVALID_FACE && out[pos] != face {
			pos++
		}
		out[pos] = face
	}
}

/**
 * pentagonIndexCount returns the number of pentagons (same at any resolution)
 *
 * @return count int of pentagon indexes
 */
func pentagonIndexCount() int {
	return NUM_PENTAGONS
}

/**
* Generates all pentagons at the specified resolution
*
* @param res The resolution to produce pentagons at.
* @param out Output array. Must be of size pentagonIndexCount().
 */
func getPentagonIndexes(res int, out *[]H3Index) {
	for bc := 0; bc < NUM_BASE_CELLS; bc++ {
		if _isBaseCellPentagon(bc) {
			var pentagon H3Index
			setH3Index(&pentagon, res, bc, 0)
			*out = append(*out, pentagon)
		}
	}
}

/**
* Returns whether or not a resolution is a Class III grid. Note that odd
* resolutions are Class III and even resolutions are Class II.
* @param res The H3 resolution.
* @return 1 if the resolution is a Class III grid, and 0 if the resolution is
*         a Class II grid.
 */
func isResClassIII(res int) bool {
	return res%2 == 1
}
