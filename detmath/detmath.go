// Package detmath is the arithmetic of the learner model that package math would give differently
// on different builds: exp, log, the logistic function σ, a power and the retention of a skill. It
// gives the same bits on every build of the gc compiler — every architecture, GOAMD64 level and
// processor, WebAssembly included — so that the state the model folds from the log is the same
// wherever it is computed (R09 in docs/decisions.md; docs/architecture.md, section 3.6).
//
// exp and log are a port of the pure Go code of package math (port.go), which math.Exp and
// math.Log run on wasm and 386 but not where they have assembler of their own. Every function here
// keeps its argument within bounds, so that its result is finite and never subnormal: no NaN and
// no infinity comes out of it, whatever comes in. The engine may take this package as it is.
package detmath

import "math"

// The bounds of Exp's argument: exp(±708) is finite and normal, and so is σ at ±708.
const (
	expMin = -708.0
	expMax = 708.0
)

// RetentionFactor is F in the retention R = 1 / √(1 + F · Δt / τ): 0.9^(−1/C) − 1 at C = 1/2,
// 19/81, so that the retention after a pause of τ is 0.9 (docs/architecture.md, section 3.2).
const RetentionFactor = 19.0 / 81

// Exp returns e**x for x within [-708, 708], and the value at the nearer bound outside it; NaN
// gives the value at -708.
func Exp(x float64) float64 {
	return exp(clamp(x, expMin, expMax))
}

// Log returns the natural logarithm of x for x within the positive normal numbers, from
// 2**-1022 to the largest float64, and the value at the nearer bound outside them: 0, a negative
// x, a subnormal and NaN give log(2**-1022), +Inf gives the log of the largest float64.
func Log(x float64) float64 {
	return log(clamp(x, minNormal, math.MaxFloat64))
}

// minNormal is the smallest positive normal float64, 2**-1022.
const minNormal = 0x1p-1022

// Sigmoid returns the logistic function σ(x) = 1 / (1 + e**(−x)) for x within the bounds of Exp,
// and the value at the nearer bound outside them; NaN gives the value at -708. It lies within
// [3.3e-308, 1], and is exactly 1 from x ≈ 36.7 up.
func Sigmoid(x float64) float64 {
	x = clamp(x, expMin, expMax)
	return float64(1 / (1 + exp(-x)))
}

// Pow returns a**s as e**(s · ln a), with a kept within the bounds of Log, s within the finite
// floats and s · ln a within the bounds of Exp; NaN goes to the lower bound of each. Its two exact
// branches keep the rules built on it exact at their ends: a**0 is 1 and a**1 is a
// (docs/architecture.md, section 3.6, rule 3). With both kept finite, s · ln a is never 0 · ∞:
// 1**s is 1 for every s.
func Pow(a, s float64) float64 {
	a = clamp(a, minNormal, math.MaxFloat64)
	s = clamp(s, -math.MaxFloat64, math.MaxFloat64)
	switch s {
	case 0:
		return 1
	case 1:
		return a
	}
	return Exp(float64(s * log(a)))
}

// Retention returns R(Δt, τ) = 1 / √(1 + F · Δt / τ), the retention of a skill dt days after its
// last practice, at stability tau days (docs/architecture.md, section 3.2). dt is kept within
// [0, the largest float64], so that a negative dt — a time that came back — and NaN give the
// retention right after practice, 1; tau is kept within the positive normal numbers, so that 0, a
// negative tau and NaN give the stability of 2**-1022 days, and the retention of any pause after
// practice almost 0; and F · Δt / τ is kept below the largest float64, so that the result is
// normal.
func Retention(dt, tau float64) float64 {
	dt = clamp(dt, 0, math.MaxFloat64)
	tau = clamp(tau, minNormal, math.MaxFloat64)
	x := clamp(float64(RetentionFactor*dt/tau), 0, math.MaxFloat64)
	return float64(1 / math.Sqrt(1+x))
}

// clamp keeps x within [lo, hi]; NaN goes to lo.
func clamp(x, lo, hi float64) float64 {
	switch {
	case x >= lo && x <= hi:
		return x
	case x > hi:
		return hi
	}
	return lo
}

// truncate converts x to an int toward zero, as int(x) does, but only within ±2**30, where every
// build converts alike; outside it, it saturates, and NaN gives 0 (docs/architecture.md, section
// 3.6, rule 5). The port calls it with |x| < 1100.
func truncate(x float64) int {
	const limit = 1 << 30
	switch {
	case x >= -limit && x <= limit:
		return int(x)
	case x > limit:
		return limit
	case x < -limit:
		return -limit
	}
	return 0
}
