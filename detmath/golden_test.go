package detmath_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/altnovatrail/learner-model/detmath"
)

// The golden files of detmath, one for every build: they tell a difference in this arithmetic
// from a difference in the model. testdata/golden.txt holds inputs and the bits of the results;
// testdata/sweep.txt holds, for each function, the SHA-256 of the bits of its results over a long
// run of inputs drawn by PCG with fixed seeds, the same on every build. Both change only with a
// version of learner-model: DETMATH_UPDATE=1 writes them anew, from the build it runs on.

const update = "DETMATH_UPDATE"

// A line of the golden table is a function's name, its arguments and its result, each float64 as
// its bits in hexadecimal; the files are written readable by all and writable by their owner.
const (
	nameAndResult  = 2
	hexBase        = 16
	goldenFileMode = 0o644
)

// A function of detmath under test, by the name the golden files give it, with the seed of its own
// stream of inputs.
type function struct {
	name  string
	seed  uint64
	arity int
	call  func(args []float64) float64
}

var functions = []function{
	{name: "exp", seed: 1, arity: 1, call: func(a []float64) float64 { return detmath.Exp(a[0]) }},
	{name: "log", seed: 2, arity: 1, call: func(a []float64) float64 { return detmath.Log(a[0]) }},
	{name: "sigmoid", seed: 3, arity: 1, call: func(a []float64) float64 { return detmath.Sigmoid(a[0]) }},
	{name: "pow", seed: 4, arity: 2, call: func(a []float64) float64 { return detmath.Pow(a[0], a[1]) }},
	{name: "retention", seed: 5, arity: 2, call: func(a []float64) float64 { return detmath.Retention(a[0], a[1]) }},
}

func functionNamed(name string) (function, bool) {
	for _, f := range functions {
		if f.name == name {
			return f, true
		}
	}
	return function{}, false
}

// The points of the edges about exp and σ, each taken with both signs where the list has it.
const (
	// σ is exactly 1 in floats from 53 · ln 2 ≈ 36.737 up, where e**(−x) falls below half an
	// ulp of 1: sigmoidBelowOne is just before it, sigmoidOne just past it.
	sigmoidBelowOne = 36.7
	sigmoidOne      = 36.8
	// large is an argument far from every edge, and insideBounds one near the bounds of Exp, ±708,
	// but within them.
	large        = 100
	insideBounds = 700
	// pastBounds is past the bounds of Exp, which hold it to them, and short of exp's overflow;
	// belowOverflow is just short of the overflow, and pastUnderflow just past the underflow, where
	// exp would give 0.
	pastBounds    = 709
	belowOverflow = 709.78
	pastUnderflow = -745.2
	// farOut is past every threshold of exp, and huge near the largest float.
	farOut = 1000
	huge   = 1e300
)

// The points of the edges about log, each taken with both signs where the list has it: logRoot,
// where log is 0, and ln2Point and halfPoint, where it is ±ln 2; the list has detmath.LogBranch
// too, where log's reduction of its argument changes branch.
const (
	logRoot   = 1
	ln2Point  = 2
	halfPoint = 0.5
)

// Values every function meets at its edges: the ends of the floats, zeros, NaN and the infinities.
var edges = []float64{
	0, math.Copysign(0, -1), math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64,
	detmath.MinNormal, -detmath.MinNormal, math.MaxFloat64, -math.MaxFloat64, math.Inf(1),
	math.Inf(-1), detmath.NaN, detmath.OtherNaN, logRoot, -logRoot, halfPoint, -halfPoint, ln2Point,
	math.E, detmath.LogBranch, detmath.Below(detmath.LogBranch), detmath.ExpNearZero,
	-detmath.ExpNearZero, detmath.Below(detmath.ExpNearZero), sigmoidBelowOne, -sigmoidBelowOne, sigmoidOne, large, -large,
	insideBounds, -insideBounds, detmath.ExpMax, detmath.ExpMin, pastBounds, -pastBounds,
	belowOverflow, pastUnderflow, farOut, -farOut, huge, -huge,
}

// The first words of the seeds of the golden table's draws and of the sweep's; the other word is
// each function's own.
const (
	goldenSeed = 1
	sweepSeed  = 2
)

// Where the inputs are drawn from, beside the edges: past the bounds of Exp on both sides, and
// where σ is neither 0 nor 1; near 1 for log; for the power, the bases of the rule of τ, from 1 up
// to 1 + g = 16, and the shares S̄ within [0, 1]; and pauses and stabilities up to ten years, in
// days, from two minutes.
const (
	expReach      = 720
	sigmoidReach  = 40
	nearOneLow    = 0.5
	nearOneHigh   = 2
	tauBaseHigh   = 1 + growth
	halving       = 0.5 // the base (1/2) of the rule of τ after wrong answers
	twoMinutes    = 1.0 / 720
	powBaseHigh   = 20
	powPowerReach = 2
	tenYears      = 3650
)

