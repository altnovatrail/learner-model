package learnermodel_test

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

// Reading the state at a time t₀ (docs/architecture.md, sections 3.1–3.4 and 3.7): Read and
// Readout.Task. Every value here is synthetic.

const day = 86400

// The anchor skills of the examples of docs/architecture.md, section 3.4: school maths of grade 3
// teaches division with remainder and fractions of a whole, and olympiad maths relies on them.
const (
	school   = "example/school"
	division = "core.division-with-remainder"
	fraction = "core.fractions-of-a-whole"
	nobody   = "core.nobody-teaches"
	only     = "core.only-school-names"
)

func schoolSkill(name string) string { return school + "/skill/" + name }

// anchorGraphs gives the graphs of school maths and of olympiad maths with anchor skills: in school
// maths division is taught at band 1 and teaches its anchor; fractions has a verdict but no task,
// so it is taught at no band, and the essay has a rubric — both teach fractions of a whole, and
// neither may be mastered. Olympiad maths names the anchors as prerequisites.
func anchorGraphs(t *testing.T) lm.Graphs {
	t.Helper()
	bands := []lm.Band{{Place: 1, ID: "grade-3", Grades: []int32{3}}, {Place: 2, ID: "grade-4", Grades: []int32{4}}}
	sg, err := lm.NewGraph(school, bands, []lm.Skill{
		{ID: schoolSkill("division"), Assessability: lm.AssessabilityVerdict, Teaches: []string{division, only}},
		{ID: schoolSkill("essay"), Assessability: lm.AssessabilityRubric, Teaches: []string{fraction}},
		{ID: schoolSkill("fractions"), Assessability: lm.AssessabilityVerdict, Teaches: []string{fraction}},
	}, nil, []lm.Task{{ID: school + "/task/d", Band: lm.BandRef{Place: 1, ID: "grade-3"}, Difficulty: 3,
		Primary: []string{schoolSkill("division")}}})
	if err != nil {
		t.Fatal(err)
	}
	og, err := lm.NewGraph(olympiad, ladder(), []lm.Skill{
		{ID: topic("a"), Assessability: lm.AssessabilityVerdict, Prerequisites: []string{division, fraction, nobody, topic("b")}},
		{ID: topic("b"), Assessability: lm.AssessabilityVerdict},
	}, nil, []lm.Task{{ID: olympiad + "/task/x", Band: ladderBand(2), Difficulty: 3, Primary: []string{topic("a")}}})
	if err != nil {
		t.Fatal(err)
	}
	return lm.Graphs{school: {Graph: sg, Current: true}, olympiad: {Graph: og, Current: true}}
}

// schoolState is a state of school maths where division has been practised: mastered at band 1, or
// not.
func schoolState(mastered bool) lm.State {
	sk := lm.SkillState{Skill: schoolSkill("division"), Correction: 0.3, AnswerWeight: 8, EvidenceCount: 4,
		LastPractice: 1000, Stability: 2}
	if mastered {
		sk.Mastered, sk.MasteryBand = true, 1
	}
	return lm.State{Kind: lm.StateKind, Modules: []lm.ModuleState{{Module: school, Level: 2, AnswerWeight: 8,
		StartBand: lm.BandRef{Place: 1, ID: "grade-3"}, Skills: []lm.SkillState{sk}}}}
}

