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
	return c + float64((1-c)*detmath.Sigmoid(x))
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