// The power's own points: the bases of the rule of τ, 1 + g · (1 − R), at R = 1, at R near 1 and
// at R(τ) = 0.9, below them the bases outside its domain, which Pow holds to its least, and the
// exponents of its exact branches, a**0 = 1 and a**1 = a, half a share, and exponents outside the
// shares, [0, 1], on both sides.
const (
	growth         = 15  // g of the rule of τ
	retentionAtTau = 0.9 // R after a pause of τ
	unitBase       = 1   // 1**s is 1 for every s
	baseNearOne    = 1.5
	baseAtTau      = 1 + growth*(1-retentionAtTau) // 2.5
	zeroBase       = 0
	negativeBase   = -1
	zeroPower      = 0
	unitPower      = 1
	halfShare      = 0.5
	negativePower  = -1
	pastShares     = 2
	shareLow       = 0
	shareHigh      = 1
)

// The retention's own points: pauses at the last practice and before it, of a day and of the
// summer break of the documentation's example of forgetting; and stabilities τ₀, a few days, the
// example's τ before the break and after a right answer on coming back, and 0, which no state
// holds.
const (
	atPractice     = 0
	beforePractice = -1
	aDay           = 1
	summerBreak    = 60
	stabilityFloor = 1
	fewDays        = 2.5
	exampleTau     = 10
	tauAfterBreak  = 63.3
	zeroTau        = 0
)

// How many inputs each draw gives.
const (
	wideDraws      = 100
	narrowDraws    = 40
	tauDraws       = 60
	halvingDraws   = 20
	powDraws       = 20
	retentionDraws = 80
)

// The bits of a float64, 64 in all: its fraction, and the exponents of the normal numbers, from
// the first, 1, up to 2046; 0 is the subnormals' and 2047 the infinities' and NaNs'.
const (
	float64Bits         = 64
	fractionBits        = 52
	firstNormalExponent = 1
	normalExponents     = 2046
)

// goldenInputs gives the inputs of the golden table: the edges, and values drawn where each
// function is used.
func goldenInputs(f function) [][]float64 {
	r := rand.New(rand.NewPCG(goldenSeed, f.seed))
	var in [][]float64
	switch f.name {
	case "exp", "sigmoid":
		for _, x := range edges {
			in = append(in, []float64{x})
		}
		for range wideDraws {
			in = append(in, []float64{uniform(r, -expReach, expReach)})
		}
		for range narrowDraws {
			in = append(in, []float64{uniform(r, -sigmoidReach, sigmoidReach)})
		}
	case "log":
		for _, x := range edges {
			in = append(in, []float64{x})
		}
		for range wideDraws {
			in = append(in, []float64{positiveNormal(r)})
		}
		for range narrowDraws {
			in = append(in, []float64{uniform(r, nearOneLow, nearOneHigh)})
		}
	case "pow":
		// The rule of τ: (1 + g · (1 − R))^S̄ with g = 15 and (1/2)^(1 − S̄), with S̄ within
		// [0, 1]; and the exact branches and the edges.
		bases := []float64{halving, unitBase, baseNearOne, baseAtTau, tauBaseHigh, detmath.MinNormal,
			math.MaxFloat64, zeroBase, negativeBase, detmath.NaN}
		powers := []float64{zeroPower, unitPower, halfShare, negativePower, pastShares, math.Inf(1),
			math.Inf(-1), detmath.NaN}
		for _, a := range bases {
			for _, s := range powers {
				in = append(in, []float64{a, s})
			}
		}
		for range tauDraws {
			in = append(in, []float64{uniform(r, unitBase, tauBaseHigh), uniform(r, shareLow, shareHigh)})
		}
		for range halvingDraws {
			in = append(in, []float64{halving, uniform(r, shareLow, shareHigh)})
		}
		for range powDraws {
			in = append(in, []float64{positiveNormal(r), uniform(r, -powPowerReach, powPowerReach)})
		}
	case "retention":
		pauses := []float64{atPractice, aDay, beforePractice, twoMinutes, summerBreak, tenYears, huge,
			math.Inf(1), detmath.NaN}
		taus := []float64{stabilityFloor, fewDays, exampleTau, tauAfterBreak, zeroTau, detmath.MinNormal,
			math.MaxFloat64, math.Inf(1), detmath.NaN}
		for _, dt := range pauses {
			for _, tau := range taus {
				in = append(in, []float64{dt, tau})
			}
		}
		for range retentionDraws {
			in = append(in, []float64{uniform(r, atPractice, tenYears), uniform(r, stabilityFloor, tenYears)})
		}
	}
	return in
}

// uniform draws from [lo, hi), with every * held by float64(), so that the inputs are the same on
// every build: the rule of float64() holds in tests too.
func uniform(r *rand.Rand, lo, hi float64) float64 {
	u := r.Float64()
	return lo + float64((hi-lo)*u)
}

// positiveNormal draws a positive normal float64, its exponent and fraction uniform.
func positiveNormal(r *rand.Rand) float64 {
	exponent := firstNormalExponent + r.Uint64N(normalExponents)
	return math.Float64frombits(exponent<<fractionBits | r.Uint64()>>(float64Bits-fractionBits))
}

