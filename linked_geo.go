package h3

const NORMALIZATION_SUCCESS = 0
const NORMALIZATION_ERR_MULTIPLE_POLYGONS = 1
const NORMALIZATION_ERR_UNASSIGNED_HOLES = 2

/**
 * Add a linked polygon to the current polygon
 * @param  polygon Polygon to add link to
 * @return         Pointer to new polygon
 */
func addNewLinkedPolygon(polygon *LinkedGeoPolygon) *LinkedGeoPolygon {
	next := &LinkedGeoPolygon{}
	polygon.next = next
	return next
}

/**
 * Add a new linked loop to the current polygon
 * @param  polygon Polygon to add loop to
 * @return         Pointer to loop
 */
func addNewLinkedLoop(polygon *LinkedGeoPolygon) *LinkedGeoLoop {
	loop := &LinkedGeoLoop{}
	return addLinkedLoop(polygon, loop)
}

/**
* Add an existing linked loop to the current polygon
* @param  polygon Polygon to add loop to
* @return         Pointer to loop
 */
func addLinkedLoop(polygon *LinkedGeoPolygon, loop *LinkedGeoLoop) *LinkedGeoLoop {
	last := polygon.last
	if last == nil {
		polygon.first = loop
	} else {
		last.next = loop
	}

	polygon.last = loop
	return loop
}

/**
* Add a new linked coordinate to the current loop
* @param  loop   Loop to add coordinate to
* @param  vertex Coordinate to add
* @return        Pointer to the coordinate
 */
func addLinkedCoord(loop *LinkedGeoLoop, vertex *GeoCoord) *LinkedGeoCoord {
	coord := &LinkedGeoCoord{
		vertex: *vertex,
	}

	last := loop.last
	if last == nil {
		loop.first = coord
	} else {
		last.next = coord
	}
	loop.last = coord

	return coord
}

/**
* Free all allocated memory for a linked geo loop. The caller is
* responsible for freeing memory allocated to input loop struct.
* @param loop Loop to free
 */
func destroyLinkedGeoLoop(loop *LinkedGeoLoop) {
	var nextCoord *LinkedGeoCoord

	for currentCoord := loop.first; currentCoord != nil; currentCoord = nextCoord {
		nextCoord = currentCoord.next
		currentCoord = nil
	}
}

/**
* Free all allocated memory for a linked geo structure. The caller is
* responsible for freeing memory allocated to input polygon struct.
* @param polygon Pointer to the first polygon in the structure
 */
func destroyLinkedPolygon(polygon *LinkedGeoPolygon) {
	// flag to skip the input polygon
	skip := true
	var nextPolygon *LinkedGeoPolygon
	var nextLoop *LinkedGeoLoop
	for currentPolygon := polygon; currentPolygon != nil; currentPolygon = nextPolygon {
		for currentLoop := currentPolygon.first; currentLoop != nil; currentLoop = nextLoop {
			destroyLinkedGeoLoop(currentLoop)
			nextLoop = currentLoop.next
			currentLoop = nil
		}
		nextPolygon = currentPolygon.next
		if skip {
			// do not free the input polygon
			skip = false
		} else {
			currentPolygon = nil
		}
	}
}

/**
* Count the number of polygons in a linked list
* @param  polygon Starting polygon
* @return         Count
 */
func countLinkedPolygons(polygon *LinkedGeoPolygon) int {
	count := 0
	for polygon != nil {
		count++
		polygon = polygon.next
	}
	return count
}

/**
* Count the number of linked loops in a polygon
* @param  polygon Polygon to count loops for
* @return         Count
 */
func countLinkedLoops(polygon *LinkedGeoPolygon) int {
	loop := polygon.first
	count := 0
	for loop != nil {
		count++
		loop = loop.next
	}
	return count
}

/**
* Count the number of coordinates in a loop
* @param  loop Loop to count coordinates for
* @return      Count
 */
func countLinkedCoords(loop *LinkedGeoLoop) int {
	coord := loop.first
	count := 0
	for coord != nil {
		count++
		coord = coord.next
	}
	return count
}

/**
* Count the number of polygons containing a given loop.
* @param  loop         Loop to count containers for
* @param  polygons     Polygons to test
* @param  bboxes       Bounding boxes for polygons, used in point-in-poly check
* @param  polygonCount Number of polygons in the test array
* @return              Number of polygons containing the loop
 */
func countContainers(loop *LinkedGeoLoop, polygons []*LinkedGeoPolygon, bboxes []*BBox, polygonCount int) int {
	containerCount := 0
	for i := 0; i < polygonCount; i++ {
		if loop != polygons[i].first && pointInside(polygons[i].first, bboxes[i], &loop.first.vertex) {
			containerCount++
		}
	}
	return containerCount
}

/**
* Given a list of nested containers, find the one most deeply nested.
* @param  polygons     Polygon containers to check
* @param  bboxes       Bounding boxes for polygons, used in point-in-poly check
* @param  polygonCount Number of polygons in the list
* @return              Deepest container, or null if list is empty
 */
