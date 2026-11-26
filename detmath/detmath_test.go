package detmath

import (
	"math"
	"math/rand/v2"
	"runtime"
	"testing"
)

// Where math.Exp and math.Log run the same pure Go as the port, it must give their bits; elsewhere
// they run assembler of their own, and the two may differ in the last bits: both are within 1 ulp
// of e**x, so within 2 of each other.
func sameCodeAsMath() bool { return runtime.GOARCH == "wasm" || runtime.GOARCH == "386" }

const mathULPs = 2

// epsilon is 2**-52, the spacing of the floats at 1: x · epsilon is about one ulp of x.
const epsilon = 0x1p-52

// nearEdges is the |x| from which exp's argument is at the edges of overflow and underflow.
const nearEdges = 709

// The values just past exp's thresholds, where it gives +Inf and 0.
const (
	aboveOverflow  = 709.79
	belowUnderflow = -745.14
)

// randomInputs is how many random values a test of the port draws.
const randomInputs = 100_000

// The seeds of the random values: portSeed is the first word of every test of the port against
// math, and the other word is each test's own.
const (
	portSeed  = 3
	expSeed   = 1
	logSeed   = 2
	builtSeed = 4
)

// logHigh is the top of the range of log's random arguments, from 0.
const logHigh = 1e6

// ulps is the distance between two floats of one sign in units in the last place.
func ulps(a, b float64) uint64 {
	x, y := math.Float64bits(a), math.Float64bits(b)
	if x > y {
		return x - y
	}
	return y - x
}

// portInputs gives the special values and a run of values over the whole domain and as raw bits.
func portInputs(seed uint64, lo, hi float64) []float64 {
	in := []float64{
		0, math.Copysign(0, -1), math.Inf(1), math.Inf(-1), NaN, OtherNaN,
		math.SmallestNonzeroFloat64, minNormal, math.MaxFloat64, -math.MaxFloat64, 1, -1,
		expNearZero, -expNearZero, Below(expNearZero),
		expOverflow, aboveOverflow, expUnderflow, belowUnderflow, math.Sqrt2 / 2, Below(math.Sqrt2 / 2),
	}
	r := rand.New(rand.NewPCG(portSeed, seed))
	for range randomInputs {
		u := r.Float64()
		in = append(in, lo+float64((hi-lo)*u), math.Float64frombits(r.Uint64()))
	}
	return in
}

func TestExpAgainstMath(t *testing.T) {
	// The whole range where exp is finite and nonzero, and a little past it on each side.
	for _, x := range portInputs(expSeed, math.Floor(expUnderflow), math.Ceil(expOverflow)) {
		againstMath(t, "exp", x, exp(x), math.Exp(x))
	}
}

func TestLogAgainstMath(t *testing.T) {
	for _, x := range portInputs(logSeed, 0, logHigh) {
		// math.Log's assembler on amd64 takes a subnormal for one near 2**-1023: log(5e-324) is
		// -709.09 there, not -744.44. The port is right, and Log keeps subnormals out anyway.
		if !sameCodeAsMath() && x != 0 && math.Abs(x) < minNormal {
			continue
		}
		againstMath(t, "log", x, log(x), math.Log(x))
	}
}

// againstMath holds the port's result to math's: its bits where both run the same code. Elsewhere
// a NaN stays a NaN, though its payload is math's own choice, and a finite nonzero result is within
// mathULPs; at the edges of overflow and underflow, within 1 ulp of them, math's assembler may round
// to the infinity or to zero where the port gives the last finite value, and the port is not held
// to it there: Exp keeps its argument well inside.
func againstMath(t *testing.T, name string, x, got, want float64) {
	t.Helper()
	finite := func(y float64) bool { return !math.IsNaN(y) && !math.IsInf(y, 0) && y != 0 }
	switch {
	case sameCodeAsMath():
		if math.Float64bits(got) != math.Float64bits(want) {
			t.Fatalf("%s(%v) = %v (%016x), math gives %v (%016x)",
				name, x, got, math.Float64bits(got), want, math.Float64bits(want))
		}
	case math.IsNaN(got) != math.IsNaN(want):
		t.Fatalf("%s(%v) = %v, math gives %v", name, x, got, want)
	case finite(got) && finite(want):
		if ulps(got, want) > mathULPs {
			t.Fatalf("%s(%v) = %v, %d ulp from math's %v", name, x, got, ulps(got, want), want)
		}
	case !math.IsNaN(got) && got != want && math.Abs(x) < nearEdges:
		t.Fatalf("%s(%v) = %v, math gives %v", name, x, got, want)
	}
}

