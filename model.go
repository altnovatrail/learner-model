package learnermodel

import (
	"math"

	"github.com/altnovatrail/learner-model/detmath"
)

// The constants of this version of learner-model (docs/architecture.md, sections 3.1 and 3.3).
// Changing one is a new version, and only through sim's bench.
const (
	// bandStep is how far a band of a module's scale stands above the one before it.
	bandStep = 2.5
	// startSpread is σ₀, the spread of the trial series around the starting point, and
	// startSpread2 its square.
	startSpread  = 2.5
	startSpread2 = startSpread * startSpread

	// The step of Elo: K₀ of the level θ and of a skill's correction δ, how fast both decay with
	// the weight of answers behind them, and the floor of θ's.
	levelStep0      = 0.2
	correctionStep0 = 0.4
	stepDecay       = 0.05
	levelStepFloor  = 0.05

	// information is I, what one answer at the middle of the corridor tells of a level, and
	// correctionPrior the weight of answers a new correction is as uncertain as, both in v.
	information     = 0.149981
	correctionPrior = 20.0

	// masteryThreshold is σ of the cautious level against the middle task of a band that masters
	// a skill at that band: P = 0.775 on a task of five options.
	masteryThreshold = 0.71875

	// trialWeight is the weight of answers the trial series lasts, and masteryWeight the weight
	// of answers a skill needs before it is mastered; weightTolerance is ε, with which a sum of
	// weights is compared with either (docs/architecture.md, section 3.5).
	trialWeight     = 5.0
	masteryWeight   = 5.0
	weightTolerance = 1e-9

	// failuresToLose is how many failed pieces of evidence in a row take a mastery away.
	failuresToLose = 2

	// gridStep is the step of the grid the trial series' maximum is first looked for on, and
	// halvings the most halvings of the best cell after it. The halving stops sooner where the
	// cell has narrowed to neighbouring floats; near 0, where floats are dense, it stops at 128,
	// when the cell of 0.1 is below 3e-40 wide, far past what a level needs. Either way the result
	// is the same on every build.
	gridStep = 0.05
	halvings = 128
)

// The constants of forgetting (docs/architecture.md, section 3.2): the retention is
// R(Δt, τ) = (1 + F · Δt / τ)^(−C) with C = 1/2, so that R^(−1/C) is 1/R², and F =
// detmath.RetentionFactor = 19/81, so that R(τ) = 0.9.
const (
	// stabilityFloor is τ₀, in days: the stability of a skill after its first answer, the least it
	// falls to, and the least pause before a review date.
	stabilityFloor = 1.0
	// stabilityGrowth is g: a right answer at retention R multiplies τ by 1 + g · (1 − R).
	stabilityGrowth = 15.0
	// secondsPerDay turns a pause in whole seconds into days; time knows no calendar.
	secondsPerDay = 86400
	// maxPauseSeconds is the longest pause, 2⁵³ s, past which a difference of seconds would not
	// convert to a float exactly (docs/architecture.md, section 3.6, rule 5): some 285 million
	// years.
	maxPauseSeconds = 1 << 53
)

// pauseDays is Δt = max(0, t − last) in days, from whole seconds: the difference is taken in
// integers, which no overflow can wrap, and is kept below 2⁵³ before it becomes a float.
func pauseDays(t, last int64) float64 {
	if t <= last {
		return 0
	}
	d := uint64(t) - uint64(last) // t > last, so the true difference fits in a uint64
	if d > maxPauseSeconds {
		d = maxPauseSeconds
	}
	return float64(float64(d) / secondsPerDay)
}

// retentionAt is R of a skill's state at time t, in whole seconds: 1 at its last practice and
// before it, and falling after it with its stability.
func retentionAt(sk SkillState, t int64) float64 {
	return detmath.Retention(pauseDays(t, sk.LastPractice), sk.Stability)
}

