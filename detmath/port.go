// This file ports exp, expmulti and log from package math of Go 1.27.3 — the pure Go code that
// math itself runs on wasm and 386 — with every float * and / held by float64(), so that no build
// fuses them (docs/architecture.md, section 3.6, rule 1). Beside those, four changes, none of which
// changes a bit: the functions of math are called by their package; int(…) of a float is truncate,
// which checks its range (rule 5); f1 *= 2 is f1 = float64(f1 * 2), since *= cannot be held; and
// NaN() is the same bits by Float64frombits, since math.NaN is not on the list of rule 2. The
// constants, the order of the operations and the special cases are Go's, and on wasm and 386 the
// port gives math.Exp and math.Log bit for bit (R09). Its notices are kept below as both licenses
// ask; they are the notices of package math itself, which every Go binary that imports math
// carries already.
//
// Copyright 2009 The Go Authors.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google LLC nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package detmath

import "math"

// The original C code, the long comment, and the constants
// below are from FreeBSD's /usr/src/lib/msun/src/e_exp.c
// and came with this notice. The go code is a simplified
// version of the original C.
//
// ====================================================
// Copyright (C) 2004 by Sun Microsystems, Inc. All rights reserved.
//
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
//
// exp(x)
// Returns the exponential of x.
//
// Method
//   1. Argument reduction:
//      Reduce x to an r so that |r| <= 0.5*ln2 ~ 0.34658.
//      Given x, find r and integer k such that
//
//               x = k*ln2 + r,  |r| <= 0.5*ln2.
//
//      Here r will be represented as r = hi-lo for better
//      accuracy.
//
//   2. Approximation of exp(r) by a special rational function on
//      the interval [0,0.34658]:
//      Write
//          R(r**2) = r*(exp(r)+1)/(exp(r)-1) = 2 + r*r/6 - r**4/360 + ...
//      We use a special Remez algorithm on [0,0.34658] to generate
//      a polynomial of degree 5 to approximate R. The maximum error
//      of this polynomial approximation is bounded by 2**-59. In
//      other words,
//          R(z) ~ 2.0 + P1*z + P2*z**2 + P3*z**3 + P4*z**4 + P5*z**5
//      (where z=r*r, and the values of P1 to P5 are listed below)
//      and
//          |                  5          |     -59
//          | 2.0+P1*z+...+P5*z   -  R(z) | <= 2
//          |                             |
//      The computation of exp(r) thus becomes
//                             2*r
//              exp(r) = 1 + -------
//                            R - r
//                                 r*R1(r)
//                     = 1 + r + ----------- (for better accuracy)
//                                2 - R1(r)
//      where
//                               2       4             10
//              R1(r) = r - (P1*r  + P2*r  + ... + P5*r   ).
//
//   3. Scale back to obtain exp(x):
//      From step 1, we have
//         exp(x) = 2**k * exp(r)
//
// Special cases:
//      exp(INF) is INF, exp(NaN) is NaN;
//      exp(-INF) is 0, and
//      for finite argument, only exp(0)=1 is exact.
//
// Accuracy:
//      according to an error analysis, the error is always less than
//      1 ulp (unit in the last place).
//
// Misc. info.
//      For IEEE double
//          if x >  7.09782712893383973096e+02 then exp(x) overflow
//          if x < -7.45133219101941108420e+02 then exp(x) underflow
//
// Constants:
// The hexadecimal values are the intended ones for the following
// constants. The decimal values may be used, provided that the
// compiler will convert from decimal to binary accurately enough
// to produce the hexadecimal values shown.

func exp(x float64) float64 {
	const (
		Ln2Hi = 6.93147180369123816490e-01
		Ln2Lo = 1.90821492927058770002e-10
		Log2e = 1.44269504088896338700e+00

		Overflow  = 7.09782712893383973096e+02
		Underflow = -7.45133219101941108420e+02
		NearZero  = 1.0 / (1 << 28) // 2**-28
	)

	// special cases
	switch {
	case math.IsNaN(x):
		return x
	case x > Overflow: // handles case where x is +∞
		return math.Inf(1)
	case x < Underflow: // handles case where x is -∞
		return 0
	case -NearZero < x && x < NearZero:
		return 1 + x
	}

	// reduce; computed as r = hi - lo for extra precision.
	var k int
	switch {
	case x < 0:
		k = truncate(float64(Log2e*x) - 0.5)
	case x > 0:
		k = truncate(float64(Log2e*x) + 0.5)
	}
	hi := x - float64(float64(k)*Ln2Hi)
	lo := float64(float64(k) * Ln2Lo)

	// compute
	return expmulti(hi, lo, k)
}

