package learnermodel

import (
	"math"

	"github.com/altnovatrail/learner-model/detmath"
)

// The level of the trial series (docs/architecture.md, sections 3.1 and 3.3): after every piece of
// evidence of the series, θ is the maximum of the weighted likelihood of all its answers so far,
// with a normal prior of spread σ₀ around the starting point θ₀,
//
//	L(θ) = −(θ − θ₀)² / (2σ₀²) + Σᵢ uᵢ · [Sᵢ · ln Pᵢ + (1 − Sᵢ) · ln(1 − Pᵢ)].
//
// L can have two peaks, so the maximum is looked for on a grid across all of θ₀ ± Σu · σ₀², where
// it must lie, and the best cell is then halved on the sign of L′, as SPEC 2.2.1 of standalone
// does.

// trialPoint is an answer of the series as the likelihood weighs it, worked out once.
type trialPoint struct {
	beta    float64
	guess   float64
	correct bool
	weight  float64
}

func trialPoints(answers []TrialAnswer) []trialPoint {
	points := make([]trialPoint, 0, len(answers))
	for _, a := range answers {
		points = append(points, trialPoint{
			beta:    beta(a.BandPlace, a.Difficulty),
			guess:   guess(a.Slots),
			correct: a.Verdict == VerdictCorrect,
			weight:  a.Weight,
		})
	}
	return points
}

// trialLevel is the θ of the series: the maximum of L, or start with no answers.
func trialLevel(start float64, answers []TrialAnswer) float64 {
	points := trialPoints(answers)
	if len(points) == 0 {
		return start
	}
	// No answer pulls the maximum with a slope steeper than its weight, and the prior pulls it
	// back by (θ − θ₀)/σ₀², so it lies within Σu · σ₀² of the start.
	var total float64
	for _, p := range points {
		total += p.weight
	}
	reach := float64(total * startSpread2)
	cells := gridCells(math.Ceil(float64(2 * reach / gridStep)))

	best, highest := start, math.Inf(-1)
	for cell := 0; cell <= cells; cell++ {
		level := start - reach + float64(float64(cell)*gridStep)
		if l := logLikelihood(level, start, points); l > highest {
			best, highest = level, l
		}
	}

	// The maximum is within a cell of the best point, where L′ turns from rising to falling.
	low, high := best-gridStep, best+gridStep
	for range halvings {
		middle := low + float64((high-low)/2)
		if middle <= low || middle >= high {
			break
		}
		if logLikelihoodSlope(middle, start, points) > 0 {
			low = middle
		} else {
			high = middle
		}
	}
	return low + float64((high-low)/2)
}

// logLikelihood is L at a level. The term of an answer takes its branch by the outcome, not by a
// factor (1 − S), and ln(1 − P) is ln(1 − c) + ln σ(−x), formed without 1 − P: far above a task σ
// is exactly 1, and 0 · ln 0 would be NaN (docs/architecture.md, section 3.6, rule 4).
func logLikelihood(level, start float64, points []trialPoint) float64 {
	d := level - start
	sum := -float64(d * d / (2 * startSpread2))
	for _, p := range points {
		x := level - p.beta
		var term float64
		if p.correct {
			term = detmath.Log(probability(p.guess, x))
		} else {
			term = detmath.Log(1-p.guess) + detmath.Log(detmath.Sigmoid(-x))
		}
		sum += float64(p.weight * term)
	}
	return sum
}

// logLikelihoodSlope is L′ at a level: the prior pulling back, a right answer pulling up by
// u · (1 − c) · σ(1 − σ) / P and a wrong one down by u · σ.
func logLikelihoodSlope(level, start float64, points []trialPoint) float64 {
	pull := -float64((level - start) / startSpread2)
	for _, p := range points {
		s := detmath.Sigmoid(level - p.beta)
		if p.correct {
			rise := float64((1 - p.guess) * s * (1 - s))
			pull += float64(p.weight * rise / probability(p.guess, level-p.beta))
		} else {
			pull -= float64(p.weight * s)
		}
	}
	return pull
}