func mustRead(t *testing.T, s lm.State, g lm.Graphs, at int64, grade lm.Grade, engine ...string) lm.Readout {
	t.Helper()
	r, err := lm.Read(s, g, at, grade, engine)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func anchorOf(t *testing.T, r lm.Readout, id string) lm.AnchorReadout {
	t.Helper()
	for _, a := range r.Anchors {
		if a.Anchor == id {
			return a
		}
	}
	t.Fatalf("no anchor %s in %+v", id, r.Anchors)
	return lm.AnchorReadout{}
}

func hasAnchor(r lm.Readout, id string) bool {
	return slices.ContainsFunc(r.Anchors, func(a lm.AnchorReadout) bool { return a.Anchor == id })
}

func wantAnchor(t *testing.T, what string, a lm.AnchorReadout, status lm.AnchorStatus, reason lm.AnchorReason) {
	t.Helper()
	if a.Status != status || a.Reason != reason {
		t.Errorf("%s: %s is %d for %d, want %d for %d", what, a.Anchor, a.Status, a.Reason, status, reason)
	}
}

// The examples of docs/architecture.md, section 3.4, with each status and reason.
func TestAnchors(t *testing.T) {
	g := anchorGraphs(t)
	grade3 := lm.Grade{Value: 3, Known: true}

	// A child of grade 3 with both maths, school maths not yet begun.
	r := mustRead(t, lm.State{}, g, 0, grade3, olympiad, school)
	wantAnchor(t, "no evidence", anchorOf(t, r, division), lm.AnchorHasNot, lm.AnchorReasonNoEvidence)
	wantAnchor(t, "skills without a verdict", anchorOf(t, r, fraction), lm.AnchorUnknown, lm.AnchorReasonNoVerdictTeacher)
	wantAnchor(t, "no one teaches", anchorOf(t, r, nobody), lm.AnchorUnknown, lm.AnchorReasonNoTeacher)
	if a := anchorOf(t, r, fraction); len(a.Teachers) != 2 || a.Teachers[0].Skill != schoolSkill("essay") ||
		a.Teachers[1].HighestBand.Place != 0 || a.Teachers[0].Practised {
		t.Errorf("the teachers of fractions: %+v", a.Teachers)
	}
	if a := anchorOf(t, r, division); len(a.Teachers) != 1 ||
		a.Teachers[0] != (lm.TeacherReadout{Skill: schoolSkill("division"), HighestBand: lm.BandRef{Place: 1, ID: "grade-3"}}) {
		t.Errorf("the teacher of division without practice: %+v", a.Teachers)
	}

	// Division practised, not mastered; then mastered.
	r = mustRead(t, schoolState(false), g, 2000, grade3, olympiad, school)
	wantAnchor(t, "not mastered", anchorOf(t, r, division), lm.AnchorHasNot, lm.AnchorReasonNotMastered)
	r = mustRead(t, schoolState(true), g, 2000, grade3, olympiad, school)
	a := anchorOf(t, r, division)
	wantAnchor(t, "mastered", a, lm.AnchorHas, lm.AnchorReasonNone)
	if tr := a.Teachers[0]; !tr.Practised || !tr.Mastered || tr.MasteryBand.ID != "grade-3" || tr.Estimate != 2.3 ||
		!(tr.Retention < 1) || tr.Review.At == 0 {
		t.Errorf("the teacher of division, mastered: %+v", tr)
	}

	// A learner of grade 5 with olympiad maths alone: no one teaches; the anchor school maths alone
	// names is not read.
	r = mustRead(t, lm.State{}, g, 0, lm.Grade{Value: 5, Known: true}, olympiad)
	wantAnchor(t, "olympiad alone", anchorOf(t, r, division), lm.AnchorUnknown, lm.AnchorReasonNoTeacher)
	if hasAnchor(r, only) {
		t.Errorf("an anchor of a module not the engine's: %+v", r.Anchors)
	}
	r = mustRead(t, lm.State{}, g, 0, grade3, olympiad, school)
	wantAnchor(t, "named by school maths", anchorOf(t, r, only), lm.AnchorHasNot, lm.AnchorReasonNoEvidence)

	// School maths turned off: a skill mastered there still holds the anchor; one not mastered
	// leaves it unknown.
	r = mustRead(t, schoolState(true), g, 2000, grade3, olympiad)
	wantAnchor(t, "a module turned off, mastered", anchorOf(t, r, division), lm.AnchorHas, lm.AnchorReasonNone)
	r = mustRead(t, schoolState(false), g, 2000, grade3, olympiad)
	wantAnchor(t, "a module turned off, not mastered", anchorOf(t, r, division), lm.AnchorUnknown, lm.AnchorReasonNoTeacher)

	// A module with no current version gives no graph, and its mastery is not seen.
	gone := anchorGraphs(t)
	gone[school] = lm.ModuleGraph{}
	r = mustRead(t, schoolState(true), gone, 2000, grade3, olympiad)
	if a := anchorOf(t, r, division); a.Status != lm.AnchorUnknown || len(a.Teachers) != 0 {
		t.Errorf("a module with no current version: %+v", a)
	}
	if len(r.Modules) != 2 || r.Modules[1].Module != school || r.Modules[1].Skills[0].MasteryBand != (lm.BandRef{Place: 1}) {
		t.Errorf("a module of the state with no current version: %+v", r.Modules)
	}

	// A module of the engine with no evidence and no current version is left out.
	r = mustRead(t, lm.State{}, gone, 0, grade3, olympiad, school)
	if len(r.Modules) != 1 || r.Modules[0].Module != olympiad {
		t.Errorf("a module of the engine with neither a state nor a current version: %+v", r.Modules)
	}
}

// A module of the engine with no tasks may have no bands (the format's section 11): read with no
// state, it gives no starting point rather than an error, so that the engine still gets its
// reading, and its skills still teach their anchors (Д15 in RUN.md).
func TestReadModuleWithNoBands(t *testing.T) {
	const lessons = "example/lessons"
	spoken := "core.spoken-numbers"
	lg, err := lm.NewGraph(lessons, nil, []lm.Skill{
		{ID: lessons + "/skill/counting-aloud", Assessability: lm.AssessabilityVerdict, Teaches: []string{spoken}},
	}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	g := anchorGraphs(t)
	g[lessons] = lm.ModuleGraph{Graph: lg, Current: true}
	for _, grade := range []lm.Grade{{Value: 3, Known: true}, {}} {
		r := mustRead(t, lm.State{}, g, 0, grade, olympiad, lessons)
		if len(r.Modules) != 2 || r.Modules[0].Module != lessons || r.Modules[1].Module != olympiad {
			t.Fatalf("grade %+v: the modules read: %+v", grade, r.Modules)
		}
		if m := r.Modules[0]; m.HasEvidence || m.StartBand != (lm.BandRef{}) || m.TrialOver || m.LastInTrial ||
			len(m.Skills) != 0 || len(m.Traps) != 0 {
			t.Errorf("grade %+v: the module with no bands: %+v", grade, m)
		}
		if m := r.Modules[1]; m.HasEvidence || m.StartBand.Place < 1 {
			t.Errorf("grade %+v: the module beside it: %+v", grade, m)
		}
		a := anchorOf(t, r, spoken)
		wantAnchor(t, "a teacher with no band", a, lm.AnchorUnknown, lm.AnchorReasonNoVerdictTeacher)
		if len(a.Teachers) != 1 || a.Teachers[0].HighestBand != (lm.BandRef{}) || a.Teachers[0].Practised {
			t.Errorf("grade %+v: the teacher with no band: %+v", grade, a.Teachers)
		}
		if _, err := r.Task(lm.Task{ID: lessons + "/task/x", Difficulty: 3, Band: lm.BandRef{Place: 1}}); !errors.Is(err, lm.ErrInvalid) ||
			!strings.Contains(err.Error(), "has no bands") {
			t.Errorf("grade %+v: a candidate of the module with no bands: %v", grade, err)
		}
	}
}

// The two flags of the trial series (R104 of spec): before it, in it, at its last evidence and
// after it; evidence that moves nothing leaves them as they are.
func TestTrialFlags(t *testing.T) {
	g := olympiadGraphs(t, map[string][]int32{topic("a"): {1, 2, 3}})
	grade3 := lm.Grade{Value: 3, Known: true}
	flags := func(what string, s lm.State, over, last bool) {
		t.Helper()
		r := mustRead(t, s, g, 1<<20, grade3, olympiad)
		if m := r.Modules[0]; m.TrialOver != over || m.LastInTrial != last {
			t.Errorf("%s: the series over %v and the last evidence in it %v, want %v and %v", what, m.TrialOver, m.LastInTrial, over, last)
		}
	}
	var s lm.State
	r := mustRead(t, s, g, 0, grade3, olympiad)
	if m := r.Modules[0]; m.HasEvidence || m.StartBand != ladderBand(2) || m.TrialOver || m.LastInTrial || len(m.Skills) != 0 {
		t.Errorf("a module of the engine before its first evidence: %+v", m)
	}
	for i := range 4 {
		s = mustAppend(t, s, choiceAnswer(topic("a"), 2, 3, i%2 == 0, ladderBand(2), int64(i)), g)
		flags("in the series", s, false, true)
	}
	s = mustAppend(t, s, choiceAnswer(topic("a"), 2, 3, false, ladderBand(2), 4), g)
	flags("at its last evidence", s, true, true)

	ai := choiceAnswer(topic("a"), 2, 3, false, ladderBand(2), 5)
	ai.Task.SkillTagsProvenance = lm.ProvenanceAI
	flags("after evidence that weighs nothing", mustAppend(t, s, ai, g), true, true)

	s = mustAppend(t, s, choiceAnswer(topic("a"), 2, 3, true, ladderBand(2), 6), g)
	flags("after it", s, true, false)
}

// mastered is the skill of the example of docs/architecture.md, section 3.2: mastered at the band
// `3-4`, θ + δ = 4.5, n = 100, m = 20 and τ = 10, practised at t₀.
func mastered(t0 int64) lm.State {
	return lm.State{Kind: lm.StateKind, Modules: []lm.ModuleState{{Module: olympiad, Level: 4, AnswerWeight: 100,
		StartBand: ladderBand(2), Skills: []lm.SkillState{{Skill: topic("a"), Correction: 0.5, AnswerWeight: 20,
			EvidenceCount: 30, LastPractice: t0, Stability: 10, Mastered: true, MasteryBand: 2,
			Traps: []lm.TrapCount{{ID: olympiad + "/trap/order", Count: 2, Last: t0}}}}}}}
}

// candidate is a task of the band `3-4` of difficulty 3 with five options, β = 2.5.
func candidate(skills ...string) lm.Task {
	return lm.Task{ID: olympiad + "/task/c", Band: ladderBand(2), Difficulty: 3, Primary: skills,
		Items: []lm.TaskItem{{Place: 1, Slots: []lm.Slot{{Check: lm.CheckKindChoice, Options: 5, HasOptions: true}}}}}
}

func taskP(t *testing.T, r lm.Readout, task lm.Task) []lm.ItemProbability {
	t.Helper()
	out, err := r.Task(task)
	if err != nil {
		t.Fatal(err)
	}
	return out.Items
}

func near3(t *testing.T, what string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-3 {
		t.Errorf("%s = %.6f, want %.3f", what, got, want)
	}
}

// The example of docs/architecture.md, section 3.2, read: at its last practice and after the
// summer break of 60 days.
func TestReadSkill(t *testing.T) {
	g := olympiadGraphs(t, map[string][]int32{topic("a"): {1, 2, 3}})
	const t0 = 1762041600
	s := mastered(t0)

	r := mustRead(t, s, g, t0, lm.Grade{}, olympiad)
	sk := r.Modules[0].Skills[0]
	if sk.Skill != topic("a") || sk.EvidenceCount != 30 || sk.LastPractice != t0 || sk.Estimate != 4.5 ||
		sk.Retention != 1 || !sk.Mastered || sk.MasteryBand != ladderBand(2) || sk.Review.Already ||
		len(sk.Traps) != 1 || sk.Traps[0].Count != 2 {
		t.Errorf("the skill at its last practice: %+v", sk)
	}
	near3(t, "v", sk.Uncertainty, 0.2327)
	near3(t, "q", sk.CautiousProbability, 0.820)
	if days := float64(float64(sk.Review.At-t0) / day); math.Abs(days-12.9) > 0.05 {
		t.Errorf("the review after τ = 10 in %.3f days, want 12.9", days)
	}
	near3(t, "P without forgetting", taskP(t, r, candidate(topic("a")))[0].Probability, 0.905)

	r = mustRead(t, s, g, t0+60*day, lm.Grade{}, olympiad)
	sk = r.Modules[0].Skills[0]
	near3(t, "R(60; 10)", sk.Retention, 0.645)
	near3(t, "R · q", sk.CautiousProbability, 0.645*0.820)
	if !r.Modules[0].TrialOver || r.Modules[0].LastInTrial || !r.Modules[0].HasEvidence || len(r.Modules[0].Traps) != 0 {
		t.Errorf("the module past its series: %+v", r.Modules[0])
	}
	// A trap with no skill that has a state is the module's, and reading gives it there.
	withTrap := mastered(t0)
	withTrap.Modules[0].Traps = []lm.TrapCount{{ID: olympiad + "/trap/skip", Count: 3, Last: t0}}
	if m := mustRead(t, withTrap, g, t0, lm.Grade{}).Modules[0]; len(m.Traps) != 1 || m.Traps[0] != withTrap.Modules[0].Traps[0] {
		t.Errorf("the traps of the module: %+v", m.Traps)
	}
	near3(t, "P with forgetting", taskP(t, r, candidate(topic("a")))[0].Probability, 0.654)

	// A skill that dropped below the threshold even without forgetting is due already.
	low := mastered(t0)
	low.Modules[0].Level = 2
	if sk := mustRead(t, low, g, t0, lm.Grade{}).Modules[0].Skills[0]; !sk.Review.Already {
		t.Errorf("a review due already: %+v", sk.Review)
	}
	// A skill not mastered has no mastery fields.
	plain := mastered(t0)
	plain.Modules[0].Skills[0].Mastered, plain.Modules[0].Skills[0].MasteryBand = false, 0
	if sk := mustRead(t, plain, g, t0, lm.Grade{}).Modules[0].Skills[0]; sk.Mastered ||
		sk.MasteryBand != (lm.BandRef{}) || sk.CautiousProbability != 0 || sk.Review != (lm.ReviewDate{}) {
		t.Errorf("a skill not mastered: %+v", sk)
	}
}

// Readout.Task: the starting point of a module with no state, R̄ = 1 in the trial series, items
// without a verdict, and the transfer.
func TestReadTask(t *testing.T) {
	g := olympiadGraphs(t, map[string][]int32{topic("a"): {1, 2, 3}})

	// A learner of grade 3 starts at θ₀ = 2.5: P = 0.2 + 0.8 · σ(0) = 0.6 (section 3.1).
	r := mustRead(t, lm.State{}, g, 0, lm.Grade{Value: 3, Known: true}, olympiad)
	task := candidate(topic("a"))
	task.Items = append(task.Items,
		lm.TaskItem{Place: 3, Slots: []lm.Slot{{Check: 9}}}, // a kind of check of a later version
		lm.TaskItem{Place: 2})                               // no slot: a rubric
	items := taskP(t, r, task)
	if len(items) != 3 || items[0].Place != 1 || !items[0].HasVerdict ||
		items[1] != (lm.ItemProbability{Place: 2}) || items[2] != (lm.ItemProbability{Place: 3}) {
		t.Errorf("the items of a candidate: %+v", items)
	}
	near3(t, "P at the starting point", items[0].Probability, 0.6)

	// A time of format 0.4 has a verdict and no guessing: P = σ(0) = 0.5 at the starting point.
	times := candidate(topic("a"))
	times.Items[0].Slots = []lm.Slot{{Check: lm.CheckKindTime}}
	if got := taskP(t, r, times); len(got) != 1 || !got[0].HasVerdict {
		t.Errorf("an item of a time: %+v", got)
	} else {
		near3(t, "P of a time at the starting point", got[0].Probability, 0.5)
	}

	// In the trial series P is without retention, however long ago the skill was practised; past
	// it, retention lowers P.
	trial := mastered(0)
	trial.Modules[0].AnswerWeight, trial.Modules[0].LastInTrial = 3, true
	soon, late := taskP(t, mustRead(t, trial, g, 0, lm.Grade{}), task), taskP(t, mustRead(t, trial, g, 300*day, lm.Grade{}), task)
	if soon[0] != late[0] {
		t.Errorf("a pause in the trial series changed P: %v and %v", soon[0], late[0])
	}
	past := mastered(0)
	if p0, p1 := taskP(t, mustRead(t, past, g, 0, lm.Grade{}), task), taskP(t, mustRead(t, past, g, 300*day, lm.Grade{}), task); !(p1[0].Probability < p0[0].Probability) {
		t.Errorf("a pause past the trial series left P at %v and %v", p0[0], p1[0])
	}

	// A retired skill reads as its replacement.
	og, err := lm.NewGraph(olympiad, ladder(), []lm.Skill{{ID: topic("a"), Assessability: lm.AssessabilityVerdict}},
		[]lm.RetiredSkill{{ID: topic("old"), ReplacedBy: []string{topic("a")}}},
		[]lm.Task{{ID: olympiad + "/task/x", Band: ladderBand(2), Difficulty: 3, Primary: []string{topic("a")}}})
	if err != nil {
		t.Fatal(err)
	}
	rg := mustRead(t, past, lm.Graphs{olympiad: {Graph: og, Current: true}}, 30*day, lm.Grade{})
	if a, b := taskP(t, rg, candidate(topic("a"))), taskP(t, rg, candidate(topic("old"))); a[0] != b[0] {
		t.Errorf("P by the retired skill %v, by its replacement %v", b[0], a[0])
	}
}

func TestReadRefuses(t *testing.T) {
	g := anchorGraphs(t)
	unordered := anchorGraphs(t)
	sg := unordered[school]
	sg.Graph.Skills = slices.Clone(sg.Graph.Skills)
	slices.Reverse(sg.Graph.Skills)
	unordered[school] = sg
	trialWithout := schoolState(false)
	trialWithout.Modules[0].AnswerWeight = 3

	cases := []struct {
		what   string
		s      lm.State
		g      lm.Graphs
		engine []string
	}{
		{"a module of the engine twice", lm.State{}, g, []string{school, school}},
		{"a module of the engine by no module ID", lm.State{}, g, []string{"example"}},
		{"a module of the engine by an ID of a skill", lm.State{}, g, []string{schoolSkill("division")}},
		{"a module of the engine with no entry", lm.State{}, g, []string{"example/none"}},
		{"a module of the state with no entry", schoolState(true), lm.Graphs{olympiad: g[olympiad]}, nil},
		{"a graph out of order", lm.State{}, unordered, []string{school}},
		{"a module in the series whose last evidence was not in it", trialWithout, g, nil},
		{"a module whose band is out of its place", lm.State{}, lm.Graphs{"example/empty": {Graph: lm.Graph{Module: "example/empty",
			Bands: []lm.Band{{Place: 2, ID: "band-2"}}}, Current: true}}, []string{"example/empty"}},
		{"a skill taught at a band the graph has not", lm.State{}, lm.Graphs{"example/empty": {Graph: lm.Graph{Module: "example/empty",
			Skills: []lm.Skill{{ID: "example/empty/skill/a", Assessability: lm.AssessabilityVerdict, BandsTaught: []int32{1}}}},
			Current: true}}, []string{"example/empty"}},
	}
	for _, c := range cases {
		if _, err := lm.Read(c.s, c.g, 0, lm.Grade{}, c.engine); !errors.Is(err, lm.ErrInvalid) {
			t.Errorf("%s: %v", c.what, err)
		}
	}

	r := mustRead(t, lm.State{}, g, 0, lm.Grade{}, olympiad)
	choice1 := candidate(topic("a"))
	choice1.Items[0].Slots[0].Options = 1
	hardest := candidate(topic("a"))
	hardest.Difficulty = 6
	lowest := candidate(topic("a"))
	lowest.Band.Place = 0
	twice := candidate(topic("a"))
	twice.Items = append(twice.Items, twice.Items[0])
	tasks := []struct {
		what string
		r    lm.Readout
		task lm.Task
	}{
		{"a choice of one option", r, choice1},
		{"a difficulty of 6", r, hardest},
		{"a band place of 0", r, lowest},
		{"an item place twice", r, twice},
		{"an ID that names no module", r, lm.Task{ID: "example", Difficulty: 3, Band: ladderBand(1)}},
		{"a module with no entry", r, lm.Task{ID: "example/none/task/x", Difficulty: 3, Band: ladderBand(1)}},
		{"a module with no current version", mustRead(t, lm.State{}, lm.Graphs{olympiad: {}}, 0, lm.Grade{}), candidate(topic("a"))},
		{"the zero Readout", lm.Readout{}, candidate(topic("a"))},
		{"a graph out of order, of a module Read did not read", mustRead(t, lm.State{}, unordered, 0, lm.Grade{}),
			lm.Task{ID: school + "/task/d", Difficulty: 3, Band: lm.BandRef{Place: 1, ID: "grade-3"}}},
		{"a band out of its place, of a module Read did not read", mustRead(t, lm.State{}, lm.Graphs{"example/odd": {Graph: lm.Graph{
			Module: "example/odd", Bands: []lm.Band{{Place: 2, ID: "band-2"}}}, Current: true}}, 0, lm.Grade{}),
			lm.Task{ID: "example/odd/task/x", Difficulty: 3, Band: lm.BandRef{Place: 1}}},
	}
	for _, c := range tasks {
		if _, err := c.r.Task(c.task); !errors.Is(err, lm.ErrInvalid) {
			t.Errorf("%s: %v", c.what, err)
		}
	}
	empty := mustRead(t, lm.State{}, lm.Graphs{"example/empty": {Graph: lm.Graph{Module: "example/empty"}, Current: true}}, 0, lm.Grade{})
	if _, err := empty.Task(lm.Task{ID: "example/empty/task/x", Difficulty: 3, Band: lm.BandRef{Place: 1}}); !errors.Is(err, lm.ErrInvalid) {
		t.Errorf("a candidate of a module with no bands: %v", err)
	}
}

// Read and Task on the random logs of the property: they read every state the fold gives, and
// every number they give is a probability or a time it may be.
func TestReadRandom(t *testing.T) {
	for seed := range uint64(30) {
		r := rand.New(rand.NewPCG(seed, 20))
		w := newWorld(t, r)
		var s lm.State
		var engine []string
		for _, m := range w.modules {
			if r.IntN(2) == 0 && w.graphs[m.id].Current {
				engine = append(engine, m.id)
			}
		}
		at := int64(1762041600)
		var tasks []lm.Task
		for n := range 1 + r.IntN(40) {
			at += int64(r.IntN(200000))
			e := journalEvidence(t, r, w, n, at)
			s = mustAppend(t, s, e, w.graphs)
			if w.graphs[e.Task.ID[:len("example/m0")]].Current {
				tasks = append(tasks, e.Task)
			}
		}
		rd, err := lm.Read(s, w.graphs, at+int64(r.IntN(100*day)), lm.Grade{Value: 3, Known: r.IntN(2) == 0}, engine)
		if err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}
		for _, m := range rd.Modules {
			for _, sk := range m.Skills {
				if !(sk.Retention > 0 && sk.Retention <= 1) || !(sk.Uncertainty > 0) ||
					(sk.Mastered && !(sk.CautiousProbability >= 0 && sk.CautiousProbability <= 1)) {
					t.Fatalf("seed %d: %+v", seed, sk)
				}
			}
		}
		for _, task := range tasks {
			out, err := rd.Task(task)
			if err != nil {
				t.Fatalf("seed %d, task %s: %v", seed, task.ID, err)
			}
			for _, p := range out.Items {
				if p.HasVerdict && !(p.Probability > 0 && p.Probability <= 1) {
					t.Fatalf("seed %d, task %s: %+v", seed, task.ID, p)
				}
			}
		}
	}
}
