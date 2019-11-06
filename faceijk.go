package h3

import (
	"fmt"
	"math"
)

/**
 * @brief Face number and ijk coordinates on that face-centered coordinate
 * system
 */
type FaceIJK struct {
	face  int      ///< face number
	coord CoordIJK ///< ijk coordinates on that face
}

/**
 * @brief Information to transform into an adjacent face IJK system
 */
type FaceOrientIJK struct {
	face      int      ///< face number
	translate CoordIJK ///< res 0 translation relative to primary face
	ccwRot60  int      ///< number of 60 degree ccw rotations relative to primary
}

// indexes for faceNeighbors table
/** IJ quadrant faceNeighbors table direction */
const IJ = 1

/** KI quadrant faceNeighbors table direction */
const KI = 2

/** JK quadrant faceNeighbors table direction */
const JK = 3

/** Invalid face index */
const INVALID_FACE = -1

type Overage int

/** Digit representing overage type */
const (
	/** No overage (on original face) */
	NO_OVERAGE Overage = 0
	/** On face edge (only occurs on substrate grids) */
	FACE_EDGE Overage = 1
	/** Overage on new face interior */
	NEW_FACE Overage = 2
)

/** square root of 7 */
const M_SQRT7 = 2.6457513110645905905016157536392604257102

/** @brief icosahedron face centers in Lat/Lon radians */
var faceCenterGeo = [NUM_ICOSA_FACES]GeoCoord{
	{0.803582649718989942, 1.248397419617396099},   // face  0
	{1.307747883455638156, 2.536945009877921159},   // face  1
	{1.054751253523952054, -1.347517358900396623},  // face  2
	{0.600191595538186799, -0.450603909469755746},  // face  3
	{0.491715428198773866, 0.401988202911306943},   // face  4
	{0.172745327415618701, 1.678146885280433686},   // face  5
	{0.605929321571350690, 2.953923329812411617},   // face  6
	{0.427370518328979641, -1.888876200336285401},  // face  7
	{-0.079066118549212831, -0.733429513380867741}, // face  8
	{-0.230961644455383637, 0.506495587332349035},  // face  9
	{0.079066118549212831, 2.408163140208925497},   // face 10
	{0.230961644455383637, -2.635097066257444203},  // face 11
	{-0.172745327415618701, -1.463445768309359553}, // face 12
	{-0.605929321571350690, -0.187669323777381622}, // face 13
	{-0.427370518328979641, 1.252716453253507838},  // face 14
	{-0.600191595538186799, 2.690988744120037492},  // face 15
	{-0.491715428198773866, -2.739604450678486295}, // face 16
	{-0.803582649718989942, -1.893195233972397139}, // face 17
	{-1.307747883455638156, -0.604647643711872080}, // face 18
	{-1.054751253523952054, 1.794075294689396615},  // face 19
}

/** @brief icosahedron face centers in x/y/z on the unit sphere */
var faceCenterPoint = [NUM_ICOSA_FACES]Vec3d{
	{0.2199307791404606, 0.6583691780274996, 0.7198475378926182},    // face  0
	{-0.2139234834501421, 0.1478171829550703, 0.9656017935214205},   // face  1
	{0.1092625278784797, -0.4811951572873210, 0.8697775121287253},   // face  2
	{0.7428567301586791, -0.3593941678278028, 0.5648005936517033},   // face  3
	{0.8112534709140969, 0.3448953237639384, 0.4721387736413930},    // face  4
	{-0.1055498149613921, 0.9794457296411413, 0.1718874610009365},   // face  5
	{-0.8075407579970092, 0.1533552485898818, 0.5695261994882688},   // face  6
	{-0.2846148069787907, -0.8644080972654206, 0.4144792552473539},  // face  7
	{0.7405621473854482, -0.6673299564565524, -0.0789837646326737},  // face  8
	{0.8512303986474293, 0.4722343788582681, -0.2289137388687808},   // face  9
	{-0.7405621473854481, 0.6673299564565524, 0.0789837646326737},   // face 10
	{-0.8512303986474292, -0.4722343788582682, 0.2289137388687808},  // face 11
	{0.1055498149613919, -0.9794457296411413, -0.1718874610009365},  // face 12
	{0.8075407579970092, -0.1533552485898819, -0.5695261994882688},  // face 13
	{0.2846148069787908, 0.8644080972654204, -0.4144792552473539},   // face 14
	{-0.7428567301586791, 0.3593941678278027, -0.5648005936517033},  // face 15
	{-0.8112534709140971, -0.3448953237639382, -0.4721387736413930}, // face 16
	{-0.2199307791404607, -0.6583691780274996, -0.7198475378926182}, // face 17
	{0.2139234834501420, -0.1478171829550704, -0.9656017935214205},  // face 18
	{-0.1092625278784796, 0.4811951572873210, -0.8697775121287253},  // face 19
}