// The functions built on the port, against the same formulas over math.Exp and math.Log: σ within
// a few ulp, the power within the error that a logarithm a few ulp off carries into the exponent,
// some |s · ln a| ulp, and the retention, which takes nothing of the port, bit for bit.
func TestBuiltOnThePort(t *testing.T) {
	const (
		sigmoidReach  = 700  // σ's arguments are within ±sigmoidReach
		sigmoidULPs   = 4    // σ's ulp from the formula over math.Exp
		powBaseHigh   = 20   // the bases of the power are within [0, powBaseHigh)
		powPowerReach = 2    // its exponents within ±powPowerReach
		powULPs       = 4    // its ulp, and as many again for each unit of |s · ln a|
		tenYears      = 3650 // the pauses and the stabilities run up to ten years, in days
	)
	r := rand.New(rand.NewPCG(builtSeed, 0))
	for range randomInputs {
		x := float64(r.Float64()*(2*sigmoidReach)) - sigmoidReach
		if got, want := Sigmoid(x), float64(1/(1+math.Exp(-x))); ulps(got, want) > sigmoidULPs {
			t.Fatalf("Sigmoid(%v) = %v, %d ulp from %v", x, got, ulps(got, want), want)
		}
		a, s := float64(r.Float64()*powBaseHigh), float64(r.Float64()*(2*powPowerReach))-powPowerReach
		y := float64(s * math.Log(a))
		want := math.Exp(y)
		bound := float64(float64(powULPs+float64(powULPs*math.Abs(y))) * float64(epsilon*want))
		if got := Pow(a, s); math.Abs(got-want) > bound {
			t.Fatalf("Pow(%v, %v) = %v, against %v", a, s, got, want)
		}
		dt, tau := float64(r.Float64()*tenYears), 1+float64(r.Float64()*tenYears)
		if got, want := Retention(dt, tau), float64(1/math.Sqrt(1+float64(RetentionFactor*dt/tau))); got != want {
			t.Fatalf("Retention(%v, %v) = %v, not %v", dt, tau, got, want)
		}
	}
}

func TestBounds(t *testing.T) {
	const (
		a              = 2.5                // a base of the power
		tau            = 10                 // a stability, in days
		retentionAtTau = 0.8999999999999999 // R(τ, τ): 0.9 as the floats round it, an ulp below
		sigmoidOne     = 37                 // from about 36.7 up σ is exactly 1
		farBelowExp    = -1000              // far below the bounds of Exp
		nearZero       = 1e-150             // the retention at the largest pause is below it
	)
	inf := math.Inf(1)
	for _, c := range []struct {
		name      string
		got, want float64
	}{
		{"Exp(NaN)", Exp(NaN), Exp(expMin)},
		{"Exp(-Inf)", Exp(-inf), Exp(expMin)},
		{"Exp(+Inf)", Exp(inf), Exp(expMax)},
		{"Exp(0)", Exp(0), 1},
		{"Log(0)", Log(0), Log(minNormal)},
		{"Log(-1)", Log(-1), Log(minNormal)},
		{"Log(NaN)", Log(NaN), Log(minNormal)},
		{"Log(subnormal)", Log(math.SmallestNonzeroFloat64), Log(minNormal)},
		{"Log(+Inf)", Log(inf), Log(math.MaxFloat64)},
		{"Log(1)", Log(1), 0},
		{"Sigmoid(0)", Sigmoid(0), 0.5},
		{"Sigmoid(-1000)", Sigmoid(farBelowExp), Sigmoid(expMin)},
		{"Sigmoid(37)", Sigmoid(sigmoidOne), 1},
		{"Pow(a, 0)", Pow(a, 0), 1},
		{"Pow(a, 1)", Pow(a, 1), a},
		{"Pow(1/2, 1)", Pow(0.5, 1), 0.5},
		{"Pow(0, 1)", Pow(0, 1), minNormal},
		{"Pow(NaN, 0)", Pow(NaN, 0), 1},
		{"Pow(4, 1/2)", Pow(4, 0.5), 2},
		{"Retention(0, τ)", Retention(0, tau), 1},
		{"Retention(-1, τ)", Retention(-1, tau), 1},
		{"Retention(NaN, τ)", Retention(NaN, tau), 1},
		{"Retention(0, 0)", Retention(0, 0), 1},
		{"Retention(1, -τ)", Retention(1, -tau), Retention(1, 0)},
		{"Retention(1, NaN)", Retention(1, NaN), Retention(1, 0)},
		{"Retention(-1, -τ)", Retention(-1, -tau), 1},
		{"Sigmoid(NaN)", Sigmoid(NaN), Sigmoid(expMin)},
		{"Sigmoid(+Inf)", Sigmoid(inf), Sigmoid(expMax)},
		{"Pow(1, +Inf)", Pow(1, inf), 1},
		{"Pow(1, -Inf)", Pow(1, -inf), 1},
		{"Pow(1, NaN)", Pow(1, NaN), 1},
		{"Pow(2, NaN)", Pow(2, NaN), Exp(expMin)},
		{"Pow(2, +Inf)", Pow(2, inf), Exp(expMax)},
		{"Retention(τ, τ)", Retention(tau, tau), retentionAtTau},
	} {
		if math.Float64bits(c.got) != math.Float64bits(c.want) {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}
	for _, r := range []float64{Retention(inf, 1), Retention(1, 0), Retention(inf, 0)} {
		if r < minNormal || r >= nearZero {
			t.Errorf("the retention at the largest pause = %v: not a normal number near 0", r)
		}
	}
	if got := Sigmoid(expMin); got < minNormal {
		t.Errorf("Sigmoid(expMin) = %v: subnormal", got)
	}
}

func TestTruncate(t *testing.T) {
	for _, c := range []struct {
		x    float64
		want int
	}{
		{0, 0}, {1.9, 1}, {-1.9, -1}, {1074.5, 1074}, {-1075.5, -1075},
		{truncateLimit, truncateLimit}, {1e300, truncateLimit}, {-1e300, -truncateLimit},
		{math.Inf(-1), -truncateLimit}, {NaN, 0},
	} {
		if got := truncate(c.x); got != c.want {
			t.Errorf("truncate(%v) = %d, want %d", c.x, got, c.want)
		}
	}
}
