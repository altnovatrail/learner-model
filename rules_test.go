package learnermodel_test

import (
	"errors"
	"fmt"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

// The rules in place of a calibration: the start a band below the band of the
// grade, the band the level stands in, an adult's correction and the guard.

// settled is a learner of olympiad maths past the trial series, at a level, with one skill a.
func settled(level float64) lm.State {
	return lm.State{Kind: lm.StateKind, Modules: []lm.ModuleState{{
		Module: olympiad, Level: level, AnswerWeight: 20, StartBand: ladderBand(1),
		Skills: []lm.SkillState{{Skill: topic("a"), AnswerWeight: 8, EvidenceCount: 8, LastPractice: 1000, Stability: 1}},
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
	m := moduleOf(t, mustRead(t, s, g, 1000, grade, olympiad), olympiad)
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
	// series would leave θ at θ₀ = −2.5, and a right one lifts it from there.
	s := mustAppend(t, lm.State{}, choiceAnswer(topic("a"), 1, 3, true, ladderBand(1), 1000), g)
	if m := s.Modules[0]; m.StartBand != ladderBand(1) || !(m.Level > -2.5 && m.Level < 0) {
		t.Errorf("after a right answer from the start of grade 1: start %+v, θ %v", m.StartBand, m.Level)
	}
}

// The band the level stands in: the band whose middle is nearest, the upper of two as near, the
// highest above it, none more than half a band below the lowest, and none without bands to name.
func TestLevelBand(t *testing.T) {
	g := olympiadGraphs(t, map[string][]int32{topic("a"): {1, 2, 3}})
	for level, want := range map[float64]string{
		-1.2500001: "none", -1.25: "1/1-2", 0: "1/1-2", 1.2499999: "1/1-2", 1.25: "2/3-4",
		3.75: "3/5-6", 5: "3/5-6", 100: "3/5-6",
	} {
		if got := levelBandOf(t, settled(level), g, lm.Grade{}); got != want {
			t.Errorf("θ = %v stands in band %s, want %s", level, got, want)
		}
	}
	// A module with no current version names no band, nor does one with no bands.
	waiting := lm.Graphs{olympiad: {Current: false}}
	if m := mustRead(t, settled(0), waiting, 1000, lm.Grade{}).Modules[0]; m.HasLevelBand || m.LevelBand != (lm.BandRef{}) {
		t.Errorf("a module with no current version: %+v", m)
	}
	bare, err := lm.NewGraph(olympiad, nil, []lm.Skill{{ID: topic("a"), Assessability: lm.AssessabilityVerdict}}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	noBands := lm.Graphs{olympiad: {Graph: bare, Current: true}}
	for _, s := range []lm.State{settled(0), {}} {
		if m := moduleOf(t, mustRead(t, s, noBands, 1000, lm.Grade{}, olympiad), olympiad); m.HasLevelBand {
			t.Errorf("a module with no bands: %+v", m)
		}
	}
}

// An adult's correction sets the level at the middle of the band set, in one step, and a repeat
// changes nothing; it moves no skill, ends the trial series and starts the guard's count anew.
func TestCorrection(t *testing.T) {
	g := olympiadGraphs(t, map[string][]int32{topic("a"): {1, 2, 3}})
	harder := lm.LevelCorrection(olympiad, ladderBand(2), ladderBand(3), 5000)
	before := settled(0.3)
	before.Modules[0].FailuresInARow = 2
	s, e, err := lm.Append(before, harder, g)
	if err != nil {
		t.Fatal(err)
	}
	m := s.Modules[0]
	if m.Level != 2.5 || m.AnswerWeight != 20 || m.FailuresInARow != 0 || m.LastInTrial || m.TrialEndedByAdult ||
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
	old.Kind, old.LevelSetByAdult = 2, false
	if got := mustAppend(t, before, old, g); !sameState(t, got, before) {
		t.Errorf("the correction as kind 2 reads it moves the state: %+v", got.Modules[0])
	}

	// In the trial series the correction ends it: no trial answers, and the next answer is a step
	// of Elo, past the series.
	trial := mustAppend(t, lm.State{}, choiceAnswer(topic("a"), 1, 3, true, ladderBand(2), 1000), g)
	trial = mustAppend(t, trial, lm.LevelCorrection(olympiad, ladderBand(3), ladderBand(2), 2000), g)
	if m := trial.Modules[0]; m.Level != 5 || len(m.Trial) != 0 || !m.TrialEndedByAdult || m.LastInTrial || m.AnswerWeight != 1 {
		t.Errorf("a correction in the trial series: %+v", m)
	}
	if r := moduleOf(t, mustRead(t, trial, g, 2000, lm.Grade{}, olympiad), olympiad); !r.TrialOver || r.LastInTrial {
		t.Errorf("the reading after a correction in the series: %+v", r)
	}
	next := mustAppend(t, trial, choiceAnswer(topic("a"), 3, 3, true, ladderBand(2), 3000), g)
	if m := next.Modules[0]; len(m.Trial) != 0 || m.LastInTrial || !(m.Level > 5) || !m.TrialEndedByAdult {
		t.Errorf("an answer after a correction in the series: %+v", m)
	}

	// A correction is the first evidence: the module starts with the level set, its start band the
	// evidence's; one to the band the start stands in keeps the level and ends the series.
	first := mustAppend(t, lm.State{}, lm.LevelCorrection(olympiad, ladderBand(3), ladderBand(2), 1000), g)
	if m := first.Modules[0]; m.Level != 5 || m.StartBand != ladderBand(2) || !m.TrialEndedByAdult || len(m.Skills) != 0 {
		t.Errorf("a correction as the first evidence: %+v", m)
	}
	if kept := mustAppend(t, lm.State{}, lm.LevelCorrection(olympiad, ladderBand(1), ladderBand(2), 1000), g); len(kept.Modules) != 1 ||
		kept.Modules[0].Level != 0 || !kept.Modules[0].TrialEndedByAdult {
		t.Errorf("a correction to the band of the start: %+v", kept)
	}
	// In the band already, the correction still ends the series and starts the guard's count anew,
	// and the level stays where it is.
	inBand := mustAppend(t, lm.State{}, choiceAnswer(topic("a"), 1, 3, true, ladderBand(2), 1000), g)
	level := inBand.Modules[0].Level
	inBand = mustAppend(t, inBand, lm.LevelCorrection(olympiad, ladderBand(1), ladderBand(2), 2000), g)
	if m := inBand.Modules[0]; m.Level != level || len(m.Trial) != 0 || !m.TrialEndedByAdult {
		t.Errorf("a correction to the band of the level in the series: %+v, θ was %v", m, level)
	}
	counted := settled(0.3)
	counted.Modules[0].FailuresInARow = 2
	counted = mustAppend(t, counted, lm.LevelCorrection(olympiad, ladderBand(1), ladderBand(1), 2000), g)
	if m := counted.Modules[0]; m.Level != 0.3 || m.FailuresInARow != 0 {
		t.Errorf("a correction to the band of the level with a count of the guard: %+v", m)
	}
	// A module with no current version moves nothing, as with any evidence.
	if s := mustAppend(t, before, harder, lm.Graphs{olympiad: {Current: false}}); !sameState(t, s, before) {
		t.Errorf("a correction of a module with no current version: %+v", s.Modules[0])
	}
	// Appended and recomputed are the same.
	log := []lm.Evidence{
		choiceAnswer(topic("a"), 1, 3, true, ladderBand(2), 1000),
		lm.LevelCorrection(olympiad, ladderBand(3), ladderBand(2), 2000),
		choiceAnswer(topic("a"), 3, 3, true, ladderBand(2), 3000),
	}
	if got, err := lm.Recompute(log, g); err != nil || !sameState(t, got, next) {
		t.Errorf("Recompute of the corrected log: %v", err)
	}
}

// What a correction must be: what LevelCorrection gives, of kind 3 or later; a band above the
// module's current version's highest is its highest.
func TestCorrectionRefused(t *testing.T) {
	g := olympiadGraphs(t, map[string][]int32{topic("a"): {1, 2, 3}})
	good := func() lm.Evidence { return lm.LevelCorrection(olympiad, ladderBand(2), ladderBand(1), 1000) }
	for name, change := range map[string]func(*lm.Evidence){
		"kind 2":             func(e *lm.Evidence) { e.Kind = 2 },
		"another ID":         func(e *lm.Evidence) { e.Task.ID = olympiad + "/task/x" },
		"an ID of no module": func(e *lm.Evidence) { e.Task.ID, e.Task.Band.ID = "level/", "" },
		"a difficulty":       func(e *lm.Evidence) { e.Task.Difficulty = 4 },
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
	far := lm.LevelCorrection(olympiad, lm.BandRef{Place: 4, ID: "7-8"}, ladderBand(1), 1000)
	if s, _, err := lm.Append(settled(0), far, g); err != nil || s.Modules[0].Level != 5 {
		t.Errorf("a correction to a band the graph does not have: %v, %+v", err, s)
	}

	// With no current version the band is not known, and the correction waits like any evidence.
	if _, _, err := lm.Append(settled(0), far, lm.Graphs{olympiad: {Current: false}}); err != nil {
		t.Errorf("a correction of a module with no current version: %v", err)
	}
}

// stuckOn is a wrong answer of choiceAnswer the learner was stuck on: the solution shown.
func stuckOn(skillID string, place, difficulty int32, at int64) lm.Evidence {
	return choiceAnswer(skillID, place, difficulty, false, ladderBand(1), at)
}

// The guard: three tasks in a row the model gave as fitting and the learner was stuck on lower the
// level by a band, once a run; a task put right, or one the model took for too hard, is no failure.
func TestGuard(t *testing.T) {
	a := topic("a")
	g := olympiadGraphs(t, map[string][]int32{a: {1, 2, 3}})
	// standalone's three failures, stuck on: the third lowers θ by a band after its step.
	s := settled(0.9)
	s.Modules[0].Skills[0].Correction = 0.3
	thetas := []float64{0.8185, 0.7437, 0.6755 - 2.5}
	for i, want := range thetas {
		s = mustAppend(t, s, stuckOn(a, 1, 3, 1000), g)
		near(t, "θ after failure "+itoa(i+1), s.Modules[0].Level, want)
		if s.Modules[0].FailuresInARow != int32(i+1) {
			t.Errorf("failure %d: the count %d", i+1, s.Modules[0].FailuresInARow)
		}
	}
	// The next task of the band 1-2 is too hard by the model's estimate now: stuck on it, it is no
	// failure, and the count starts anew.
	s = mustAppend(t, s, stuckOn(a, 1, 3, 1000), g)
	if m := s.Modules[0]; m.FailuresInARow != 0 || !(m.Level > -2.5) {
		t.Errorf("stuck on a task the model took for too hard: %+v", m)
	}
	// A choice of two far above the level is too hard too, though its P is never below ½.
	two := settled(0.9)
	two.Modules[0].FailuresInARow = 2
	hard := stuckOn(a, 3, 5, 1000)
	hard.Task.Items[0].Slots[0].Options = 2
	if m := mustAppend(t, two, hard, g).Modules[0]; m.FailuresInARow != 0 || !(m.Level > 0.5) {
		t.Errorf("stuck on a choice of two far above the level: %+v", m)
	}
	// Failures put right on a rung are none.
	right := settled(0.9)
	for range 3 {
		right = mustAppend(t, right, putRight(choiceAnswer(a, 1, 3, false, ladderBand(1), 1000)), g)
	}
	if m := right.Modules[0]; m.FailuresInARow != 0 || !(m.Level > 0.5) {
		t.Errorf("three failures put right: %+v", m)
	}

	// A run past the third lowers no more: the model's level must come back to fitting tasks, and a
	// task that is no failure must come between.
	low := settled(0.9)
	low.Modules[0].Skills[0].Correction = 4
	for range 4 {
		low = mustAppend(t, low, stuckOn(a, 1, 3, 1000), g)
	}
	if m := low.Modules[0]; m.FailuresInARow != 4 || !(m.Level > -2.5 && m.Level < -1.5) {
		t.Errorf("a run of four: %+v", m)
	}
	low = mustAppend(t, low, choiceAnswer(a, 1, 1, true, ladderBand(1), 1000), g)
	if m := low.Modules[0]; m.FailuresInARow != 0 {
		t.Errorf("a right answer after the run: %+v", m)
	}

	// The floor: no lower than a band below the lowest band, and never upward from below it.
	for _, c := range []struct{ from, at float64 }{{-2, -2.5}, {-3, -3}} {
		f := settled(c.from)
		f.Modules[0].Skills[0].Correction = 4
		f.Modules[0].FailuresInARow = 2
		f = mustAppend(t, f, stuckOn(a, 1, 1, 1000), g)
		if m := f.Modules[0]; !(m.Level <= c.at+0.01 && m.Level >= c.at-0.2) || m.FailuresInARow != 3 {
			t.Errorf("the guard from θ = %v: %+v, want about %v", c.from, m, c.at)
		}
	}

	// A correction starts the count anew, and the series counts nothing.
	corrected := settled(0.9)
	corrected.Modules[0].FailuresInARow = 2
	corrected = mustAppend(t, corrected, lm.LevelCorrection(olympiad, ladderBand(2), ladderBand(1), 1000), g)
	if corrected.Modules[0].FailuresInARow != 0 {
		t.Errorf("a correction leaves the count at %d", corrected.Modules[0].FailuresInARow)
	}
	var series lm.State
	for i := range 4 {
		series = mustAppend(t, series, choiceAnswer(a, 1, 1, false, ladderBand(3), int64(i)), g)
	}
	if m := series.Modules[0]; m.FailuresInARow != 0 || len(m.Trial) != 4 {
		t.Errorf("stuck in the trial series: %+v", m)
	}
}

// A page of many items is a failure when fewer than half are solved at the end, and the model
// expected at least half right.
func TestGuardPage(t *testing.T) {
	g := workbookGraphs(t)
	add := []string{skill("addition")}
	base := lm.State{Kind: lm.StateKind, Modules: []lm.ModuleState{{
		Module: workbook, Level: 2, AnswerWeight: 20, StartBand: lm.BandRef{Place: 1, ID: "grade-2"}, FailuresInARow: 2,
	}}}
	// One of four right: stuck. Two of four: no failure.
	if m := mustAppend(t, base, page(4, 1, add, 1000), g).Modules[0]; m.FailuresInARow != 3 || !(m.Level < 0) {
		t.Errorf("a page of one right in four: %+v", m)
	}
	if m := mustAppend(t, base, page(4, 2, add, 1000), g).Modules[0]; m.FailuresInARow != 0 {
		t.Errorf("a page of two right in four: %+v", m)
	}
	// Wrong at first, right at the end: solved.
	fixed := page(4, 1, add, 1000)
	for i := 1; i < 4; i++ {
		fixed.Items[i].FinalVerdict, fixed.Items[i].HasFinalVerdict = lm.VerdictCorrect, true
	}
	if m := mustAppend(t, base, fixed, g).Modules[0]; m.FailuresInARow != 0 {
		t.Errorf("a page put right at the end: %+v", m)
	}
	// Right at first, not yet at the end: not solved.
	undone := page(4, 4, add, 1000)
	for i := 1; i < 4; i++ {
		undone.Items[i].FinalVerdict, undone.Items[i].HasFinalVerdict = lm.VerdictNotYet, true
	}
	if m := mustAppend(t, base, undone, g).Modules[0]; m.FailuresInARow != 3 {
		t.Errorf("a page undone at the end: %+v", m)
	}
}

// The guard's state: past the series only, never below 0, and only in a state of kind 2.
func TestGuardState(t *testing.T) {
	g := olympiadGraphs(t, map[string][]int32{topic("a"): {1, 2, 3}})
	inTrial := settled(0)
	inTrial.Modules[0].AnswerWeight, inTrial.Modules[0].LastInTrial, inTrial.Modules[0].FailuresInARow = 2, true, 1
	if _, _, err := lm.Append(inTrial, stuckOn(topic("a"), 1, 3, 1000), g); !errors.Is(err, lm.ErrInvalid) {
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
		s.Kind = 1
		change(&s.Modules[0])
		if _, err := lm.EncodeState(s); !errors.Is(err, lm.ErrInvalid) {
			t.Errorf("the guard's state in a state of kind 1: %v", err)
		}
	}
	// Kind 1 reads without the fields and refuses them; kind 2 requires them.
	one := settled(0)
	one.Kind = 1
	b, err := lm.EncodeState(one)
	if err != nil {
		t.Fatal(err)
	}
	if s, err := lm.DecodeState(b); err != nil || s.Kind != 1 || s.Modules[0].FailuresInARow != 0 {
		t.Errorf("a state of kind 1: %+v, %v", s, err)
	}
	two, err := lm.EncodeState(settled(0))
	if err != nil {
		t.Fatal(err)
	}
	// The module is the last field; its last two fields, 9 and 10 of 0, are 0x48 0x00 0x50 0x00.
	// Moving them into kind 1 or dropping them from kind 2 leaves the module's length to mend.
	if _, err := lm.DecodeState(withKind(two, 1)); !errors.Is(err, lm.ErrCorrupt) {
		t.Errorf("a state of kind 1 with fields 9 and 10: %v", err)
	}
	if _, err := lm.DecodeState(withKind(b, 2)); !errors.Is(err, lm.ErrCorrupt) {
		t.Errorf("a state of kind 2 without fields 9 and 10: %v", err)
	}
}

// withKind gives the bytes of a state with another kind: its first field, 0x08 and the kind.
func withKind(b []byte, kind byte) []byte {
	out := append([]byte{}, b...)
	out[1] = kind
	return out
}

// rulesLog is a learner of grade 1 in olympiad maths, who starts below its lowest band, and in
// school's workbook: the trial series from θ₀ = −2.5, an adult's correction to the band 3-4 and its
// repeat, which changes nothing, three tasks of that band the learner is stuck on, of two skills, which
// lower the level by a band, one more the model now takes for too hard, and a right answer; then the adult
// sets the workbook's band grade-3 before its first page, which ends its trial series.
func rulesLog() []lm.Evidence {
	start := ladderBand(1)
	at := int64(1762041600)
	next := func() int64 {
		at += 3600
		return at
	}
	var log []lm.Evidence
	for i := range 5 {
		log = append(log, choiceAnswer(topic("a"), 1, 3, i != 2, start, next()))
	}
	log = append(log,
		lm.LevelCorrection(olympiad, ladderBand(2), start, next()),
		lm.LevelCorrection(olympiad, ladderBand(2), start, next()))
	for i := range 4 {
		log = append(log, stuckOn(topic([]string{"a", "b"}[i%2]), 2, 2, next()))
	}
	log = append(log, choiceAnswer(topic("a"), 1, 3, true, start, next()))
	grade2 := lm.BandRef{Place: 1, ID: "grade-2"}
	log = append(log,
		lm.LevelCorrection(workbook, lm.BandRef{Place: 2, ID: "grade-3"}, grade2, next()),
		page(4, 3, []string{skill("addition")}, next()))
	return log
}

func init() {
	goldenJournals = append(goldenJournals, goldenJournal{"rules",
		"a learner of grade 1 in olympiad maths and the workbook: the start below the lowest band, an adult's corrections and the guard",
		foldGraphs, rulesLog, 0})
}