/** @brief icosahedron face ijk axes as azimuth in radians from face center to
 * vertex 0/1/2 respectively
 */
var faceAxesAzRadsCII = [NUM_ICOSA_FACES][3]float64{
	{5.619958268523939882, 3.525563166130744542, 1.431168063737548730}, // face  0
	{5.760339081714187279, 3.665943979320991689, 1.571548876927796127}, // face  1
	{0.780213654393430055, 4.969003859179821079, 2.874608756786625655}, // face  2
	{0.430469363979999913, 4.619259568766391033, 2.524864466373195467}, // face  3
	{6.130269123335111400, 4.035874020941915804, 1.941478918548720291}, // face  4
	{2.692877706530642877, 0.598482604137447119, 4.787272808923838195}, // face  5
	{2.982963003477243874, 0.888567901084048369, 5.077358105870439581}, // face  6
	{3.532912002790141181, 1.438516900396945656, 5.627307105183336758}, // face  7
	{3.494305004259568154, 1.399909901866372864, 5.588700106652763840}, // face  8
	{3.003214169499538391, 0.908819067106342928, 5.097609271892733906}, // face  9
	{5.930472956509811562, 3.836077854116615875, 1.741682751723420374}, // face 10
	{0.138378484090254847, 4.327168688876645809, 2.232773586483450311}, // face 11
	{0.448714947059150361, 4.637505151845541521, 2.543110049452346120}, // face 12
	{0.158629650112549365, 4.347419854898940135, 2.253024752505744869}, // face 13
	{5.891865957979238535, 3.797470855586042958, 1.703075753192847583}, // face 14
	{2.711123289609793325, 0.616728187216597771, 4.805518392002988683}, // face 15
	{3.294508837434268316, 1.200113735041072948, 5.388903939827463911}, // face 16
	{3.804819692245439833, 1.710424589852244509, 5.899214794638635174}, // face 17
	{3.664438879055192436, 1.570043776661997111, 5.758833981448388027}, // face 18
	{2.361378999196363184, 0.266983896803167583, 4.455774101589558636}, // face 19
}

