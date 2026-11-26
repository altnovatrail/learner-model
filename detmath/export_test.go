package detmath

import "math"

// The bounds and thresholds of the package under exported names, for its external tests.
const (
	MinNormal    = minNormal
	ExpMin       = expMin
	ExpMax       = expMax
	ExpOverflow  = expOverflow
	ExpUnderflow = expUnderflow
	ExpNearZero  = expNearZero
)

// NaN is the NaN that math.NaN gives, and OtherNaN one unlike it — negative, signalling, with a
// payload of 1 —, so that a test holds every NaN to the bounds, not math's alone. Both are taken
// by their bits, and Below by its bits too: the tests take math.NaN and math.Nextafter no more
// than the code does.
var (
	NaN      = math.Float64frombits(uvnan)
	OtherNaN = math.Float64frombits(0xFFF0000000000001)
)

// Below is the float just below a positive one.
func Below(x float64) float64 { return math.Float64frombits(math.Float64bits(x) - 1) }
