// Package fma is the control of fma-diff in the justfile: the three shapes the compiler fuses, one
// of which at least it fuses on every target the check compiles for, so that the check must see a
// difference here. Nothing else builds it, and no rule of code holds it.
package fma

// MulAdd is a*b + c: a fused multiply-add on every target that fuses.
func MulAdd(a, b, c float64) float64 { return a*b + c }

// QuarterAdd is x/4 + y: a division by a constant whose reciprocal is exact becomes a
// multiplication, and fuses as one.
func QuarterAdd(x, y float64) float64 { return x/4 + y }

// MulSub is c - a*b: a fused multiply-subtract where the target has one; amd64 has none.
func MulSub(a, b, c float64) float64 { return c - a*b }