/** @brief Definition of which faces neighbor each other. */
var faceNeighbors = [NUM_ICOSA_FACES][4]FaceOrientIJK{
	{
		// face 0
		{0, CoordIJK{0, 0, 0}, 0}, // central face
		{4, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{1, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{5, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 1
		{1, CoordIJK{0, 0, 0}, 0}, // central face
		{0, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{2, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{6, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 2
		{2, CoordIJK{0, 0, 0}, 0}, // central face
		{1, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{3, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{7, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 3
		{3, CoordIJK{0, 0, 0}, 0}, // central face
		{2, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{4, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{8, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 4
		{4, CoordIJK{0, 0, 0}, 0}, // central face
		{3, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{0, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{9, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 5
		{5, CoordIJK{0, 0, 0}, 0},  // central face
		{10, CoordIJK{2, 2, 0}, 3}, // ij quadrant
		{14, CoordIJK{2, 0, 2}, 3}, // ki quadrant
		{0, CoordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 6
		{6, CoordIJK{0, 0, 0}, 0},  // central face
		{11, CoordIJK{2, 2, 0}, 3}, // ij quadrant
		{10, CoordIJK{2, 0, 2}, 3}, // ki quadrant
		{1, CoordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 7
		{7, CoordIJK{0, 0, 0}, 0},  // central face
		{12, CoordIJK{2, 2, 0}, 3}, // ij quadrant
		{11, CoordIJK{2, 0, 2}, 3}, // ki quadrant
		{2, CoordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 8
		{8, CoordIJK{0, 0, 0}, 0},  // central face
		{13, CoordIJK{2, 2, 0}, 3}, // ij quadrant
		{12, CoordIJK{2, 0, 2}, 3}, // ki quadrant
		{3, CoordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 9
		{9, CoordIJK{0, 0, 0}, 0},  // central face
		{14, CoordIJK{2, 2, 0}, 3}, // ij quadrant
		{13, CoordIJK{2, 0, 2}, 3}, // ki quadrant
		{4, CoordIJK{0, 2, 2}, 3},  // jk quadrant
	},
	{
		// face 10
		{10, CoordIJK{0, 0, 0}, 0}, // central face
		{5, CoordIJK{2, 2, 0}, 3},  // ij quadrant
		{6, CoordIJK{2, 0, 2}, 3},  // ki quadrant
		{15, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 11
		{11, CoordIJK{0, 0, 0}, 0}, // central face
		{6, CoordIJK{2, 2, 0}, 3},  // ij quadrant
		{7, CoordIJK{2, 0, 2}, 3},  // ki quadrant
		{16, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 12
		{12, CoordIJK{0, 0, 0}, 0}, // central face
		{7, CoordIJK{2, 2, 0}, 3},  // ij quadrant
		{8, CoordIJK{2, 0, 2}, 3},  // ki quadrant
		{17, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 13
		{13, CoordIJK{0, 0, 0}, 0}, // central face
		{8, CoordIJK{2, 2, 0}, 3},  // ij quadrant
		{9, CoordIJK{2, 0, 2}, 3},  // ki quadrant
		{18, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 14
		{14, CoordIJK{0, 0, 0}, 0}, // central face
		{9, CoordIJK{2, 2, 0}, 3},  // ij quadrant
		{5, CoordIJK{2, 0, 2}, 3},  // ki quadrant
		{19, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 15
		{15, CoordIJK{0, 0, 0}, 0}, // central face
		{16, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{19, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{10, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 16
		{16, CoordIJK{0, 0, 0}, 0}, // central face
		{17, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{15, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{11, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 17
		{17, CoordIJK{0, 0, 0}, 0}, // central face
		{18, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{16, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{12, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 18
		{18, CoordIJK{0, 0, 0}, 0}, // central face
		{19, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{17, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{13, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	},
	{
		// face 19
		{19, CoordIJK{0, 0, 0}, 0}, // central face
		{15, CoordIJK{2, 0, 2}, 1}, // ij quadrant
		{18, CoordIJK{2, 2, 0}, 5}, // ki quadrant
		{14, CoordIJK{0, 2, 2}, 3}, // jk quadrant
	}}

/** @brief direction from the origin face to the destination face, relative to
 * the origin face's coordinate system, or -1 if not adjacent.
 */
var adjacentFaceDir = [NUM_ICOSA_FACES][NUM_ICOSA_FACES]int{
	{0, KI, -1, -1, IJ, JK, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, // face 0
	{IJ, 0, KI, -1, -1, -1, JK, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, // face 1
	{-1, IJ, 0, KI, -1, -1, -1, JK, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, // face 2
	{-1, -1, IJ, 0, KI, -1, -1, -1, JK, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, // face 3
	{KI, -1, -1, IJ, 0, -1, -1, -1, -1, JK, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, // face 4
	{JK, -1, -1, -1, -1, 0, -1, -1, -1, -1, IJ, -1, -1, -1, KI, -1, -1, -1, -1, -1}, // face 5
	{-1, JK, -1, -1, -1, -1, 0, -1, -1, -1, KI, IJ, -1, -1, -1, -1, -1, -1, -1, -1}, // face 6
	{-1, -1, JK, -1, -1, -1, -1, 0, -1, -1, -1, KI, IJ, -1, -1, -1, -1, -1, -1, -1}, // face 7
	{-1, -1, -1, JK, -1, -1, -1, -1, 0, -1, -1, -1, KI, IJ, -1, -1, -1, -1, -1, -1}, // face 8
	{-1, -1, -1, -1, JK, -1, -1, -1, -1, 0, -1, -1, -1, KI, IJ, -1, -1, -1, -1, -1}, // face 9
	{-1, -1, -1, -1, -1, IJ, KI, -1, -1, -1, 0, -1, -1, -1, -1, JK, -1, -1, -1, -1}, // face 10
	{-1, -1, -1, -1, -1, -1, IJ, KI, -1, -1, -1, 0, -1, -1, -1, -1, JK, -1, -1, -1}, // face 11
	{-1, -1, -1, -1, -1, -1, -1, IJ, KI, -1, -1, -1, 0, -1, -1, -1, -1, JK, -1, -1}, // face 12
	{-1, -1, -1, -1, -1, -1, -1, -1, IJ, KI, -1, -1, -1, 0, -1, -1, -1, -1, JK, -1}, // face 13
	{-1, -1, -1, -1, -1, KI, -1, -1, -1, IJ, -1, -1, -1, -1, 0, -1, -1, -1, -1, JK}, // face 14
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, JK, -1, -1, -1, -1, 0, IJ, -1, -1, KI}, // face 15
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, JK, -1, -1, -1, KI, 0, IJ, -1, -1}, // face 16
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, JK, -1, -1, -1, KI, 0, IJ, -1}, // face 17
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, JK, -1, -1, -1, KI, 0, IJ}, // face 18
	{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, JK, IJ, -1, -1, KI, 0}, // face 19
}

/** @brief overage distance table */
var maxDimByCIIres = []int{
	2,        // res  0
	-1,       // res  1
	14,       // res  2
	-1,       // res  3
	98,       // res  4
	-1,       // res  5
	686,      // res  6
	-1,       // res  7
	4802,     // res  8
	-1,       // res  9
	33614,    // res 10
	-1,       // res 11
	235298,   // res 12
	-1,       // res 13
	1647086,  // res 14
	-1,       // res 15
	11529602, // res 16
}

/** @brief unit scale distance table */
var unitScaleByCIIres = []int{
	1,       // res  0
	-1,      // res  1
	7,       // res  2
	-1,      // res  3
	49,      // res  4
	-1,      // res  5
	343,     // res  6
	-1,      // res  7
	2401,    // res  8
	-1,      // res  9
	16807,   // res 10
	-1,      // res 11
	117649,  // res 12
	-1,      // res 13
	823543,  // res 14
	-1,      // res 15
	5764801, // res 16
}

/**
 * Encodes a coordinate on the sphere to the address *FaceIJK of the containing
 * cell at the specified resolution.
 *
 * @param g The spherical coordinates to encode.
 * @param res The desired H3 resolution for the encoding.
 * @param h The address *FaceIJK of the containing cell at resolution res.
 */
func _geoToFaceIjk(g *GeoCoord, res int, h *FaceIJK) {
	// first convert to hex2d
	var v Vec2d
	_geoToHex2d(g, res, &h.face, &v)

	_hex2dToCoordIJK(&v, &h.coord)
}

/**
 * Encodes a coordinate on the sphere to the corresponding icosahedral face and
 * containing 2D hex coordinates relative to that face center.
 *
 * @param g The spherical coordinates to encode.
 * @param res The desired H3 resolution for the encoding.
 * @param face The icosahedral face containing the spherical coordinates.
 * @param v The 2D hex coordinates of the cell containing the point.
 */
func _geoToHex2d(g *GeoCoord, res int, face *int, v *Vec2d) {
	var v3d Vec3d
	_geoToVec3d(g, &v3d)

	// determine the icosahedron face
	*face = 0
	sqd := _pointSquareDist(&faceCenterPoint[0], &v3d)
	for f := 1; f < NUM_ICOSA_FACES; f++ {
		sqdT := _pointSquareDist(&faceCenterPoint[f], &v3d)
		if sqdT < sqd {
			*face = f
			sqd = sqdT
		}
	}

	// cos(r) = 1 - 2 * sin^2(r/2) = 1 - 2 * (sqd / 4) = 1 - sqd/2
	r := math.Acos(1 - sqd/2)

	if r < EPSILON {
		v.x = 0.0
		v.y = 0.0
		return
	}

	// now have face and r, now find CCW theta from CII i-axis
	theta := _posAngleRads(faceAxesAzRadsCII[*face][0] - _posAngleRads(_geoAzimuthRads(&faceCenterGeo[*face], g)))

	// adjust theta for Class III (odd resolutions)
	if isResClassIII(res) {
		theta = _posAngleRads(theta - M_AP7_ROT_RADS)
	}

	// perform gnomonic scaling of r
	r = math.Tan(r)

	// scale for current resolution length u
	r /= RES0_U_GNOMONIC
	for i := 0; i < res; i++ {
		r *= M_SQRT7
	}

	// we now have (r, theta) in hex2d with theta ccw from x-axes

	// convert to local x,y
	v.x = r * math.Cos(theta)
	v.y = r * math.Sin(theta)
}

/**
 * Determines the center poin int spherical coordinates of a cell given by 2D
 * hex coordinates on a particular icosahedral face.
 *
 * @param v The 2D hex coordinates of the cell.
 * @param face The icosahedral face upon which the 2D hex coordinate system is
 *             centered.
 * @param res The H3 resolution of the cell.
 * @param substrate Indicates whether or not this grid is actually a substrate
 *        grid relative to the specified resolution.
 * @param g The spherical coordinates of the cell center point.
 */
func _hex2dToGeo(v *Vec2d, face int, res int, substrate int, g *GeoCoord) {
	// calculate (r, theta) in hex2d
	r := _v2dMag(v)

	if r < EPSILON {
		*g = faceCenterGeo[face]
		return
	}

	theta := math.Atan2(v.y, v.x)

	// scale for current resolution length u
	for i := 0; i < res; i++ {
		r /= M_SQRT7
	}

	// scale accordingly if this is a substrate grid
	if substrate != 0 {
		r /= 3.0
		if isResClassIII(res) {
			r /= M_SQRT7
		}
	}

	r *= RES0_U_GNOMONIC

	// perform inverse gnomonic scaling of r
	r = math.Atan(r)

	// adjust theta for Class III
	// if a substrate grid, then it's already been adjusted for Class III
	if substrate == 0 && isResClassIII(res) {
		theta = _posAngleRads(theta + M_AP7_ROT_RADS)
	}

	// find theta as an azimuth
	theta = _posAngleRads(faceAxesAzRadsCII[face][0] - theta)

	// now find the poat int (r,theta) from the face center
	_geoAzDistanceRads(&faceCenterGeo[face], theta, r, g)
}

/**
 * Determines the center poin int spherical coordinates of a cell given by
 * a address *FaceIJK at a specified resolution.
 *
 * @param h The address *FaceIJK of the cell.
 * @param res The H3 resolution of the cell.
 * @param g The spherical coordinates of the cell center point.
 */
func _faceIjkToGeo(h *FaceIJK, res int, g *GeoCoord) {
	var v Vec2d
	_ijkToHex2d(&h.coord, &v)
	_hex2dToGeo(&v, h.face, res, 0, g)
}

/**
 * Generates the cell boundary in spherical coordinates for a pentagonal cell
 * given by a address *FaceIJK at a specified resolution.
 *
 * @param h The address *FaceIJK of the pentagonal cell.
 * @param res The H3 resolution of the cell.
 * @param g The spherical coordinates of the cell boundary.
 */
func _faceIjkPentToGeoBoundary(h *FaceIJK, res int, g *GeoBoundary) {
	adjRes := res
	centerIJK := *h
	fijkVerts := make([]FaceIJK, NUM_PENT_VERTS)

	_faceIjkPentToVerts(&centerIJK, &adjRes, fijkVerts)

	// convert each vertex to Lat/Lon
	// adjust the face of each vertex as appropriate and introduce
	// edge-crossing vertices as needed
	g.numVerts = 0
	var lastFijk *FaceIJK

	for vert := 0; vert < NUM_PENT_VERTS+1; vert++ {

		v := vert % NUM_PENT_VERTS
		fijk := fijkVerts[v]

		_adjustPentVertOverage(&fijk, adjRes)

		// all Class III pentagon edges cross icosa edges
		// note that Class II pentagons have vertices on the edge,
		// not edge intersections
		if isResClassIII(res) && vert > 0 {
			// find hex2d of the two vertexes on the last face

			tmpFijk := fijk

			var orig2d0 Vec2d
			_ijkToHex2d(&lastFijk.coord, &orig2d0)

			currentToLastDir := adjacentFaceDir[tmpFijk.face][lastFijk.face]

			fijkOrient := &faceNeighbors[tmpFijk.face][currentToLastDir]

			tmpFijk.face = fijkOrient.face
			ijk := &tmpFijk.coord

			// rotate and translate for adjacent face
			for i := 0; i < fijkOrient.ccwRot60; i++ {
				_ijkRotate60ccw(ijk)
			}

			transVec := fijkOrient.translate
			_ijkScale(&transVec, unitScaleByCIIres[adjRes]*3)
			_ijkAdd(ijk, &transVec, ijk)
			_ijkNormalize(ijk)

			var orig2d1 Vec2d
			_ijkToHex2d(ijk, &orig2d1)

			// find the appropriate icosa face edge vertexes
			maxDim := float64(maxDimByCIIres[adjRes])
			v0 := Vec2d{3.0 * maxDim, 0.0}
			v1 := Vec2d{-1.5 * maxDim, 3.0 * M_SQRT3_2 * maxDim}
			v2 := Vec2d{-1.5 * maxDim, -3.0 * M_SQRT3_2 * maxDim}

			var edge0 *Vec2d
			var edge1 *Vec2d

			f := adjacentFaceDir[tmpFijk.face][fijk.face]
			switch f {
			case IJ:
				edge0 = &v0
				edge1 = &v1
				break
			case JK:
				edge0 = &v1
				edge1 = &v2
				break
			default:
				if f != KI {
					panic(fmt.Errorf("invalid direction. expect %d, but got %d", KI, f))
				}
				edge0 = &v2
				edge1 = &v0
				break
			}

			// find the intersection and add the Lat/Lon poto int the result
			var inter Vec2d
			_v2dIntersect(&orig2d0, &orig2d1, edge0, edge1, &inter)

			if len(g.Verts) <= g.numVerts {
				g.Verts = append(g.Verts, GeoCoord{})
			}
			_hex2dToGeo(&inter, tmpFijk.face, adjRes, 1, &g.Verts[g.numVerts])
			g.numVerts++
		}

		// convert vertex to Lat/Lon and add to the result
		// vert == NUM_PENT_VERTS is only used to test for possible intersection
		// on last edge
		if vert < NUM_PENT_VERTS {
			var vec Vec2d
			_ijkToHex2d(&fijk.coord, &vec)

			coord := GeoCoord{}
			_hex2dToGeo(&vec, fijk.face, adjRes, 1, &coord)

			if len(g.Verts) > g.numVerts {
				g.Verts[g.numVerts] = coord
			} else {
				g.Verts = append(g.Verts, coord)
			}
			g.numVerts++
		}

		lastFijk = &fijk
	}
}

/**
 * Get the vertices of a pentagon cell as substrate addresses *FaceIJK
 *
 * @param fijk The address *FaceIJK of the cell.
 * @param res The H3 resolution of the cell. This may be adjusted if
 *            necessary for the substrate grid resolution.
 * @param fijkVerts Output array for the vertices
 */
func _faceIjkPentToVerts(fijk *FaceIJK, res *int, fijkVerts []FaceIJK) {
	// get the correct set of substrate vertices for this resolution
	var verts []CoordIJK
	if isResClassIII(*res) {
		// the vertexes of an origin-centered pentagon in a Class III resolution on
		// a substrate grid with aperture sequence 33r7r. The aperture 3 gets us the
		// vertices, and the 3r7r gets us to Class II. vertices listed ccw from the
		// i-axes
		verts = []CoordIJK{
			{5, 4, 0}, // 0
			{1, 5, 0}, // 1
			{0, 5, 4}, // 2
			{0, 1, 5}, // 3
			{4, 0, 5}, // 4
		}
	} else {
		// the vertexes of an origin-centered pentagon in a Class II resolution on a
		// substrate grid with aperture sequence 33r. The aperture 3 gets us the
		// vertices, and the 3r gets us back to Class II.
		// vertices listed ccw from the i-axes
		verts = []CoordIJK{
			{2, 1, 0}, // 0
			{1, 2, 0}, // 1
			{0, 2, 1}, // 2
			{0, 1, 2}, // 3
			{1, 0, 2}, // 4
		}

	}

	// adjust the center poto int be in an aperture 33r substrate grid
	// these should be composed for speed
	_downAp3(&fijk.coord)
	_downAp3r(&fijk.coord)

	// if res is Class III we need to add a cw aperture 7 to get to
	// icosahedral Class II
	if isResClassIII(*res) {
		_downAp7r(&fijk.coord)
		*res += 1
	}

	// The center pois int now in the same substrate grid as the origin
	// cell vertices. Add the center posubstate int coordinates
	// to each vertex to translate the vertices to that cell.
	for v := 0; v < NUM_PENT_VERTS; v++ {
		fijkVerts[v].face = fijk.face
		_ijkAdd(&fijk.coord, &verts[v], &fijkVerts[v].coord)
		_ijkNormalize(&fijkVerts[v].coord)
	}
}

/**
 * Generates the cell boundary in spherical coordinates for a cell given by a
 * address *FaceIJK at a specified resolution.
 *
 * @param h The address *FaceIJK of the cell.
 * @param res The H3 resolution of the cell.
 * @param isPentagon Whether or not the cell is a pentagon.
 * @param g The spherical coordinates of the cell boundary.
 */
func _faceIjkToGeoBoundary(h *FaceIJK, res int, isPentagon bool, g *GeoBoundary) {
	if isPentagon {
		_faceIjkPentToGeoBoundary(h, res, g)
		return
	}

	centerIJK := *h

	adjRes := res

	fijkVerts := make([]FaceIJK, NUM_HEX_VERTS)

	_faceIjkToVerts(&centerIJK, &adjRes, fijkVerts)

	// convert each vertex to Lat/Lon
	// adjust the face of each vertex as appropriate and introduce
	// edge-crossing vertices as needed
	g.numVerts = 0
	lastFace := -1
	lastOverage := NO_OVERAGE

	for vert := 0; vert < NUM_HEX_VERTS+1; vert++ {
		v := vert % NUM_HEX_VERTS

		fijk := fijkVerts[v]

		pentLeading4 := 0
		overage := _adjustOverageClassII(&fijk, adjRes, pentLeading4, 1)

		/*
		   Check for edge-crossing. Each face of the underlying icosahedron is a
		   different projection plane. So if an edge of the hexagon crosses an
		   icosahedron edge, an additional vertex must be introduced at that
		   intersection point. Then each half of the cell edge can be projected
		   to geographic coordinates using the appropriate icosahedron face
		   projection. Note that Class II cell edges have vertices on the face
		   edge, with no edge line intersections.
		*/
		if isResClassIII(res) && vert > 0 && fijk.face != lastFace &&
			lastOverage != FACE_EDGE {
			// find hex2d of the two vertexes on original face
			lastV := (v + 5) % NUM_HEX_VERTS
			var orig2d0 Vec2d
			_ijkToHex2d(&fijkVerts[lastV].coord, &orig2d0)

			var orig2d1 Vec2d
			_ijkToHex2d(&fijkVerts[v].coord, &orig2d1)

			// find the appropriate icosa face edge vertexes
			maxDim := float64(maxDimByCIIres[adjRes])
			v0 := Vec2d{3.0 * maxDim, 0.0}
			v1 := Vec2d{-1.5 * maxDim, 3.0 * M_SQRT3_2 * maxDim}
			v2 := Vec2d{-1.5 * maxDim, -3.0 * M_SQRT3_2 * maxDim}

			face2 := lastFace
			if lastFace == centerIJK.face {
				face2 = fijk.face
			}

			var edge0 *Vec2d
			var edge1 *Vec2d

			f := adjacentFaceDir[centerIJK.face][face2]
			switch f {
			case IJ:
				edge0 = &v0
				edge1 = &v1
				break
			case JK:
				edge0 = &v1
				edge1 = &v2
				break
			default:
				if f != KI {
					panic(fmt.Errorf("invalid direction. expect %d, but got %d", KI, f))
				}

				edge0 = &v2
				edge1 = &v0
				break
			}

			// find the intersection and add the Lat/Lon poto int the result
			var inter Vec2d

			_v2dIntersect(&orig2d0, &orig2d1, edge0, edge1, &inter)

			/*
			   If a poof int intersection occurs at a hexagon vertex, then each
			   adjacent hexagon edge will lie completely on a single icosahedron
			   face, and no additional vertex is required.
			*/
			isIntersectionAtVertex := _v2dEquals(&orig2d0, &inter) || _v2dEquals(&orig2d1, &inter)
			if !isIntersectionAtVertex {
				if len(g.Verts) == g.numVerts {
					g.Verts = append(g.Verts, GeoCoord{})
				}
				_hex2dToGeo(&inter, centerIJK.face, adjRes, 1, &g.Verts[g.numVerts])
				g.numVerts++
			}
		}

		// convert vertex to Lat/Lon and add to the result
		// vert == NUM_HEX_VERTS is only used to test for possible intersection
		// on last edge
		if vert < NUM_HEX_VERTS {
			var vec Vec2d
			_ijkToHex2d(&fijk.coord, &vec)
			if len(g.Verts) == g.numVerts {
				g.Verts = append(g.Verts, GeoCoord{})
			}
			_hex2dToGeo(&vec, fijk.face, adjRes, 1, &g.Verts[g.numVerts])
			g.numVerts++
		}

		lastFace = fijk.face
		lastOverage = overage
	}
}

/**
 * Get the vertices of a cell as substrate addresses *FaceIJK
 *
 * @param fijk The address *FaceIJK of the cell.
 * @param res The H3 resolution of the cell. This may be adjusted if
 *            necessary for the substrate grid resolution.
 * @param fijkVerts Output array for the vertices
 */
func _faceIjkToVerts(fijk *FaceIJK, res *int, fijkVerts []FaceIJK) {
	// get the correct set of substrate vertices for this resolution
	var verts []CoordIJK

	if isResClassIII(*res) {
		// the vertexes of an origin-centered cell in a Class III resolution on a
		// substrate grid with aperture sequence 33r7r. The aperture 3 gets us the
		// vertices, and the 3r7r gets us to Class II.
		// vertices listed ccw from the i-axes
		verts = []CoordIJK{
			{5, 4, 0}, // 0
			{1, 5, 0}, // 1
			{0, 5, 4}, // 2
			{0, 1, 5}, // 3
			{4, 0, 5}, // 4
			{5, 0, 1}, // 5
		}
	} else {
		// the vertexes of an origin-centered cell in a Class II resolution on a
		// substrate grid with aperture sequence 33r. The aperture 3 gets us the
		// vertices, and the 3r gets us back to Class II.
		// vertices listed ccw from the i-axes
		verts = []CoordIJK{
			{2, 1, 0}, // 0
			{1, 2, 0}, // 1
			{0, 2, 1}, // 2
			{0, 1, 2}, // 3
			{1, 0, 2}, // 4
			{2, 0, 1}, // 5
		}
	}

	// adjust the center poto int be in an aperture 33r substrate grid
	// these should be composed for speed
	_downAp3(&fijk.coord)
	_downAp3r(&fijk.coord)

	// if res is Class III we need to add a cw aperture 7 to get to
	// icosahedral Class II
	if isResClassIII(*res) {
		_downAp7r(&fijk.coord)
		*res += 1
	}

	// The center pois int now in the same substrate grid as the origin
	// cell vertices. Add the center posubstate int coordinates
	// to each vertex to translate the vertices to that cell.
	for v := 0; v < NUM_HEX_VERTS; v++ {
		fijkVerts[v].face = fijk.face
		_ijkAdd(&fijk.coord, &verts[v], &fijkVerts[v].coord)
		_ijkNormalize(&fijkVerts[v].coord)
	}
}

/**
 * Adjusts a address *FaceIJK in place so that the resulting cell address is
 * relative to the correct icosahedral face.
 *
 * @param fijk The address *FaceIJK of the cell.
 * @param res The H3 resolution of the cell.
 * @param pentLeading4 Whether or not the cell is a pentagon with a leading
 *        digit 4.
 * @param substrate Whether or not the cell is in a substrate grid.
 * @return 0 if on original face (no overage); 1 if on face edge (only occurs
 *         on substrate grids); 2 if overage on new face interior
 */
func _adjustOverageClassII(fijk *FaceIJK, res int, pentLeading4 int, substrate int) Overage {
	overage := NO_OVERAGE

	ijk := &fijk.coord

	// get the maximum dimension value; scale if a substrate grid
	maxDim := maxDimByCIIres[res]
	if substrate != 0 {
		maxDim *= 3
	}

	// check for overage
	if substrate != 0 && ijk.i+ijk.j+ijk.k == maxDim {
		// on edge
		overage = FACE_EDGE
	} else if ijk.i+ijk.j+ijk.k > maxDim {
		// overage
		overage = NEW_FACE

		var fijkOrient *FaceOrientIJK
		if ijk.k > 0 {
			if ijk.j > 0 {
				// jk "quadrant"
				fijkOrient = &faceNeighbors[fijk.face][JK]
			} else {
				// ik "quadrant"
				fijkOrient = &faceNeighbors[fijk.face][KI]

				// adjust for the pentagonal missing sequence
				if pentLeading4 != 0 {
					// translate origin to center of pentagon
					var origin CoordIJK
					_setIJK(&origin, maxDim, 0, 0)

					var tmp CoordIJK
					_ijkSub(ijk, &origin, &tmp)
					// rotate to adjust for the missing sequence
					_ijkRotate60cw(&tmp)
					// translate the origin back to the center of the triangle
					_ijkAdd(&tmp, &origin, ijk)
				}
			}
		} else {
			// ij "quadrant"
			fijkOrient = &faceNeighbors[fijk.face][IJ]
		}

		fijk.face = fijkOrient.face

		// rotate and translate for adjacent face
		for i := 0; i < fijkOrient.ccwRot60; i++ {
			_ijkRotate60ccw(ijk)
		}

		transVec := fijkOrient.translate
		unitScale := unitScaleByCIIres[res]
		if substrate != 0 {
			unitScale *= 3
		}
		_ijkScale(&transVec, unitScale)
		_ijkAdd(ijk, &transVec, ijk)
		_ijkNormalize(ijk)

		// overage points on pentagon boundaries can end up on edges
		if substrate != 0 && ijk.i+ijk.j+ijk.k == maxDim {
			// on edge
			overage = FACE_EDGE
		}
	}

	return overage
}

/**
 * Adjusts a address *FaceIJK for a pentagon vertex in a substrate grid in
 * place so that the resulting cell address is relative to the correct
 * icosahedral face.
 *
 * @param fijk The address *FaceIJK of the cell.
 * @param res The H3 resolution of the cell.
 */
func _adjustPentVertOverage(fijk *FaceIJK, res int) Overage {
	pentLeading4 := 0
	var overage Overage
	for {
		overage = _adjustOverageClassII(fijk, res, pentLeading4, 1)

		if !(overage == NEW_FACE) {
			break
		}
	}

	return overage
}
