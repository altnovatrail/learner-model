package learnermodel_test

import (
	"math"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

// The test vectors of standalone: the worked examples of SPEC.md of
// mathtrail/mathtrail-standalone, at commit 6b29dd9 of its main. They hold to the four decimals
// printed there, 1e-4, as SPEC.md asks: a module of olympiad maths whose tasks are a choice of
// five, c = 0.2, counts exactly as standalone while retention is 1. Each answer of a table is
// evidence of a task of one item with one choice slot; the level `1-2`, `3-4` or `5-6` is the band
// at place 1, 2 or 3, and the topic is the task's primary skill. The numbers are standalone's,
// typed in from its tables; no code of standalone is here. Its first example starts the step of Elo
// from no answers, which evidence never does past the trial series, and runs inside the package, in
// step_internal_test.go.

const olympiad = "example/olympiad"

func topic(name string) string { return olympiad + "/skill/" + name }

// standaloneOptions is the number of options of standalone's tasks, a choice of five, c = 0.2.
const standaloneOptions = 5

// standaloneTolerance is what the tables promise: four decimals.
const standaloneTolerance = 1e-4

func ladder() []lm.Band {
	return []lm.Band{
		{Place: 1, ID: "1-2", Grades: []int32{1, 2}},
		{Place: 2, ID: "3-4", Grades: []int32{3, 4}},
		{Place: 3, ID: "5-6", Grades: []int32{5, 6}},
	}
}

func ladderBand(place int32) lm.BandRef {
	return lm.BandRef{Place: place, ID: ladder()[place-lm.FirstPlace].ID}
}

// olympiadGraphs gives the graphs of the olympiad module, each topic taught at the bands named.
func olympiadGraphs(t *testing.T, taught map[string][]int32) lm.Graphs {
	t.Helper()
	var skills []lm.Skill
	var tasks []lm.Task
	for _, id := range sortedKeys(taught) {
		skills = append(skills, lm.Skill{ID: id, Assessability: lm.AssessabilityVerdict})
		for _, place := range taught[id] {
			tasks = append(tasks, lm.Task{ID: olympiad + "/task/x", Band: ladderBand(place), Difficulty: lm.MiddleDifficulty, Primary: []string{id}})
		}
	}
	g, err := lm.NewGraph(olympiad, ladder(), skills, nil, tasks)
	if err != nil {
		t.Fatal(err)
	}
	return lm.Graphs{olympiad: {Graph: g, Current: true}}
}

// choiceSpec is one answer of standalone: its skill, the band place and difficulty of its task,
// whether it is right, the start band and the practice time.
type choiceSpec struct {
	skill             string
	place, difficulty int32
	correct           bool
	start             lm.BandRef
	at                int64
}

// choiceAnswer is the evidence of one answer of standalone: a task of one item, a choice of five.
func choiceAnswer(c choiceSpec) lm.Evidence {
	skillID, place, difficulty, correct, start, at := c.skill, c.place, c.difficulty, c.correct, c.start, c.at
	verdict, outcome := lm.VerdictCorrect, lm.OutcomeFirstSubmission
	if !correct {
		verdict, outcome = lm.VerdictNotYet, lm.OutcomeSolutionShown
	}
	return lm.Evidence{
		Kind: lm.EvidenceKind,
		Task: lm.Task{
			ID: olympiad + "/task/x", Band: ladderBand(place), Difficulty: difficulty, Primary: []string{skillID},
			KeyProvenance: lm.ProvenanceAuthor, SkillTagsProvenance: lm.ProvenanceAuthor,
			Items: []lm.TaskItem{{Place: 1, Slots: []lm.Slot{{Check: lm.CheckKindChoice, Options: standaloneOptions, HasOptions: true}}}},
		},
		StartBand:    start,
		Outcome:      lm.Outcome{Kind: outcome},
		PracticeTime: at,
		Items:        []lm.ItemAnswer{{Place: 1, Submitted: true, FirstVerdict: verdict, HasFirstVerdict: true}},
	}
}

// putRight makes a wrong answer of choiceAnswer one put right on the first rung of the ladder: the
// solution is not shown, and the item is solved at the end.
func putRight(e lm.Evidence) lm.Evidence {
	e.Outcome = lm.Outcome{Kind: lm.OutcomeAfterHint, Rung: 1, HasRung: true}
	e.Items = []lm.ItemAnswer{{
		Place: 1, Submitted: true, FirstVerdict: lm.VerdictNotYet, HasFirstVerdict: true,
		FinalVerdict: lm.VerdictCorrect, HasFinalVerdict: true, WrongBeforeFirstCorrect: 1, HasFirstCorrect: true,
	}}
	return e
}

func near(t *testing.T, what string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > standaloneTolerance {
		t.Errorf("%s = %.6f, want %.4f", what, got, want)
	}
}

func skillOf(t *testing.T, s lm.State, id string) lm.SkillState {
	t.Helper()
	for _, sk := range s.Modules[0].Skills {
		if sk.Skill == id {
			return sk
		}
	}
	t.Fatalf("no state of skill %s", id)
	return lm.SkillState{}
}

// standaloneRow is a row of a table: the task, the answer and θ and δ after it.
type standaloneRow struct {
	topic        string
	place, d     int32
	correct      bool
	theta, delta float64
}

// Three failures: a settled child fails difficulty 3 of `1-2` three times, then answers one of
// difficulty 2. Each failure is put right on the first rung of the ladder, so that the guard, which
// counts the tasks a learner is stuck on, does not lower the level; TestGuard leaves the
// three unsolved.
func TestStandaloneThreeFailures(t *testing.T) {
	const (
		at          = 1000 // the time of every answer and of the practice before them
		skillBefore = 8    // m and the evidence of the skill before the first answer
	)
	a := topic("combinatorics.enumeration")
	g := olympiadGraphs(t, map[string][]int32{a: {1}})
	s := lm.State{Kind: lm.StateKind, Modules: []lm.ModuleState{{
		Module: olympiad, Level: 0.9, AnswerWeight: 20, StartBand: ladderBand(1),
		Skills: []lm.SkillState{{Skill: a, Correction: 0.3, AnswerWeight: skillBefore, EvidenceCount: skillBefore, LastPractice: at, Stability: 1}},
	}}}
	rows := []standaloneRow{
		{topic: a, place: 1, d: 3, correct: false, theta: 0.8185, delta: 0.0672},
		{topic: a, place: 1, d: 3, correct: false, theta: 0.7437, delta: -0.1442},
		{topic: a, place: 1, d: 3, correct: false, theta: 0.6755, delta: -0.3353},
		{topic: a, place: 1, d: 2, correct: true, theta: 0.6910, delta: -0.2924},
	}
	failures := []int32{1, 2, 3, 0}
	for i, r := range rows {
		var err error
		e := choiceAnswer(choiceSpec{skill: r.topic, place: r.place, difficulty: r.d, correct: r.correct, start: ladderBand(1), at: at})
		if !r.correct {
			e = putRight(e)
		}
		s, _, err = lm.Append(s, e, g)
		if err != nil {
			t.Fatal(err)
		}
		sk := skillOf(t, s, a)
		near(t, "θ after answer "+itoa(i+1), s.Modules[0].Level, r.theta)
		near(t, "δ after answer "+itoa(i+1), sk.Correction, r.delta)
		if sk.FailuresInARow != failures[i] || sk.Mastered {
			t.Errorf("answer %d: failures in a row %d, mastered %v; want %d, not mastered", i+1, sk.FailuresInARow, sk.Mastered, failures[i])
		}
		if sk.AnswerWeight != float64(skillBefore+1+i) || sk.EvidenceCount != int64(skillBefore+1+i) || sk.LastPractice != at {
			t.Errorf("answer %d: m %v, evidence %d, last practice %d", i+1, sk.AnswerWeight, sk.EvidenceCount, sk.LastPractice)
		}
	}
}

// The trial series of a child of grade 3, who starts at θ₀ = 2.5 and stands about a level lower;
// every topic is new, and no δ moves. The start is a band below the band of the grade, so
// θ₀ = 2.5 is the start of a learner of grade 5, whose band 5-6 stands at 5.
func TestStandaloneTrialSeries(t *testing.T) {
	start, err := lm.BandByGrade(lm.Grade{Value: 5, Known: true}, ladder())
	if err != nil || start != ladderBand(3) {
		t.Fatalf("BandByGrade(5) = %+v, %v; want the band 5-6", start, err)
	}
	rows := []standaloneRow{
		{topic: topic("logic.ordering"), place: 2, d: 2, correct: false, theta: 0.6407, delta: 0},
		{topic: topic("combinatorics.enumeration"), place: 1, d: 3, correct: true, theta: 1.1189, delta: 0},
		{topic: topic("logic.knights_liars"), place: 2, d: 1, correct: false, theta: 0.2536, delta: 0},
		{topic: topic("counting.gaps"), place: 1, d: 2, correct: true, theta: 0.4532, delta: 0},
		{topic: topic("time.clocks"), place: 1, d: 2, correct: true, theta: 0.6039, delta: 0},
	}
	taught := map[string][]int32{}
	for _, r := range rows {
		taught[r.topic] = []int32{r.place}
	}
	g := olympiadGraphs(t, taught)

	var s lm.State
	var log []lm.Evidence
	for i, r := range rows {
		var e lm.Evidence
		s, e, err = lm.Append(s, choiceAnswer(choiceSpec{skill: r.topic, place: r.place, difficulty: r.d, correct: r.correct, start: start, at: int64(i)}), g)
		if err != nil {
			t.Fatal(err)
		}
		log = append(log, e)
		m := s.Modules[0]
		near(t, "θ after trial answer "+itoa(i+1), m.Level, r.theta)
		if sk := skillOf(t, s, r.topic); sk.Correction != 0 || sk.AnswerWeight != 1 {
			t.Errorf("trial answer %d: δ %v and m %v, want 0 and 1", i+1, sk.Correction, sk.AnswerWeight)
		}
		wantTrial := i + 1
		if i == len(rows)-1 {
			wantTrial = 0 // the fifth answer ends the series
		}
		if len(m.Trial) != wantTrial || m.AnswerWeight != float64(i+1) || m.StartBand != start {
			t.Errorf("trial answer %d: %d trial answers, n %v, start %+v", i+1, len(m.Trial), m.AnswerWeight, m.StartBand)
		}
	}
	again, err := lm.Recompute(log, g)
	if err != nil {
		t.Fatal(err)
	}
	if !sameState(t, again, s) {
		t.Errorf("Recompute gives another state than the appends")
	}
}

// The floor of θ's step, a mastery at a band below the task's, one failure that holds it and a
// second in a row that takes it away. The topic is taught at all three bands.
func TestStandaloneFloorAndMastery(t *testing.T) {
	const (
		moduleBefore = 79 // n before the first answer
		skillBefore  = 9  // m and the evidence of the skill before it
	)
	a := topic("combinatorics.enumeration")
	g := olympiadGraphs(t, map[string][]int32{a: {1, 2, 3}})
	s := lm.State{Kind: lm.StateKind, Modules: []lm.ModuleState{{
		Module: olympiad, Level: 1.2, AnswerWeight: moduleBefore, StartBand: ladderBand(1),
		Skills: []lm.SkillState{{Skill: a, Correction: 0.5, AnswerWeight: skillBefore, EvidenceCount: skillBefore, Stability: 1}},
	}}}
	rows := []standaloneRow{
		{topic: a, place: 2, d: 1, correct: true, theta: 1.2093, delta: 0.5511},
		{topic: a, place: 1, d: 4, correct: false, theta: 1.1720, delta: 0.3524},
		{topic: a, place: 1, d: 3, correct: false, theta: 1.1292, delta: 0.1312},
	}
	mastered := []bool{true, true, false}
	for i, r := range rows {
		var err error
		s, _, err = lm.Append(s, choiceAnswer(choiceSpec{skill: r.topic, place: r.place, difficulty: r.d, correct: r.correct, start: ladderBand(1), at: 0}), g)
		if err != nil {
			t.Fatal(err)
		}
		sk := skillOf(t, s, a)
		near(t, "θ after answer "+itoa(i+1), s.Modules[0].Level, r.theta)
		near(t, "δ after answer "+itoa(i+1), sk.Correction, r.delta)
		wantBand := int32(lm.FirstPlace) // the lowest band, below the task's
		if !mastered[i] {
			wantBand = 0
		}
		if sk.Mastered != mastered[i] || sk.MasteryBand != wantBand {
			t.Errorf("answer %d: mastered %v at band %d, want %v at %d", i+1, sk.Mastered, sk.MasteryBand, mastered[i], wantBand)
		}
		if s.Modules[0].AnswerWeight != float64(moduleBefore+1+i) || sk.AnswerWeight != float64(skillBefore+1+i) {
			t.Errorf("answer %d: n %v and m %v", i+1, s.Modules[0].AnswerWeight, sk.AnswerWeight)
		}
	}
}
