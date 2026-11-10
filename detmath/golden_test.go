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

// The golden files of detmath, one for every build (docs/architecture.md, section 3.6): they tell
// a difference in this arithmetic from a difference in the model. testdata/golden.txt holds
// inputs and the bits of the results; testdata/sweep.txt holds, for each function, the SHA-256 of
// the bits of its results over a long run of inputs drawn by PCG with fixed seeds, the same on
// every build. Both change only with a version of learner-model: DETMATH_UPDATE=1 writes them
// anew, from the build it runs on.

const update = "DETMATH_UPDATE"

// A function of detmath under test, by the name the golden files give it, with the seed of its own
// stream of inputs.
type function struct {
	name  string
	seed  uint64
	arity int
	call  func(args []float64) float64
}

var functions = []function{
	{"exp", 1, 1, func(a []float64) float64 { return detmath.Exp(a[0]) }},
	{"log", 2, 1, func(a []float64) float64 { return detmath.Log(a[0]) }},
	{"sigmoid", 3, 1, func(a []float64) float64 { return detmath.Sigmoid(a[0]) }},
	{"pow", 4, 2, func(a []float64) float64 { return detmath.Pow(a[0], a[1]) }},
	{"retention", 5, 2, func(a []float64) float64 { return detmath.Retention(a[0], a[1]) }},
}

func functionNamed(name string) (function, bool) {
	for _, f := range functions {
		if f.name == name {
			return f, true
		}
	}
	return function{}, false
}

// Values every function meets at its edges: the ends of the floats, zeros, NaN and the infinities.
var edges = []float64{
	0, math.Copysign(0, -1), math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64, 0x1p-1022,
	-0x1p-1022, math.MaxFloat64, -math.MaxFloat64, math.Inf(1), math.Inf(-1),
	math.Float64frombits(0x7FF8000000000001), math.Float64frombits(0xFFF0000000000001),
	1, -1, 0.5, -0.5, 2, math.E, math.Sqrt2 / 2, below(math.Sqrt2 / 2), 0x1p-28,
	-0x1p-28, below(0x1p-28), 36.7, -36.7, 36.8, 100, -100, 700, -700, 708, -708,
	709, -709, 709.78, -745.2, 1000, -1000, 1e300, -1e300,
}

// goldenInputs gives the inputs of the golden table: the edges, and values drawn where each
// function is used.
func goldenInputs(f function) [][]float64 {
	r := rand.New(rand.NewPCG(1, f.seed))
	var in [][]float64
	switch f.name {
	case "exp", "sigmoid":
		for _, x := range edges {
			in = append(in, []float64{x})
		}
		for range 100 {
			in = append(in, []float64{uniform(r, -720, 720)})
		}
		for range 40 {
			in = append(in, []float64{uniform(r, -40, 40)})
		}
	case "log":
		for _, x := range edges {
			in = append(in, []float64{x})
		}
		for range 100 {
			in = append(in, []float64{positiveNormal(r)})
		}
		for range 40 {
			in = append(in, []float64{uniform(r, 0.5, 2)})
		}
	case "pow":
		// The rule of τ of docs/architecture.md, section 3.3: (1 + g · (1 − R))^S̄ with g up to 15, a base
		// within [1, 16], and b^(1 − S̄) with b = 0.35, with S̄ within [0, 1]; and the exact
		// branches and the edges.
		for _, a := range []float64{0.35, 0.5, 1, 1.5, 2.5, 16, 0x1p-1022, math.MaxFloat64, 0, -1, nan} {
			for _, s := range []float64{0, 1, 0.5, -1, 2, math.Inf(1), math.Inf(-1), nan} {
				in = append(in, []float64{a, s})
			}
		}
		for range 60 {
			in = append(in, []float64{uniform(r, 1, 16), uniform(r, 0, 1)})
		}
		for range 20 {
			in = append(in, []float64{0.5, uniform(r, 0, 1)})
		}
		for range 20 {
			in = append(in, []float64{positiveNormal(r), uniform(r, -2, 2)})
		}
	case "retention":
		for _, dt := range []float64{0, 1, -1, 1.0 / 720, 60, 3650, 1e300, math.Inf(1), nan} {
			for _, tau := range []float64{1, 2.5, 10, 63.3, 0, 0x1p-1022, math.MaxFloat64, math.Inf(1), nan} {
				in = append(in, []float64{dt, tau})
			}
		}
		for range 80 {
			in = append(in, []float64{uniform(r, 0, 3650), uniform(r, 1, 3650)})
		}
	}
	return in
}

// uniform draws from [lo, hi), with every * held by float64(), so that the inputs are the same on
// every build (rule 1 holds in tests too: R09, the amendment of 2026-11-07).
func uniform(r *rand.Rand, lo, hi float64) float64 {
	u := r.Float64()
	return lo + float64((hi-lo)*u)
}

// positiveNormal draws a positive normal float64, its exponent and fraction uniform.
func positiveNormal(r *rand.Rand) float64 {
	exponent := 1 + r.Uint64N(2046)
	return math.Float64frombits(exponent<<52 | r.Uint64()>>12)
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
		if !ok || len(fields) != f.arity+2 {
			t.Fatalf("%s:%d: not a line of the table: %q", path, i+1, line)
		}
		bits := make([]uint64, len(fields)-1)
		for j, field := range fields[1:] {
			if bits[j], err = strconv.ParseUint(field, 16, 64); err != nil {
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
	if err := os.WriteFile(path, b.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

// The length of the run of each function in the sweep.
const sweepLength = 100_000

// sweepInputs draws the run of inputs for f: half from where f is used, half as raw bits, every
// pattern of a float64 among them — subnormals, infinities, NaNs of every payload.
func sweepInputs(f function, yield func(args []float64)) {
	r := rand.New(rand.NewPCG(2, f.seed))
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
				args[j] = uniform(r, 0, 20)
			case f.name == "pow":
				args[j] = uniform(r, -2, 2)
			case f.name == "retention":
				args[j] = uniform(r, 0, 3650)
			default:
				args[j] = uniform(r, -720, 720)
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
			if math.IsNaN(y) || math.IsInf(y, 0) || (y != 0 && math.Abs(y) < 0x1p-1022) {
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
		if err := os.WriteFile(path, b.Bytes(), 0o644); err != nil {
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

// The NaN that math.NaN gives, and the float just below a positive one: rule 2 keeps tests to the
// functions of math with an exact result, and these need none.
var nan = math.Float64frombits(0x7FF8000000000001)

func below(x float64) float64 { return math.Float64frombits(math.Float64bits(x) - 1) }