// exp1 returns e**r × 2**k where r = hi - lo and |r| ≤ ln(2)/2.
func expmulti(hi, lo float64, k int) float64 {
	const (
		P1 = 1.66666666666666657415e-01  /* 0x3FC55555; 0x55555555 */
		P2 = -2.77777777770155933842e-03 /* 0xBF66C16C; 0x16BEBD93 */
		P3 = 6.61375632143793436117e-05  /* 0x3F11566A; 0xAF25DE2C */
		P4 = -1.65339022054652515390e-06 /* 0xBEBBBD41; 0xC5D26BF1 */
		P5 = 4.13813679705723846039e-08  /* 0x3E663769; 0x72BEA4D0 */
	)

	r := hi - lo
	t := float64(r * r)
	c := r - float64(t*(P1+float64(t*(P2+float64(t*(P3+float64(t*(P4+float64(t*P5)))))))))
	y := 1 - ((lo - float64((r*c)/(2-c))) - hi)
	// TODO(rsc): make sure Ldexp can handle boundary k
	return math.Ldexp(y, k)
}

// The original C code, the long comment, and the constants
// below are from FreeBSD's /usr/src/lib/msun/src/e_log.c
// and came with this notice. The go code is a simpler
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
// __ieee754_log(x)
// Return the logarithm of x
//
// Method :
//   1. Argument Reduction: find k and f such that
//			x = 2**k * (1+f),
//	   where  sqrt(2)/2 < 1+f < sqrt(2) .
//
//   2. Approximation of log(1+f).
//	Let s = f/(2+f) ; based on log(1+f) = log(1+s) - log(1-s)
//		 = 2s + 2/3 s**3 + 2/5 s**5 + .....,
//	     	 = 2s + s*R
//      We use a special Reme algorithm on [0,0.1716] to generate
//	a polynomial of degree 14 to approximate R.  The maximum error
//	of this polynomial approximation is bounded by 2**-58.45. In
//	other words,
//		        2      4      6      8      10      12      14
//	    R(z) ~ L1*s +L2*s +L3*s +L4*s +L5*s  +L6*s  +L7*s
//	(the values of L1 to L7 are listed in the program) and
//	    |      2          14          |     -58.45
//	    | L1*s +...+L7*s    -  R(z) | <= 2
//	    |                             |
//	Note that 2s = f - s*f = f - hfsq + s*hfsq, where hfsq = f*f/2.
//	In order to guarantee error in log below 1ulp, we compute log by
//		log(1+f) = f - s*(f - R)		(if f is not too large)
//		log(1+f) = f - (hfsq - s*(hfsq+R)).	(better accuracy)
//
//	3. Finally,  log(x) = k*Ln2 + log(1+f).
//			    = k*Ln2_hi+(f-(hfsq-(s*(hfsq+R)+k*Ln2_lo)))
//	   Here Ln2 is split into two floating point number:
//			Ln2_hi + Ln2_lo,
//	   where n*Ln2_hi is always exact for |n| < 2000.
//
// Special cases:
//	log(x) is NaN with signal if x < 0 (including -INF) ;
//	log(+INF) is +INF; log(0) is -INF with signal;
//	log(NaN) is that NaN with no signal.
//
// Accuracy:
//	according to an error analysis, the error is always less than
//	1 ulp (unit in the last place).
//
// Constants:
// The hexadecimal values are the intended ones for the following
// constants. The decimal values may be used, provided that the
// compiler will convert from decimal to binary accurately enough
// to produce the hexadecimal values shown.

func log(x float64) float64 {
	const (
		Ln2Hi = 6.93147180369123816490e-01 /* 3fe62e42 fee00000 */
		Ln2Lo = 1.90821492927058770002e-10 /* 3dea39ef 35793c76 */
		L1    = 6.666666666666735130e-01   /* 3FE55555 55555593 */
		L2    = 3.999999999940941908e-01   /* 3FD99999 9997FA04 */
		L3    = 2.857142874366239149e-01   /* 3FD24924 94229359 */
		L4    = 2.222219843214978396e-01   /* 3FCC71C5 1D8E78AF */
		L5    = 1.818357216161805012e-01   /* 3FC74664 96CB03DE */
		L6    = 1.531383769920937332e-01   /* 3FC39A09 D078C69F */
		L7    = 1.479819860511658591e-01   /* 3FC2F112 DF3E5244 */
	)

	// special cases
	switch {
	case math.IsNaN(x) || math.IsInf(x, 1):
		return x
	case x < 0:
		return math.Float64frombits(uvnan) // math.NaN()
	case x == 0:
		return math.Inf(-1)
	}

	// reduce
	f1, ki := math.Frexp(x)
	if f1 < math.Sqrt2/2 {
		f1 = float64(f1 * 2)
		ki--
	}
	f := f1 - 1
	k := float64(ki)

	// compute
	s := float64(f / (2 + f))
	s2 := float64(s * s)
	s4 := float64(s2 * s2)
	t1 := float64(s2 * (L1 + float64(s4*(L3+float64(s4*(L5+float64(s4*L7)))))))
	t2 := float64(s4 * (L2 + float64(s4*(L4+float64(s4*L6)))))
	R := t1 + t2
	hfsq := float64(0.5 * f * f)
	return float64(k*Ln2Hi) - ((hfsq - (float64(s*(hfsq+R)) + float64(k*Ln2Lo))) - f)
}

// uvnan is the bits of the NaN that math.NaN returns, in math/bits.go of Go 1.27.3.
const uvnan = 0x7FF8000000000001
