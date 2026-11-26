package learnermodel

import "math"

// The rules in place of a calibration of the model's constants: the start a band below the band of
// the learner's grade — startLevel in model.go —, an adult's correction of a module's level, and a
// guard that lowers the level after three failures of the module in a row.

// correct is an adult's correction of the module's level to the band at a place of its current
// version, the step of evidence that LevelCorrection gives: the level goes to the band's middle,
// where it stands in that band, unless it stands in it already. A place above the version's
// highest band is that band: a correction logged on another version must fold. It ends the
// trial series, if it lasts, and starts the count of the guard anew; when the level stands in the
// band already, past the series, with the count at 0, nothing changes and correct reports so. It
// moves no skill, and its time is no practice.
func (ms *ModuleState) correct(place int32, g *Graph) bool {
	if len(g.Bands) > 0 && int(place) > len(g.Bands) {
		place = g.Bands[len(g.Bands)-1].Place
	}
	trial := ms.inTrial()
	if b, ok := levelBand(ms.Level, g.Bands); !ok || b.Place != place {
		ms.Level = shift(place)
	} else if !trial && ms.FailuresInARow == 0 {
		return false
	}
	if trial {
		ms.TrialEndedByAdult = true
	}
	ms.Trial = nil
	ms.LastInTrial = false
	ms.FailuresInARow = 0
	return true
}

// stuck tells whether a piece of evidence with answers is a failure of the module for the guard:
// the model expected the learner to know at least half of its answers — the sum of their P above
// the guessing floor, R̄ · σ, from the state before it, is at least half their number, so that a
// choice of two, whose P is never below ½, is not fitting by its guessing alone —, and fewer than
// half of the items with an answer were solved at the end of the attempt. An item is solved when
// its final verdict is right, or, with no final verdict recorded, its first one; an item the step
// reads as "not solved", the solution shown, is not. So a wrong first answer the learner put right
// on a rung of the ladder is no failure, nor is a task the model already took for too hard: the
// guard is for tasks the model gave as fitting and the learner was stuck on, which tell that the
// level stands too high.
func stuck(items []item, answers []answer) bool {
	var expected float64
	for _, a := range answers {
		expected += a.known
	}
	var answered, solved int
	for _, it := range items {
		if !it.answered {
			continue
		}
		answered++
		a := it.answer
		if (a.HasFinalVerdict && a.FinalVerdict == VerdictCorrect) ||
			(!a.HasFinalVerdict && a.HasFirstVerdict && a.FirstVerdict == VerdictCorrect) {
			solved++
		}
	}
	return fewerThanHalf(solved, answered) && atLeastHalf(expected, answered)
}

// fewerThanHalf tells whether part is fewer than half of whole.
func fewerThanHalf(part, whole int) bool { return 2*part < whole }

// atLeastHalf tells whether x is at least half of whole: x is doubled, which is exact, so that no
// rounding of a half decides it.
func atLeastHalf(x float64, whole int) bool { return float64(2*x) >= float64(whole) }

// guard counts a piece of evidence past the trial series: a failure adds one to the count, and the
// third in a row lowers the level by a band, but not below levelFloor nor upward from below it;
// evidence that is no failure starts the count anew. A run of failures lowers the level once, at
// its third: the next band down waits for a piece of evidence that is no failure and three more.
func (ms *ModuleState) guard(failed bool) {
	if !failed {
		ms.FailuresInARow = 0
		return
	}
	if ms.FailuresInARow < math.MaxInt32 {
		ms.FailuresInARow++
	}
	if ms.FailuresInARow == guardFailures && ms.Level > levelFloor {
		ms.Level = max(ms.Level-bandStep, levelFloor)
	}
}
