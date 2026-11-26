package learnermodel

import "math"

// SecondsPerDay is the package's day under an exported name, for its external tests.
const SecondsPerDay = secondsPerDay

// NaN is the NaN that math.NaN gives, by its bits: the tests take math.NaN no more than the code
// does.
var NaN = math.Float64frombits(0x7FF8000000000001)
