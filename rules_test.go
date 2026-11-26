package learnermodel_test

import (
	"errors"
	"fmt"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

// The rules in place of a calibration: the start a band below the band of the
// grade, the band the level stands in, an adult's correction and the guard.

// The numbers the tests of the rules share.
const (
	// rulesAt is the time of a settled learner's last practice, and of the answers and the readings
	// that follow it: no pause, so nothing is forgotten.
	rulesAt = 1000
	// settledWeight is n of a settled learner, past the trial series.
	settledWeight = 20
	// guardOneShort is a count of the guard one failure short of lowering the level.
	guardOneShort = lm.GuardFailures - 1
	// middleOfBand2 and middleOfBand3 are where the bands 3-4 and 5-6 of the ladder stand on the
	// scale, where a correction to them sets the level.
	middleOfBand2 = lm.BandStep
	middleOfBand3 = 2 * lm.BandStep
	// kindBeforeCorrection is the last kind of evidence that does not record LevelSetByAdult.
	kindBeforeCorrection = 2
)

// settled is a learner of olympiad maths past the trial series, at a level, with one skill a.
func settled(level float64) lm.State {
	return lm.State{Kind: lm.StateKind, Modules: []lm.ModuleState{{
		Module: olympiad, Level: level, AnswerWeight: settledWeight, StartBand: ladderBand(lm.FirstPlace),
		Skills: []lm.SkillState{{Skill: topic("a"), AnswerWeight: 8, EvidenceCount: 8, LastPractice: rulesAt, Stability: 1}},
	}}}
}

func moduleOf(t *testing.T, r lm.Readout, id string) lm.ModuleReadout {
	t.Helper()
	for _, m := range r.Modules {
		if m.Module == id {
			return m
		}
	}
	t.Fatalf("no reading of module %s", id)
	return lm.ModuleReadout{}
}

// levelBandOf reads the band the level of olympiad maths stands in, or "none".
func levelBandOf(t *testing.T, s lm.State, g lm.Graphs, grade lm.Grade) string {
	t.Helper()
	m := moduleOf(t, mustRead(t, s, g, rulesAt, grade, olympiad), olympiad)
	if !m.HasLevelBand {
		if m.LevelBand != (lm.BandRef{}) {
			t.Errorf("no level band, but %+v", m.LevelBand)
		}
		return "none"
	}
	return band(m.LevelBand)
}

// The start: a learner of grade 3, of the band 3-4, starts at its offset less a band, at 0, in the
// band 1-2; one of grade 1 or of no grade starts below the band 1-2, at −2.5, and stands in no band.
func TestStartBelowTheGrade(t *testing.T) {
	g := olympiadGraphs(t, map[string][]int32{topic("a"): {1, 2, 3}})
	for _, c := range []struct {
		grade      lm.Grade
		start, lvl string
		p          float64
	}{
		{lm.Grade{Value: 3, Known: true}, "2/3-4", "1/1-2", 0.6},
		{lm.Grade{Value: 5, Known: true}, "3/5-6", "2/3-4", 0.2 + 0.8*0.9241418},
		{lm.Grade{Value: 1, Known: true}, "1/1-2", "none", 0.2 + 0.8*0.0758582},
		{lm.Grade{}, "1/1-2", "none", 0.2 + 0.8*0.0758582},
	} {
		r := mustRead(t, lm.State{}, g, 0, c.grade, olympiad)
		m := moduleOf(t, r, olympiad)
		if band(m.StartBand) != c.start || levelBandOf(t, lm.State{}, g, c.grade) != c.lvl || m.HasEvidence {
			t.Errorf("grade %+v: start %s, level band %s; want %s and %s", c.grade, band(m.StartBand), levelBandOf(t, lm.State{}, g, c.grade), c.start, c.lvl)
		}
		// The candidate of the band 1-2 of difficulty 3, β = 0, at the starting point.
		task := candidate(topic("a"))
		task.Band = ladderBand(1)
		near3(t, "P at the start of grade "+itoa(int(c.grade.Value)), taskP(t, r, task)[0].Probability, c.p)
	}
	// The first evidence starts the module at the same point: an answer weighing nothing in the
	// series would leave θ at θ₀ = −2.5, a band's width below the lowest band, and a right one lifts
	// it from there.
	const theta0 = -lm.StartBelow
	s := mustAppend(t, lm.State{}, choiceAnswer(choiceSpec{skill: topic("a"), place: 1, difficulty: 3, correct: true, start: ladderBand(1), at: rulesAt}), g)
	if m := s.Modules[0]; m.StartBand != ladderBand(lm.FirstPlace) || !(m.Level > theta0 && m.Level < 0) {
		t.Errorf("after a right answer from the start of grade 1: start %+v, θ %v", m.StartBand, m.Level)
	}
}

// The band the level stands in: the band whose middle is nearest, the upper of two as near, the
// highest above it, none more than half a band below the lowest, and none without bands to name.
func TestLevelBand(t *testing.T) {
	g := olympiadGraphs(t, map[string][]int32{topic("a"): {1, 2, 3}})
	const (
		justBelow = 1e-7 // a level this much below a boundary between bands is on its lower side
		farAbove  = 100  // a level far above the highest band
	)
	for level, want := range map[float64]string{
		-lm.HalfBand - justBelow: "none", -lm.HalfBand: "1/1-2", 0: "1/1-2",
		lm.HalfBand - justBelow: "1/1-2", lm.HalfBand: "2/3-4",
		middleOfBand2 + lm.HalfBand: "3/5-6", middleOfBand3: "3/5-6", farAbove: "3/5-6",
	} {
		if got := levelBandOf(t, settled(level), g, lm.Grade{}); got != want {
			t.Errorf("θ = %v stands in band %s, want %s", level, got, want)
		}
	}
	// A module with no current version names no band, nor does one with no bands.
	waiting := lm.Graphs{olympiad: {Current: false}}
	if m := mustRead(t, settled(0), waiting, rulesAt, lm.Grade{}).Modules[0]; m.HasLevelBand || m.LevelBand != (lm.BandRef{}) {
		t.Errorf("a module with no current version: %+v", m)
	}
	bare, err := lm.NewGraph(olympiad, nil, []lm.Skill{{ID: topic("a"), Assessability: lm.AssessabilityVerdict}}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	noBands := lm.Graphs{olympiad: {Graph: bare, Current: true}}
	for _, s := range []lm.State{settled(0), {}} {
		if m := moduleOf(t, mustRead(t, s, noBands, rulesAt, lm.Grade{}, olympiad), olympiad); m.HasLevelBand {
			t.Errorf("a module with no bands: %+v", m)
		}
	}
}

// An adult's correction sets the level at the middle of the band set, in one step, and a repeat
// changes nothing; it moves no skill, ends the trial series and starts the guard's count anew.
func TestCorrection(t *testing.T) {
	g := olympiadGraphs(t, map[string][]int32{topic("a"): {1, 2, 3}})
	harder := lm.LevelCorrection(olympiad, ladderBand(2), ladderBand(3), 5000)
	const settledLevel = 0.3
	before := settled(settledLevel)
	before.Modules[0].FailuresInARow = guardOneShort
	s, e, err := lm.Append(before, harder, g)
	if err != nil {
		t.Fatal(err)
	}
	m := s.Modules[0]
	if m.Level != middleOfBand2 || m.AnswerWeight != settledWeight || m.FailuresInARow != 0 || m.LastInTrial || m.TrialEndedByAdult ||
		fmt.Sprint(m.Skills) != fmt.Sprint(before.Modules[0].Skills) || e.StartBand != ladderBand(1) {
		t.Errorf("after a correction to the band 3-4: %+v, the evidence's start %+v", m, e.StartBand)
	}
	if got := levelBandOf(t, s, g, lm.Grade{}); got != "2/3-4" {
		t.Errorf("after the correction the level stands in %s", got)
	}
	again := mustAppend(t, s, lm.LevelCorrection(olympiad, ladderBand(2), ladderBand(1), 6000), g)
	if !sameState(t, again, s) {
		t.Errorf("a repeat of the band changes the state: %+v", again.Modules[0])
	}
	easier := mustAppend(t, s, lm.LevelCorrection(olympiad, ladderBand(1), ladderBand(1), 6000), g)
	if easier.Modules[0].Level != 0 {
		t.Errorf("after a correction to the band 1-2: θ %v", easier.Modules[0].Level)
	}
	// A version before kind 3 skips LevelSetByAdult, and reads the correction as evidence with no
	// answer, which moves nothing.
	old := harder
	old.Kind, old.LevelSetByAdult = kindBeforeCorrection, false
	if got := mustAppend(t, before, old, g); !sameState(t, got, before) {
		t.Errorf("the correction as kind 2 reads it moves the state: %+v", got.Modules[0])
	}

	// In the trial series the correction ends it: no trial answers, and the next answer is a step
	// of Elo, past the series. n keeps the weight of the one answer before it: of one item, a
	// series of one, it weighs 3 / (1 + 2).
	const oneAnswer = 1
	trial := mustAppend(t, lm.State{}, choiceAnswer(choiceSpec{skill: topic("a"), place: 1, difficulty: 3, correct: true, start: ladderBand(2), at: 1000}), g)
	trial = mustAppend(t, trial, lm.LevelCorrection(olympiad, ladderBand(3), ladderBand(2), 2000), g)
	if m := trial.Modules[0]; m.Level != middleOfBand3 || len(m.Trial) != 0 || !m.TrialEndedByAdult || m.LastInTrial || m.AnswerWeight != oneAnswer {
		t.Errorf("a correction in the trial series: %+v", m)
	}
	if r := moduleOf(t, mustRead(t, trial, g, 2000, lm.Grade{}, olympiad), olympiad); !r.TrialOver || r.LastInTrial {
		t.Errorf("the reading after a correction in the series: %+v", r)
	}
	next := mustAppend(t, trial, choiceAnswer(choiceSpec{skill: topic("a"), place: 3, difficulty: 3, correct: true, start: ladderBand(2), at: 3000}), g)
	if m := next.Modules[0]; len(m.Trial) != 0 || m.LastInTrial || !(m.Level > middleOfBand3) || !m.TrialEndedByAdult {
		t.Errorf("an answer after a correction in the series: %+v", m)
	}

	// A correction is the first evidence: the module starts with the level set, its start band the
	// evidence's; one to the band the start stands in keeps the level and ends the series.
	first := mustAppend(t, lm.State{}, lm.LevelCorrection(olympiad, ladderBand(3), ladderBand(2), 1000), g)
	if m := first.Modules[0]; m.Level != middleOfBand3 || m.StartBand != ladderBand(2) || !m.TrialEndedByAdult || len(m.Skills) != 0 {
		t.Errorf("a correction as the first evidence: %+v", m)
	}
	if kept := mustAppend(t, lm.State{}, lm.LevelCorrection(olympiad, ladderBand(1), ladderBand(2), 1000), g); len(kept.Modules) != 1 ||
		kept.Modules[0].Level != 0 || !kept.Modules[0].TrialEndedByAdult {
		t.Errorf("a correction to the band of the start: %+v", kept)
	}
	// In the band already, the correction still ends the series and starts the guard's count anew,
	// and the level stays where it is.
	inBand := mustAppend(t, lm.State{}, choiceAnswer(choiceSpec{skill: topic("a"), place: 1, difficulty: 3, correct: true, start: ladderBand(2), at: 1000}), g)
	level := inBand.Modules[0].Level
	inBand = mustAppend(t, inBand, lm.LevelCorrection(olympiad, ladderBand(1), ladderBand(2), 2000), g)
	if m := inBand.Modules[0]; m.Level != level || len(m.Trial) != 0 || !m.TrialEndedByAdult {
		t.Errorf("a correction to the band of the level in the series: %+v, θ was %v", m, level)
	}
	counted := settled(settledLevel)
	counted.Modules[0].FailuresInARow = guardOneShort
	counted = mustAppend(t, counted, lm.LevelCorrection(olympiad, ladderBand(1), ladderBand(1), 2000), g)
	if m := counted.Modules[0]; m.Level != settledLevel || m.FailuresInARow != 0 {
		t.Errorf("a correction to the band of the level with a count of the guard: %+v", m)
	}
	// A module with no current version moves nothing, as with any evidence.
	if s := mustAppend(t, before, harder, lm.Graphs{olympiad: {Current: false}}); !sameState(t, s, before) {
		t.Errorf("a correction of a module with no current version: %+v", s.Modules[0])
	}
	// Appended and recomputed are the same.
	log := []lm.Evidence{
		choiceAnswer(choiceSpec{skill: topic("a"), place: 1, difficulty: 3, correct: true, start: ladderBand(2), at: 1000}),
		lm.LevelCorrection(olympiad, ladderBand(3), ladderBand(2), 2000),
		choiceAnswer(choiceSpec{skill: topic("a"), place: 3, difficulty: 3, correct: true, start: ladderBand(2), at: 3000}),
	}
	if got, err := lm.Recompute(log, g); err != nil || !sameState(t, got, next) {
		t.Errorf("Recompute of the corrected log: %v", err)
	}
}

// What a correction must be: what LevelCorrection gives, of kind 3 or later; a band above the
// module's current version's highest is its highest.
func TestCorrectionRefused(t *testing.T) {
	g := olympiadGraphs(t, map[string][]int32{topic("a"): {1, 2, 3}})
	good := func() lm.Evidence { return lm.LevelCorrection(olympiad, ladderBand(2), ladderBand(1), rulesAt) }
	for name, change := range map[string]func(*lm.Evidence){
		"kind 2":             func(e *lm.Evidence) { e.Kind = kindBeforeCorrection },
		"another ID":         func(e *lm.Evidence) { e.Task.ID = olympiad + "/task/x" },
		"an ID of no module": func(e *lm.Evidence) { e.Task.ID, e.Task.Band.ID = "level/", "" },
		"a difficulty":       func(e *lm.Evidence) { e.Task.Difficulty = lm.MiddleDifficulty + 1 },
		"a template":         func(e *lm.Evidence) { e.Task.Template = true },
		"the author's key":   func(e *lm.Evidence) { e.Task.KeyProvenance = lm.ProvenanceAuthor },
		"the author's tags":  func(e *lm.Evidence) { e.Task.SkillTagsProvenance = lm.ProvenanceAuthor },
		"a skill":            func(e *lm.Evidence) { e.Task.Primary = []string{topic("a")} },
		"a supporting skill": func(e *lm.Evidence) { e.Task.Supporting = []string{"core.counting"} },
		"an outcome":         func(e *lm.Evidence) { e.Outcome.Kind = lm.OutcomeFirstSubmission },
		"a technique":        func(e *lm.Evidence) { e.TechniqueOnStep = true },
		"a request":          func(e *lm.Evidence) { e.TechniqueOnRequest = true },
		"a leading question": func(e *lm.Evidence) { e.LeadingQuestionOnRequest = true },
		"explain again":      func(e *lm.Evidence) { e.ExplainAgainCount = 1 },
		"a seen solution":    func(e *lm.Evidence) { e.HasSolutionSeenFact = true },
		"a trap":             func(e *lm.Evidence) { e.Traps = []lm.FoundTrap{{ID: olympiad + "/trap/t"}} },
		"an item": func(e *lm.Evidence) {
			e.Task.Items = []lm.TaskItem{{Place: 1}}
			e.Items = []lm.ItemAnswer{{Place: 1}}
		},
	} {
		e := good()
		change(&e)
		if _, err := lm.EncodeEvidence(e); !errors.Is(err, lm.ErrInvalid) {
			t.Errorf("a correction with %s encodes: %v", name, err)
		}
		if _, _, err := lm.Append(settled(0), e, g); !errors.Is(err, lm.ErrInvalid) {
			t.Errorf("a correction with %s appends: %v", name, err)
		}
	}
	far := lm.LevelCorrection(olympiad, lm.BandRef{Place: 4, ID: "7-8"}, ladderBand(1), rulesAt)
	if s, _, err := lm.Append(settled(0), far, g); err != nil || s.Modules[0].Level != middleOfBand3 {
		t.Errorf("a correction to a band the graph does not have: %v, %+v", err, s)
	}

	// With no current version the band is not known, and the correction waits like any evidence.
	if _, _, err := lm.Append(settled(0), far, lm.Graphs{olympiad: {Current: false}}); err != nil {
		t.Errorf("a correction of a module with no current version: %v", err)
	}
}

// stuckSpec is a task the learner was stuck on: its skill, band place and difficulty, and the
// practice time.
type stuckSpec struct {
	skill             string
	place, difficulty int32
	at                int64
}

// stuckOn is a wrong answer of choiceAnswer the learner was stuck on: the solution shown.
func stuckOn(s stuckSpec) lm.Evidence {
	return choiceAnswer(choiceSpec{skill: s.skill, place: s.place, difficulty: s.difficulty, correct: false,
		start: ladderBand(lm.FirstPlace), at: s.at})
}

// The guard: three tasks in a row the model gave as fitting and the learner was stuck on lower the
// level by a band, once a run; a task put right, or one the model took for too hard, is no failure.
func TestGuard(t *testing.T) {
	a := topic("a")
	g := olympiadGraphs(t, map[string][]int32{a: {1, 2, 3}})
	const (
		standaloneFrom = 0.9  // the level standalone's failures start from
		stillHigh      = 0.5  // a level that has not fallen a band from standaloneFrom
		loweredBelow   = -1.5 // a level that has fallen a band from standaloneFrom, its steps besides
		runPast        = lm.GuardFailures + 1
		floorAbove     = 0.01 // how far above the floor the guard leaves the level, at most
		floorBelow     = 0.2  // and how far below it the step of the failure takes it
		// inSeries is how many answers of one item, each weighing 1, leave the trial series open.
		inSeries = int(lm.TrialWeight) - 1
	)
	// standalone's three failures, stuck on: the third lowers θ by a band after its step.
	s := settled(standaloneFrom)
	s.Modules[0].Skills[0].Correction = 0.3
	thetas := []float64{0.8185, 0.7437, 0.6755 - lm.BandStep}
	for i, want := range thetas {
		s = mustAppend(t, s, stuckOn(stuckSpec{skill: a, place: 1, difficulty: 3, at: rulesAt}), g)
		near(t, "θ after failure "+itoa(i+1), s.Modules[0].Level, want)
		if s.Modules[0].FailuresInARow != int32(i+1) {
			t.Errorf("failure %d: the count %d", i+1, s.Modules[0].FailuresInARow)
		}
	}
	// The next task of the band 1-2 is too hard by the model's estimate now: stuck on it, it is no
	// failure, and the count starts anew.
	s = mustAppend(t, s, stuckOn(stuckSpec{skill: a, place: 1, difficulty: 3, at: rulesAt}), g)
	if m := s.Modules[0]; m.FailuresInARow != 0 || !(m.Level > lm.LevelFloor) {
		t.Errorf("stuck on a task the model took for too hard: %+v", m)
	}
	// A choice of two far above the level is too hard too, though its P is never below ½.
	two := settled(standaloneFrom)
	two.Modules[0].FailuresInARow = guardOneShort
	hard := stuckOn(stuckSpec{skill: a, place: 3, difficulty: 5, at: rulesAt})
	hard.Task.Items[0].Slots[0].Options = lm.MinChoiceOptions
	if m := mustAppend(t, two, hard, g).Modules[0]; m.FailuresInARow != 0 || !(m.Level > stillHigh) {
		t.Errorf("stuck on a choice of two far above the level: %+v", m)
	}
	// Failures put right on a rung are none.
	right := settled(standaloneFrom)
	for range lm.GuardFailures {
		right = mustAppend(t, right, putRight(choiceAnswer(choiceSpec{skill: a, place: 1, difficulty: 3, correct: false, start: ladderBand(1), at: rulesAt})), g)
	}
	if m := right.Modules[0]; m.FailuresInARow != 0 || !(m.Level > stillHigh) {
		t.Errorf("three failures put right: %+v", m)
	}

	// A run past the third lowers no more: the model's level must come back to fitting tasks, and a
	// task that is no failure must come between.
	low := settled(standaloneFrom)
	low.Modules[0].Skills[0].Correction = 4
	for range runPast {
		low = mustAppend(t, low, stuckOn(stuckSpec{skill: a, place: 1, difficulty: 3, at: rulesAt}), g)
	}
	if m := low.Modules[0]; m.FailuresInARow != runPast || !(m.Level > lm.LevelFloor && m.Level < loweredBelow) {
		t.Errorf("a run of four: %+v", m)
	}
	low = mustAppend(t, low, choiceAnswer(choiceSpec{skill: a, place: 1, difficulty: 1, correct: true, start: ladderBand(1), at: rulesAt}), g)
	if m := low.Modules[0]; m.FailuresInARow != 0 {
		t.Errorf("a right answer after the run: %+v", m)
	}

	// The floor: no lower than a band below the lowest band, and never upward from below it.
	for _, c := range []struct{ from, at float64 }{{from: -2, at: lm.LevelFloor}, {from: -3, at: -3}} {
		f := settled(c.from)
		f.Modules[0].Skills[0].Correction = 4
		f.Modules[0].FailuresInARow = guardOneShort
		f = mustAppend(t, f, stuckOn(stuckSpec{skill: a, place: 1, difficulty: 1, at: rulesAt}), g)
		if m := f.Modules[0]; !(m.Level <= c.at+floorAbove && m.Level >= c.at-floorBelow) || m.FailuresInARow != lm.GuardFailures {
			t.Errorf("the guard from θ = %v: %+v, want about %v", c.from, m, c.at)
		}
	}

	// A correction starts the count anew, and the series counts nothing.
	corrected := settled(standaloneFrom)
	corrected.Modules[0].FailuresInARow = guardOneShort
	corrected = mustAppend(t, corrected, lm.LevelCorrection(olympiad, ladderBand(2), ladderBand(1), rulesAt), g)
	if corrected.Modules[0].FailuresInARow != 0 {
		t.Errorf("a correction leaves the count at %d", corrected.Modules[0].FailuresInARow)
	}
	var series lm.State
	for i := range inSeries {
		series = mustAppend(t, series, choiceAnswer(choiceSpec{skill: a, place: 1, difficulty: 1, correct: false, start: ladderBand(3), at: int64(i)}), g)
	}
	if m := series.Modules[0]; m.FailuresInARow != 0 || len(m.Trial) != inSeries {
		t.Errorf("stuck in the trial series: %+v", m)
	}
}

// A page of many items is a failure when fewer than half are solved at the end, and the model
// expected at least half right.
func TestGuardPage(t *testing.T) {
	g := workbookGraphs(t)
	add := []string{skill("addition")}
	base := lm.State{Kind: lm.StateKind, Modules: []lm.ModuleState{{
		Module: workbook, Level: 2, AnswerWeight: settledWeight, StartBand: lm.BandRef{Place: 1, ID: "grade-2"}, FailuresInARow: guardOneShort,
	}}}
	// One of four right: stuck. Two of four: no failure.
	if m := mustAppend(t, base, page(pageSpec{items: 4, right: 1, primary: add, at: 1000}), g).Modules[0]; m.FailuresInARow != lm.GuardFailures || !(m.Level < 0) {
		t.Errorf("a page of one right in four: %+v", m)
	}
	if m := mustAppend(t, base, page(pageSpec{items: 4, right: 2, primary: add, at: 1000}), g).Modules[0]; m.FailuresInARow != 0 {
		t.Errorf("a page of two right in four: %+v", m)
	}
	// Wrong at first, right at the end: solved.
	fixed := page(pageSpec{items: 4, right: 1, primary: add, at: 1000})
	for i := 1; i < len(fixed.Items); i++ {
		fixed.Items[i].FinalVerdict, fixed.Items[i].HasFinalVerdict = lm.VerdictCorrect, true
	}
	if m := mustAppend(t, base, fixed, g).Modules[0]; m.FailuresInARow != 0 {
		t.Errorf("a page put right at the end: %+v", m)
	}
	// Right at first, not yet at the end: not solved.
	undone := page(pageSpec{items: 4, right: 4, primary: add, at: 1000})
	for i := 1; i < len(undone.Items); i++ {
		undone.Items[i].FinalVerdict, undone.Items[i].HasFinalVerdict = lm.VerdictNotYet, true
	}
	if m := mustAppend(t, base, undone, g).Modules[0]; m.FailuresInARow != lm.GuardFailures {
		t.Errorf("a page undone at the end: %+v", m)
	}
}

// The guard's state: past the series only, never below 0, and only in a state of kind 2.
func TestGuardState(t *testing.T) {
	g := olympiadGraphs(t, map[string][]int32{topic("a"): {1, 2, 3}})
	const inTrialWeight = 2 // n of a module still in its trial series, below its weight
	inTrial := settled(0)
	inTrial.Modules[0].AnswerWeight, inTrial.Modules[0].LastInTrial, inTrial.Modules[0].FailuresInARow = inTrialWeight, true, 1
	if _, _, err := lm.Append(inTrial, stuckOn(stuckSpec{skill: topic("a"), place: 1, difficulty: 3, at: rulesAt}), g); !errors.Is(err, lm.ErrInvalid) {
		t.Errorf("a count of the guard in the trial series: %v", err)
	}
	negative := settled(0)
	negative.Modules[0].FailuresInARow = -1
	if _, err := lm.EncodeState(negative); !errors.Is(err, lm.ErrInvalid) {
		t.Errorf("a count of the guard below 0: %v", err)
	}
	for _, change := range []func(*lm.ModuleState){
		func(m *lm.ModuleState) { m.FailuresInARow = 1 },
		func(m *lm.ModuleState) { m.TrialEndedByAdult = true },
	} {
		s := settled(0)
		s.Kind = lm.FirstKind
		change(&s.Modules[0])
		if _, err := lm.EncodeState(s); !errors.Is(err, lm.ErrInvalid) {
			t.Errorf("the guard's state in a state of kind 1: %v", err)
		}
	}
	// Kind 1 reads without the fields and refuses them; kind 2 requires them.
	one := settled(0)
	one.Kind = lm.FirstKind
	b, err := lm.EncodeState(one)
	if err != nil {
		t.Fatal(err)
	}
	if s, err := lm.DecodeState(b); err != nil || s.Kind != lm.FirstKind || s.Modules[0].FailuresInARow != 0 {
		t.Errorf("a state of kind 1: %+v, %v", s, err)
	}
	two, err := lm.EncodeState(settled(0))
	if err != nil {
		t.Fatal(err)
	}
	// The module is the last field; its last two fields, 9 and 10 of 0, are 0x48 0x00 0x50 0x00.
	// Moving them into kind 1 or dropping them from kind 2 leaves the module's length to mend.
	if _, err := lm.DecodeState(withKind(two, lm.FirstKind)); !errors.Is(err, lm.ErrCorrupt) {
		t.Errorf("a state of kind 1 with fields 9 and 10: %v", err)
	}
	if _, err := lm.DecodeState(withKind(b, lm.StateKind)); !errors.Is(err, lm.ErrCorrupt) {
		t.Errorf("a state of kind 2 without fields 9 and 10: %v", err)
	}
}

// withKind gives the bytes of a state with another kind: its first field, 0x08 and the kind.
func withKind(b []byte, kind byte) []byte {
	const kindByte = 1 // after the tag 0x08 of field 1
	out := append([]byte{}, b...)
	out[kindByte] = kind
	return out
}

// rulesLog is a learner of grade 1 in olympiad maths, who starts below its lowest band, and in
// school's workbook: the trial series from θ₀ = −2.5, an adult's correction to the band 3-4 and its
// repeat, which changes nothing, three tasks of that band the learner is stuck on, of two skills, which
// lower the level by a band, one more the model now takes for too hard, and a right answer; then the adult
// sets the workbook's band grade-3 before its first page, which ends its trial series.
func rulesLog() []lm.Evidence {
	start := ladderBand(lm.FirstPlace) // the band of grade 1
	const (
		logStart = 1762041600 // 2025-11-02 00:00 UTC
		hour     = 3600       // an hour between pieces of evidence
		wrongAt  = 2          // the third answer of the trial series is wrong
	)
	at := int64(logStart)
	next := func() int64 {
		at += hour
		return at
	}
	var log []lm.Evidence
	for i := range int(lm.TrialWeight) {
		log = append(log, choiceAnswer(choiceSpec{skill: topic("a"), place: 1, difficulty: 3, correct: i != wrongAt, start: start, at: next()}))
	}
	log = append(log,
		lm.LevelCorrection(olympiad, ladderBand(2), start, next()),
		lm.LevelCorrection(olympiad, ladderBand(2), start, next()))
	for i := range lm.GuardFailures + 1 {
		log = append(log, stuckOn(stuckSpec{skill: topic([]string{"a", "b"}[i%2]), place: 2, difficulty: 2, at: next()}))
	}
	log = append(log, choiceAnswer(choiceSpec{skill: topic("a"), place: 1, difficulty: 3, correct: true, start: start, at: next()}))
	grade2 := lm.BandRef{Place: 1, ID: "grade-2"}
	log = append(log,
		lm.LevelCorrection(workbook, lm.BandRef{Place: 2, ID: "grade-3"}, grade2, next()),
		page(pageSpec{items: 4, right: 3, primary: []string{skill("addition")}, at: next()}))
	return log
}

func init() {
	goldenJournals = append(goldenJournals, goldenJournal{name: "rules",
		note:   "a learner of grade 1 in olympiad maths and the workbook: the start below the lowest band, an adult's corrections and the guard",
		graphs: foldGraphs, log: rulesLog, kind: 0})
}
