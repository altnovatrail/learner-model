package learnermodel_test

// The tolerances of the external tests: sumTolerance, of a float sum of a few terms against its
// exact value; threeDecimals, of a value a document prints to three decimals.
const (
	sumTolerance  = 1e-12
	threeDecimals = 1e-3
)

// epoch is where the synthetic logs of the tests start: 2025-11-02, 00:00 UTC, in whole seconds.
const epoch = 1762041600
