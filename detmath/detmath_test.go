package detmath

import (
	"math"
	"math/rand/v2"
	"runtime"
	"testing"
)

// Where math.Exp and math.Log run the same pure Go as the port, it must give their bits; elsewhere
// they run assembler of their own, and the two may differ in the last bits: both are within 1 ulp
// of e**x, so within 2 of each other (docs/architecture.md, section 3.6, what experience showed).
func sameCodeAsMath() bool { return runtime.GOARCH == "wasm" || runtime.GOARCH == "386" }

const mathULPs = 2

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
		0, math.Copysign(0, -1), math.Inf(1), math.Inf(-1), nan,
		math.Float64frombits(0xFFF0000000000001), math.SmallestNonzeroFloat64, 0x1p-1022,
		math.MaxFloat64, -math.MaxFloat64, 1, -1, 0x1p-28, -0x1p-28, below(0x1p-28),
		709.782712893384, 709.79, -745.1332191019411, -745.14, math.Sqrt2 / 2,
		below(math.Sqrt2 / 2),
	}
	r := rand.New(rand.NewPCG(3, seed))
	for range 100_000 {
		u := r.Float64()
		in = append(in, lo+float64((hi-lo)*u), math.Float64frombits(r.Uint64()))
	}
	return in
}

func TestExpAgainstMath(t *testing.T) {
	for _, x := range portInputs(1, -746, 710) {
		againstMath(t, "exp", x, exp(x), math.Exp(x))
	}
}

func TestLogAgainstMath(t *testing.T) {
	for _, x := range portInputs(2, 0, 1e6) {
		// math.Log's assembler on amd64 takes a subnormal for one near 2**-1023: log(5e-324) is
		// -709.09 there, not -744.44. The port is right, and Log keeps subnormals out anyway.
		if !sameCodeAsMath() && x != 0 && math.Abs(x) < 0x1p-1022 {
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
	case !math.IsNaN(got) && got != want && math.Abs(x) < 709:
		t.Fatalf("%s(%v) = %v, math gives %v", name, x, got, want)
	}
}

// The functions built on the port, against the same formulas over math.Exp and math.Log: σ within
// a few ulp, the power within the error that a logarithm a few ulp off carries into the exponent,
// some |s · ln a| ulp, and the retention, which takes nothing of the port, bit for bit.
func TestBuiltOnThePort(t *testing.T) {
	r := rand.New(rand.NewPCG(4, 0))
	for range 100_000 {
		x := float64(r.Float64()*1400) - 700
		if got, want := Sigmoid(x), float64(1/(1+math.Exp(-x))); ulps(got, want) > 4 {
			t.Fatalf("Sigmoid(%v) = %v, %d ulp from %v", x, got, ulps(got, want), want)
		}
		a, s := float64(r.Float64()*20), float64(r.Float64()*4)-2
		y := float64(s * math.Log(a))
		want := math.Exp(y)
		bound := float64(float64(4+float64(4*math.Abs(y))) * float64(0x1p-52*want))
		if got := Pow(a, s); math.Abs(got-want) > bound {
			t.Fatalf("Pow(%v, %v) = %v, against %v", a, s, got, want)
		}
		dt, tau := float64(r.Float64()*3650), 1+float64(r.Float64()*3650)
		if got, want := Retention(dt, tau), float64(1/math.Sqrt(1+float64(RetentionFactor*dt/tau))); got != want {
			t.Fatalf("Retention(%v, %v) = %v, not %v", dt, tau, got, want)
		}
	}
}

// The NaN that math.NaN gives, and the float just below a positive one: rule 2 keeps tests to the
// functions of math with an exact result, and these need none.
var nan = math.Float64frombits(uvnan)

func below(x float64) float64 { return math.Float64frombits(math.Float64bits(x) - 1) }

func TestBounds(t *testing.T) {
	inf := math.Inf(1)
	for _, c := range []struct {
		name      string
		got, want float64
	}{
		{"Exp(NaN)", Exp(nan), Exp(-708)},
		{"Exp(-Inf)", Exp(-inf), Exp(-708)},
		{"Exp(+Inf)", Exp(inf), Exp(708)},
		{"Exp(0)", Exp(0), 1},
		{"Log(0)", Log(0), Log(0x1p-1022)},
		{"Log(-1)", Log(-1), Log(0x1p-1022)},
		{"Log(NaN)", Log(nan), Log(0x1p-1022)},
		{"Log(subnormal)", Log(math.SmallestNonzeroFloat64), Log(0x1p-1022)},
		{"Log(+Inf)", Log(inf), Log(math.MaxFloat64)},
		{"Log(1)", Log(1), 0},
		{"Sigmoid(0)", Sigmoid(0), 0.5},
		{"Sigmoid(-1000)", Sigmoid(-1000), Sigmoid(-708)},
		{"Sigmoid(37)", Sigmoid(37), 1},
		{"Pow(a, 0)", Pow(2.5, 0), 1},
		{"Pow(a, 1)", Pow(2.5, 1), 2.5},
		{"Pow(1/2, 1)", Pow(0.5, 1), 0.5},
		{"Pow(0, 1)", Pow(0, 1), 0x1p-1022},
		{"Pow(NaN, 0)", Pow(nan, 0), 1},
		{"Pow(4, 1/2)", Pow(4, 0.5), 2},
		{"Retention(0, τ)", Retention(0, 10), 1},
		{"Retention(-1, τ)", Retention(-1, 10), 1},
		{"Retention(NaN, τ)", Retention(nan, 10), 1},
		{"Retention(0, 0)", Retention(0, 0), 1},
		{"Retention(1, -10)", Retention(1, -10), Retention(1, 0)},
		{"Retention(1, NaN)", Retention(1, nan), Retention(1, 0)},
		{"Retention(-1, -10)", Retention(-1, -10), 1},
		{"Sigmoid(NaN)", Sigmoid(nan), Sigmoid(-708)},
		{"Sigmoid(+Inf)", Sigmoid(inf), Sigmoid(708)},
		{"Pow(1, +Inf)", Pow(1, inf), 1},
		{"Pow(1, -Inf)", Pow(1, -inf), 1},
		{"Pow(1, NaN)", Pow(1, nan), 1},
		{"Pow(2, NaN)", Pow(2, nan), Exp(-708)},
		{"Pow(2, +Inf)", Pow(2, inf), Exp(708)},
		{"Retention(τ, τ)", Retention(10, 10), 0.8999999999999999},
	} {
		if math.Float64bits(c.got) != math.Float64bits(c.want) {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}
	for _, r := range []float64{Retention(inf, 1), Retention(1, 0), Retention(inf, 0)} {
		if r < 0x1p-1022 || r >= 1e-150 {
			t.Errorf("the retention at the largest pause = %v: not a normal number near 0", r)
		}
	}
	if got := Sigmoid(-708); got < 0x1p-1022 {
		t.Errorf("Sigmoid(-708) = %v: subnormal", got)
	}
}

func TestTruncate(t *testing.T) {
	for _, c := range []struct {
		x    float64
		want int
	}{
		{0, 0}, {1.9, 1}, {-1.9, -1}, {1074.5, 1074}, {-1075.5, -1075},
		{1 << 30, 1 << 30}, {1e300, 1 << 30}, {-1e300, -(1 << 30)}, {math.Inf(-1), -(1 << 30)},
		{nan, 0},
	} {
		if got := truncate(c.x); got != c.want {
			t.Errorf("truncate(%v) = %d, want %d", c.x, got, c.want)
		}
	}
}
