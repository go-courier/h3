package h3

import (
	"math"
)

/** pi */
const M_PI = math.Pi

/** pi / 2.0 */
const M_PI_2 = math.Pi / 2

/** 2.0 * PI */
const M_2PI = math.Pi * 2

/** pi / 180 */
const M_PI_180 = math.Pi / 180

/** 180 / pi  */
const M_180_PI = 180 / math.Pi

/** threshold epsilon */
var EPSILON = float64(7.)/3 - float64(4.)/3 - float64(1.)

/** sqrt(3) / 2.0 */
var M_SQRT3_2 = math.Sqrt(3) / 2

/** sin(60') */
var M_SIN60 = M_SQRT3_2

/** rotation angle between Class II and Class III resolution axes
 * (asin(sqrt(3.0 / 28.0))) */
var M_AP7_ROT_RADS = math.Asin(math.Sqrt(3.0 / 28.0))

/** sin(M_AP7_ROT_RADS) */
var M_SIN_AP7_ROT = math.Sin(M_AP7_ROT_RADS)

/** cos(M_AP7_ROT_RADS) */
var M_COS_AP7_ROT = math.Cos(M_AP7_ROT_RADS)

/** earth radius in kilometers using WGS84 authalic radius */
const EARTH_RADIUS_KM = 6371.007180918475

/** scaling factor from hex2d resolution 0 unit length
 * (or distance between adjacent cell center points
 * on the plane) to gnomonic unit length. */
const RES0_U_GNOMONIC = 0.38196601125010500003

/** max H3 resolution; H3 version 1 has 16 resolutions, numbered 0 through 15 */
const MAX_H3_RES = 15

/** The number of faces on an icosahedron */
const NUM_ICOSA_FACES = 20

/** The number of H3 base cells */
const NUM_BASE_CELLS = 122

/** The number of vertices in a hexagon */
const NUM_HEX_VERTS = 6

/** The number of vertices in a pentagon */
const NUM_PENT_VERTS = 5

/** The number of pentagons per resolution **/
const NUM_PENTAGONS = 12

type H3Mode int

/** H3 index modes */
const (
	H3_HEXAGON_MODE H3Mode = 1
	H3_UNIEDGE_MODE H3Mode = 2
)