func findDeepestContainer(polygons []*LinkedGeoPolygon, bboxes []*BBox, polygonCount int) *LinkedGeoPolygon {
	// Set the initial return value to the first candidate
	var parent *LinkedGeoPolygon
	if polygonCount > 0 {
		parent = polygons[0]
	}

	// If we have multiple polygons, they must be nested inside each other.
	// Find the innermost polygon by taking the one with the most containers
	// in the list.
	if polygonCount > 1 {
		max := -1
		for i := 0; i < polygonCount; i++ {
			count := countContainers(polygons[i].first, polygons, bboxes,
				polygonCount)
			if count > max {
				parent = polygons[i]
				max = count
			}
		}
	}

	return parent
}

/**
* Find the polygon to which a given hole should be allocated. Note that this
* function will return null if no parent is found.
* @param  loop         Inner loop describing a hole
* @param  polygon      Head of a linked list of polygons to check
* @param  bboxes       Bounding boxes for polygons, used in point-in-poly check
* @param  polygonCount Number of polygons to check
* @return              Pointer to parent polygon, or null if not found
 */
func findPolygonForHole(loop *LinkedGeoLoop, polygon *LinkedGeoPolygon, bboxes []BBox, polygonCount int) *LinkedGeoPolygon {
	// Early exit with no polygons
	if polygonCount == 0 {
		return nil
	}
	// Initialize arrays for candidate loops and their bounding boxes
	candidates := make([]*LinkedGeoPolygon, polygonCount)
	candidateBBoxes := make([]*BBox, polygonCount)

	// Find all polygons that contain the loop
	candidateCount := 0
	index := 0
	for polygon != nil {
		// We are guaranteed not to overlap, so just test the first point
		if pointInside(polygon.first, &bboxes[index], &loop.first.vertex) {
			candidates[candidateCount] = polygon
			candidateBBoxes[candidateCount] = &bboxes[index]
			candidateCount++
		}
		polygon = polygon.next
		index++
	}

	// The most deeply nested container is the immediate parent
	parent := findDeepestContainer(candidates, candidateBBoxes, candidateCount)

	// Free allocated memory
	candidates = nil
	candidateBBoxes = nil

	return parent
}

/**
* Normalize a LinkedGeoPolygon in-place into a structure following GeoJSON
* MultiPolygon rules: Each polygon must have exactly one outer loop, which
* must be first in the list, followed by any holes. Holes in this algorithm
* are identified by winding order (holes are clockwise), which is guaranteed
* by the h3SetToVertexGraph algorithm.
*
* Input to this function is assumed to be a single polygon including all
* loops to normalize. It's assumed that a valid arrangement is possible.
*
* @param root Root polygon including all loops
* @return     0 on success, or an error code > 0 for invalid input
 */
func normalizeMultiPolygon(root *LinkedGeoPolygon) int {
	// We assume that the input is a single polygon with loops;
	// if it has multiple polygons, don't touch it
	if root.next != nil {
		return NORMALIZATION_ERR_MULTIPLE_POLYGONS
	}

	// Count loops, exiting early if there's only one
	loopCount := countLinkedLoops(root)
	if loopCount <= 1 {
		return NORMALIZATION_SUCCESS
	}

	resultCode := NORMALIZATION_SUCCESS
	var polygon *LinkedGeoPolygon
	var next *LinkedGeoLoop

	innerCount, outerCount := 0, 0

	// Create an array to hold all of the inner loops. Note that
	// this array will never be full, as there will always be fewer
	// inner loops than outer loops.
	innerLoops := make([]*LinkedGeoLoop, loopCount)
	bboxes := make([]BBox, loopCount)

	// Get the first loop and unlink it from root
	loop := root.first

	*root = LinkedGeoPolygon{}

	// Iterate over all loops, moving inner loops into an array and
	// assigning outer loops to new polygons
	for loop != nil {
		if isClockwise(loop) {
			innerLoops[innerCount] = loop
			innerCount++
		} else {
			if polygon == nil {
				polygon = root
			} else {
				polygon = addNewLinkedPolygon(polygon)
			}
			addLinkedLoop(polygon, loop)
			bboxFrom(loop, &bboxes[outerCount])
			outerCount++
		}
		// get the next loop and unlink it from this one
		next = loop.next
		loop.next = nil
		loop = next
	}

	// Find polygon for each inner loop and assign the hole to it
	for i := 0; i < innerCount; i++ {
		polygon = findPolygonForHole(innerLoops[i], root, bboxes, outerCount)
		if polygon != nil {
			addLinkedLoop(polygon, innerLoops[i])
		} else {
			// If we can't find a polygon (possible with invalid input), then
			// we need to release the memory for the hole, because the loop has
			// been unlinked from the root and the caller will no longer have
			// a way to destroy it with destroyLinkedPolygon.
			destroyLinkedGeoLoop(innerLoops[i])
			innerLoops[i] = nil
			resultCode = NORMALIZATION_ERR_UNASSIGNED_HOLES
		}
	}

	// Free allocated memory
	innerLoops = nil
	bboxes = nil

	return resultCode
}