// chance is P = c + (1 − c) · R̄ · σ of a right answer, given σ = σ(θ + δ̄ − β) and the mean
// retention R̄ of the item's skills. At R̄ = 1 it is the probability of section 3.1 to the bit,
// since (1 − c) · 1 is 1 − c exactly.
func chance(c, sigma, rbar float64) float64 {
	return c + float64((1-c)*rbar*sigma)
}

// retentionWeight is w = R̄ · P₁ · (1 − P₁) / (P · (1 − P)), how much less an answer tells of the
// skill than it would without forgetting, P₁ being P at R̄ = 1 (docs/architecture.md, section
// 3.2). Both 1 − P₁ and 1 − P carry the factor 1 − c, which is cancelled, so that the weight is
// R̄ · P₁ · (1 − σ) / (P · (1 − R̄σ)); with no guessing, P₁ / P is 1/R̄ and it is
// (1 − σ) / (1 − R̄σ), with no small number divided by another. At R̄ = 1 it is 1 exactly. p1 and p
// are P₁ and P, as chance gives them, so that the weight and the step read one P.
func retentionWeight(c, sigma, rbar, p1, p float64) float64 {
	switch {
	case rbar == 1:
		return 1
	case c == 0:
		return float64((1 - sigma) / (1 - float64(rbar*sigma)))
	}
	return float64(rbar * p1 * (1 - sigma) / (p * (1 - float64(rbar*sigma))))
}

// stabilityAfter is τ after a piece of evidence in which share of the skill's answers are right,
// at retention r: max(τ₀, τ · (1 + g · (1 − r))^S̄ · (1/2)^(1 − S̄)) (docs/architecture.md, section
// 3.3). At S̄ = 1 and 0 the powers are 1 and their bases exactly, so it is the rule of section 3.2:
// τ · (1 + g · (1 − r)) after right answers and max(τ₀, τ / 2) after wrong ones. It is kept below
// the largest float, which no real pause comes near, so that no infinity comes out of a state
// whatever it holds.
func stabilityAfter(tau, r, share float64) float64 {
	grow := detmath.Pow(1+float64(stabilityGrowth*(1-r)), share)
	halve := detmath.Pow(0.5, 1-share)
	return max(stabilityFloor, min(float64(tau*grow*halve), math.MaxFloat64))
}

// masteryChance is q = σ(λ − shift), the cautious chance of the middle task of the band a skill is
// mastered at, without forgetting (docs/architecture.md, section 3.2).
func masteryChance(level, n float64, sk SkillState) float64 {
	return detmath.Sigmoid(cautiousLevel(level, sk.Correction, n, sk.AnswerWeight) - shift(sk.MasteryBand))
}

// reviewDate is when a mastered skill is due for review, from q, its stability and its last
// practice: Already when q < 0.71875 — the threshold does not hold even without forgetting —
// and otherwise the last practice plus the pause Δt* after which R · q falls to 0.71875, but not
// less than τ₀: R* = 0.71875 / q and Δt* = τ · (1/R*² − 1) / F. The pause becomes whole seconds
// rounded up, and a date past the largest int64 saturates. Δt* comes of three roundings, so the
// second it gives is then held to the retention itself: while R · q there still holds the
// threshold, the date moves on a second, so that "due" never comes before the threshold.
func reviewDate(q, tau float64, last int64) ReviewDate {
	if q < masteryThreshold {
		return ReviewDate{Already: true}
	}
	rStar := float64(masteryThreshold / q)
	inverse := float64(1 / float64(rStar*rStar))
	pause := float64(tau * (inverse - 1) / detmath.RetentionFactor)
	if pause <= stabilityFloor {
		// τ₀ is the date: R · q at it is below the threshold already.
		return ReviewDate{At: addSeconds(last, stabilityFloor*secondsPerDay)}
	}
	seconds := math.Ceil(float64(pause * secondsPerDay))
	for range 4 {
		if !(seconds < maxPauseSeconds) {
			break
		}
		r := detmath.Retention(float64(seconds/secondsPerDay), tau)
		if float64(r*q) <= masteryThreshold {
			break
		}
		seconds++
	}
	return ReviewDate{At: addSeconds(last, seconds)}
}

