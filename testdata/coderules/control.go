// Package coderules breaks every rule of code that rules_test.go holds learner-model to, each
// line marked "// want" with the rules it breaks, and keeps beside each break what the rules allow,
// so that the check is held to finding no more than it should either.
package coderules

import "C" // want cgo

import (
	"math"
	"os" // want import
	"strings"
	_ "unsafe" // want import

	"example.com/coderules/sub"
)

var _ = os.Args

//go:linkname now runtime.nanotime // want linkname
func now() int64 { return 0 }

func fused(a, b, c, x, y float64) float64 {
	p := a * b   // want fused
	q := x/4 + y // want fused
	r := c - a*b // want fused
	s := float64(a*b) + c
	t := float64(a*b*c) - float64(x/y)
	u := (a * b) + c // want fused
	if a*b < c {
		return float64(-(a * b)) // want fused
	}
	const third = 1.0 / 3
	return identity(a*b) + p + q + r + s + t + u + third // want fused
}

func identity(x float64) float64 { return x }

func integers(i, j, k int) int {
	n := i*j + k
	n *= 2
	n /= 3
	return n
}

func assignments(x float64) float64 {
	x *= 2 // want assignop
	x /= 3 // want assignop
	x += 1
	return x
}

func complexes(a, b float64) float64 {
	z := complex(a, b) // want complex
	_ = z              // want complex
	return a
}

func floats[T ~float32 | ~float64](a, b T) T {
	return a + b // want typeparam
}

type number interface{ ~int | ~float64 }

func numbers[T number](a, b T) T {
	a -= b       // want typeparam
	return a * b // want typeparam
}

func ints[T ~int | ~int64](a, b T) T {
	return a*b + a
}

func calls(x float64) float64 {
	return math.Sqrt(x) + math.Exp(x) + math.Floor(x) + math.Pi // want math
}

func ranges(m map[string]int, s []int) int {
	n := 0
	for k := range m { // want maprange
		n += len(k)
	}
	for _, v := range s {
		n += v
	}
	return n
}

func ids(a, b string) bool {
	return strings.ToLower(a) == b || strings.HasPrefix(a, b) // want unicode
}

func goroutines() {
	go func() {}() // want goroutine
}

// Exported is what external_test.go imports, from a package of the module itself.
func Exported(x float64) float64 { return x }

func genericKeys[M ~map[string]int](m M) int {
	n := 0
	for k := range m { // want maprange
		n += len(k)
	}
	return n
}

func increments[T ~float64](a T) T {
	a++ // want typeparam
	return a
}

// Make gives a value of sub, which this package imports and whose external tests import this one
// back: those tests see sub with its tests in it, and this package must be checked against that
// sub, not against the one it was first checked with.
func Make() sub.Value { return sub.Value{X: 1} }