func TestGolden(t *testing.T) {
	path := filepath.Join("testdata", "golden.txt")
	if os.Getenv(update) == "1" {
		writeGolden(t, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	counts := map[string]int{}
	for i, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || strings.HasPrefix(fields[0], "#") {
			continue
		}
		f, ok := functionNamed(fields[0])
		if !ok || len(fields) != f.arity+nameAndResult {
			t.Fatalf("%s:%d: not a line of the table: %q", path, i+1, line)
		}
		bits := make([]uint64, len(fields)-1)
		for j, field := range fields[1:] {
			if bits[j], err = strconv.ParseUint(field, hexBase, 64); err != nil {
				t.Fatalf("%s:%d: %v", path, i+1, err)
			}
		}
		args := make([]float64, f.arity)
		for j := range args {
			args[j] = math.Float64frombits(bits[j])
		}
		want := bits[f.arity]
		if got := math.Float64bits(f.call(args)); got != want {
			t.Errorf("%s:%d: %s%v = %016x, the table has %016x", path, i+1, f.name, args, got, want)
		}
		counts[f.name]++
	}
	for _, f := range functions {
		if counts[f.name] == 0 {
			t.Errorf("%s: no line for %s", path, f.name)
		}
	}
}

func writeGolden(t *testing.T, path string) {
	var b bytes.Buffer
	b.WriteString("# The golden table of detmath, the same on every build (detmath/golden_test.go): a\n")
	b.WriteString("# function, its arguments and its result, each a float64 as its bits in hex. It changes\n")
	b.WriteString("# only with a version of learner-model; DETMATH_UPDATE=1 go test ./detmath writes it.\n")
	for _, f := range functions {
		for _, args := range goldenInputs(f) {
			b.WriteString(f.name)
			for _, a := range args {
				fmt.Fprintf(&b, " %016x", math.Float64bits(a))
			}
			fmt.Fprintf(&b, " %016x\n", math.Float64bits(f.call(args)))
		}
	}
	if err := os.WriteFile(path, b.Bytes(), goldenFileMode); err != nil {
		t.Fatal(err)
	}
}

// The length of the run of each function in the sweep.
const sweepLength = 100_000

// sweepInputs draws the run of inputs for f: half from where f is used, half as raw bits, every
// pattern of a float64 among them — subnormals, infinities, NaNs of every payload.
func sweepInputs(f function, yield func(args []float64)) {
	r := rand.New(rand.NewPCG(sweepSeed, f.seed))
	args := make([]float64, f.arity)
	for i := range sweepLength {
		raw := i%2 == 1
		for j := range args {
			switch {
			case raw:
				args[j] = math.Float64frombits(r.Uint64())
			case f.name == "log":
				args[j] = positiveNormal(r)
			case f.name == "pow" && j == 0:
				args[j] = uniform(r, zeroBase, powBaseHigh)
			case f.name == "pow":
				args[j] = uniform(r, -powPowerReach, powPowerReach)
			case f.name == "retention":
				args[j] = uniform(r, atPractice, tenYears)
			default:
				args[j] = uniform(r, -expReach, expReach)
			}
		}
		yield(args)
	}
}

func TestSweep(t *testing.T) {
	path := filepath.Join("testdata", "sweep.txt")
	got := map[string]string{}
	for _, f := range functions {
		h := sha256.New()
		var buf []byte
		sweepInputs(f, func(args []float64) {
			y := f.call(args)
			if math.IsNaN(y) || math.IsInf(y, 0) || (y != 0 && math.Abs(y) < detmath.MinNormal) {
				t.Errorf("%s%v = %v: not a finite normal number or zero", f.name, args, y)
			}
			buf = binary.LittleEndian.AppendUint64(buf[:0], math.Float64bits(y))
			h.Write(buf)
		})
		got[f.name] = fmt.Sprintf("%x", h.Sum(nil))
	}
	if os.Getenv(update) == "1" {
		var b bytes.Buffer
		b.WriteString("# The SHA-256 of the bits of each function of detmath over its sweep, the same on every\n")
		b.WriteString("# build (detmath/golden_test.go): a function, the length of its run and the hash.\n")
		for _, f := range functions {
			fmt.Fprintf(&b, "%s %d %s\n", f.name, sweepLength, got[f.name])
		}
		if err := os.WriteFile(path, b.Bytes(), goldenFileMode); err != nil {
			t.Fatal(err)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]int{}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || strings.HasPrefix(fields[0], "#") {
			continue
		}
		if _, ok := functionNamed(fields[0]); !ok || len(fields) != 3 || fields[1] != strconv.Itoa(sweepLength) {
			t.Fatalf("%s: not a line of the sweep: %q", path, line)
		}
		if got[fields[0]] != fields[2] {
			t.Errorf("%s: the sweep of %s hashes to %s, the file has %s", path, fields[0], got[fields[0]], fields[2])
		}
		seen[fields[0]]++
	}
	for _, f := range functions {
		if seen[f.name] != 1 {
			t.Errorf("%s: %d lines for %s, not one", path, seen[f.name], f.name)
		}
	}
}