// addSeconds adds a whole number of seconds, a non-negative float, to a time, converting it only
// within the range where every build converts alike (docs/architecture.md, section 3.6, rule 5)
// and saturating at the largest int64.
func addSeconds(t int64, seconds float64) int64 {
	const limit = 1 << 62
	if !(seconds < limit) {
		return math.MaxInt64
	}
	s := int64(max(seconds, 0))
	if t > math.MaxInt64-s {
		return math.MaxInt64
	}
	return t + s
}

// shift is where the band at a place, from 1, stands on the module's scale: 2.5 · (place − 1).
func shift(place int32) float64 {
	return float64(bandStep * float64(place-1))
}

// beta is where a task of the band at a place, of a difficulty from 1 to 5, stands on the scale.
func beta(place, difficulty int32) float64 {
	return shift(place) + float64(difficulty-3)
}

// guess is c, the chance of guessing an item: the product over its slots of 1/k for a choice of
// k options, and 0 as soon as one slot is checked otherwise. An item is right only when every
// slot is, so guessing it is guessing them all.
func guess(slots []Slot) float64 {
	c := 1.0
	for _, s := range slots {
		if s.Check != CheckKindChoice {
			return 0
		}
		c = float64(c / float64(s.Options))
	}
	return c
}

// probability is P = c + (1 − c) · σ(x) of a right answer, where x = θ + δ̄ − β.
func probability(c, x float64) float64 {
	return chance(c, detmath.Sigmoid(x), 1)
}

// levelStep is K_θ after answers of weight n: 0.2 / (1 + 0.05 · n), but not below 0.05.
func levelStep(n float64) float64 {
	return max(float64(levelStep0/(1+float64(stepDecay*n))), levelStepFloor)
}

// correctionStep is K_δ after answers of weight m: 0.4 / (1 + 0.05 · m).
func correctionStep(m float64) float64 {
	return float64(correctionStep0 / (1 + float64(stepDecay*m)))
}

// seriesWeight is u = 3 / (N + 2), the weight of each of N answers of one attempt: they are one
// cluster with an intraclass correlation of 1/3, so the attempt weighs at most three answers.
func seriesWeight(n int) float64 {
	return float64(3 / float64(n+2))
}

// uncertainty is v = 1 / (1/σ₀² + I · n) + 1 / (I · (20 + m)), the uncertainty of θ + δ after
// answers of weight n in the module and m in the skill.
func uncertainty(n, m float64) float64 {
	return float64(1/(1/startSpread2+float64(information*n))) + float64(1/(information*(correctionPrior+m)))
}

// cautiousLevel is λ = θ + δ − √v.
func cautiousLevel(level, correction, n, m float64) float64 {
	return level + correction - math.Sqrt(uncertainty(n, m))
}

// masteryBand is the band a skill is mastered at by a cautious level λ, after a task of the band
// at taskPlace: the highest place the skill is taught at, not above the task's, whose middle task
// λ clears, σ(λ − shift) ≥ 0.71875. taught is ascending; false when no band qualifies.
func masteryBand(taught []int32, taskPlace int32, lambda float64) (int32, bool) {
	for i := len(taught) - 1; i >= 0; i-- {
		p := taught[i]
		if p <= taskPlace && detmath.Sigmoid(lambda-shift(p)) >= masteryThreshold {
			return p, true
		}
	}
	return 0, false
}

// gridCells converts the number of cells of the trial series' grid to an int, checking its range
// (docs/architecture.md, section 3.6, rule 5): the cells of a sum of weights below 8 are at most
// 2000, and a larger x saturates.
func gridCells(x float64) int {
	const limit = 1 << 20
	if x >= 0 && x <= limit {
		return int(x)
	}
	if x > limit {
		return limit
	}
	return 0
}
