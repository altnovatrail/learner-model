package learnermodel_test

import (
	"bytes"
	"errors"
	"math"
	"math/rand/v2"
	"strings"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

// The fold of the log (docs/architecture.md, section 3.5): appending a piece of evidence gives,
// byte for byte, the state recomputed from the log, and the rules of the step around it. Every
// value here is synthetic.

// world is a set of modules with their graphs, for the random logs of the property.
type world struct {
	graphs  lm.Graphs
	modules []worldModule
}

type worldModule struct {
	id     string
	bands  []lm.Band
	skills []string // the IDs evidence may name: skills with a verdict, one with a rubric, one unknown
}

func newWorld(t *testing.T, r *rand.Rand) world {
	t.Helper()
	w := world{graphs: lm.Graphs{}}
	for m := range 1 + r.IntN(3) {
		id := "example/m" + itoa(m)
		var bands []lm.Band
		withGrades := r.IntN(4) > 0
		for p := range 1 + r.IntN(4) {
			b := lm.Band{Place: int32(p + 1), ID: "band-" + itoa(p+1)}
			if withGrades {
				b.Grades = []int32{int32(2*p + 1), int32(2*p + 2)}
			}
			bands = append(bands, b)
		}
		var skills []lm.Skill
		var tasks []lm.Task
		mod := worldModule{id: id, bands: bands}
		for s := range 5 {
			sid := id + "/skill/s" + itoa(s)
			skills = append(skills, lm.Skill{ID: sid, Assessability: lm.AssessabilityVerdict})
			mod.skills = append(mod.skills, sid)
			for range 1 + r.IntN(2) {
				p := 1 + r.IntN(len(bands))
				tasks = append(tasks, lm.Task{ID: id + "/task/g", Band: lm.BandRef{Place: int32(p), ID: bands[p-1].ID},
					Difficulty: 3, Primary: []string{sid}})
			}
		}
		skills = append(skills, lm.Skill{ID: id + "/skill/essay", Assessability: lm.AssessabilityRubric})
		mod.skills = append(mod.skills, id+"/skill/essay", id+"/skill/unknown")
		g, err := lm.NewGraph(id, bands, skills, nil, tasks)
		if err != nil {
			t.Fatal(err)
		}
		w.graphs[id] = lm.ModuleGraph{Graph: g, Current: m != 2 || r.IntN(2) == 0}
		w.modules = append(w.modules, mod)
	}
	return w
}

func pick(r *rand.Rand, ids []string, n int) []string {
	var out []string
	for _, i := range r.Perm(len(ids))[:n] {
		out = append(out, ids[i])
	}
	return out
}

func randomSlot(r *rand.Rand) lm.Slot {
	switch r.IntN(8) {
	case 0:
		return lm.Slot{Check: lm.CheckKindFraction}
	case 1:
		return lm.Slot{Check: lm.CheckKindPredicate}
	case 2, 3, 4:
		return lm.Slot{Check: lm.CheckKindChoice, Options: int32(2 + r.IntN(4)), HasOptions: true}
	case 5:
		if r.IntN(6) == 0 {
			return lm.Slot{Check: 9} // a kind of a later version: its item gives no answer
		}
	}
	return lm.Slot{Check: lm.CheckKindNumber}
}

// journalEvidence is evidence of a random module of the world, with its start band by a random
// grade, so that two pieces of evidence of one module may carry different bands.
func journalEvidence(t *testing.T, r *rand.Rand, w world, n int, at int64) lm.Evidence {
	t.Helper()
	m := w.modules[r.IntN(len(w.modules))]
	grade := lm.Grade{Value: int32(r.IntN(9)), Known: r.IntN(5) > 0}
	start, err := lm.BandByGrade(grade, m.bands)
	if err != nil {
		t.Fatal(err)
	}
	p := 1 + r.IntN(len(m.bands))
	task := lm.Task{
		ID: m.id + "/task/t" + itoa(n), Band: lm.BandRef{Place: int32(p), ID: m.bands[p-1].ID},
		Difficulty: int32(1 + r.IntN(5)), Primary: pick(r, m.skills, 1+r.IntN(2)),
		KeyProvenance: lm.ProvenanceAuthor, SkillTagsProvenance: lm.ProvenanceAuthor,
	}
	items := 1
	if r.IntN(3) == 0 {
		items = 2 + r.IntN(19)
	}
	var answers []lm.ItemAnswer
	for i := range items {
		item := lm.TaskItem{Place: int32(i + 1)}
		if r.IntN(5) == 0 {
			item.OwnSkills, item.Primary = true, pick(r, m.skills, 1+r.IntN(2))
		}
		for range 1 + r.IntN(2) {
			item.Slots = append(item.Slots, randomSlot(r))
		}
		task.Items = append(task.Items, item)
		a := lm.ItemAnswer{Place: item.Place, Submitted: true}
		switch v := r.IntN(20); {
		case v < 9:
			a.FirstVerdict, a.HasFirstVerdict = lm.VerdictCorrect, true
		case v < 16:
			a.FirstVerdict, a.HasFirstVerdict = lm.VerdictNotYet, true
		case v < 19:
			a.Submitted = false
		default:
			a.FirstVerdict, a.HasFirstVerdict = 7, true // a verdict of a later version
		}
		answers = append(answers, a)
	}
	// Items out of their order: the step reads them by place.
	r.Shuffle(len(task.Items), func(i, j int) { task.Items[i], task.Items[j] = task.Items[j], task.Items[i] })
	return lm.Evidence{Kind: lm.EvidenceKind, Task: task, StartBand: start, PracticeTime: at,
		Outcome: lm.Outcome{Kind: lm.OutcomeFirstSubmission}, Items: answers}
}

func encodeState(t *testing.T, s lm.State) []byte {
	t.Helper()
	b, err := lm.EncodeState(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// The property of R80 of spec: for random synthetic logs and graphs, appending one piece of
// evidence at a time, with the state encoded and decoded between the steps as api keeps it, gives
// the bytes Recompute gives — of the evidence Append gave back, which api stores, and of the
// evidence as it came, since the step takes the start band of a module with a state from it.
func TestAppendIsRecompute(t *testing.T) {
	for seed := range uint64(50) {
		r := rand.New(rand.NewPCG(seed, 16))
		w := newWorld(t, r)
		var s lm.State
		var given, kept []lm.Evidence
		at := int64(1762041600)
		for n := range 1 + r.IntN(50) {
			at += int64(r.IntN(200000)) - 20000 // now and then a time that comes back
			e := journalEvidence(t, r, w, n, at)
			next, stored, err := lm.Append(s, e, w.graphs)
			if err != nil {
				t.Fatalf("seed %d, evidence %d: %v", seed, n, err)
			}
			if s, err = lm.DecodeState(encodeState(t, next)); err != nil {
				t.Fatalf("seed %d, evidence %d: %v", seed, n, err)
			}
			given, kept = append(given, e), append(kept, stored)
		}
		appended := encodeState(t, s)
		for _, log := range [][]lm.Evidence{kept, given} {
			again, err := lm.Recompute(log, w.graphs)
			if err != nil {
				t.Fatalf("seed %d: %v", seed, err)
			}
			if b := encodeState(t, again); !bytes.Equal(b, appended) {
				t.Fatalf("seed %d: the appends give\n%x\nand Recompute\n%x", seed, appended, b)
			}
		}
	}
}

const workbook = "example/workbook"

func workbookGraphs(t *testing.T) lm.Graphs {
	t.Helper()
	bands := []lm.Band{{Place: 1, ID: "grade-2", Grades: []int32{2}}, {Place: 2, ID: "grade-3", Grades: []int32{3}}}
	skills := []lm.Skill{
		{ID: skill("addition"), Assessability: lm.AssessabilityVerdict},
		{ID: skill("subtraction"), Assessability: lm.AssessabilityVerdict},
		{ID: skill("essay"), Assessability: lm.AssessabilityRubric},
	}
	tasks := []lm.Task{{ID: workbook + "/task/g", Band: lm.BandRef{Place: 1, ID: "grade-2"}, Difficulty: 3,
		Primary: []string{skill("addition"), skill("subtraction")}}}
	g, err := lm.NewGraph(workbook, bands, skills, nil, tasks)
	if err != nil {
		t.Fatal(err)
	}
	return lm.Graphs{workbook: {Graph: g, Current: true}}
}

// page is evidence of a page of the workbook: items of one number slot each, the given number of
// them right and the rest not yet, all with the given primary skills.
func page(items, right int, primary []string, at int64) lm.Evidence {
	e := lm.Evidence{
		Kind: lm.EvidenceKind,
		Task: lm.Task{ID: workbook + "/task/page", Band: lm.BandRef{Place: 1, ID: "grade-2"}, Difficulty: 1, Primary: primary,
			KeyProvenance: lm.ProvenanceAuthor, SkillTagsProvenance: lm.ProvenanceAuthor},
		StartBand: lm.BandRef{Place: 1, ID: "grade-2"}, Outcome: lm.Outcome{Kind: lm.OutcomeFirstSubmission}, PracticeTime: at,
	}
	for i := range items {
		e.Task.Items = append(e.Task.Items, lm.TaskItem{Place: int32(i + 1), Slots: []lm.Slot{{Check: lm.CheckKindNumber}}})
		v := lm.VerdictCorrect
		if i >= right {
			v = lm.VerdictNotYet
		}
		e.Items = append(e.Items, lm.ItemAnswer{Place: int32(i + 1), Submitted: true, FirstVerdict: v, HasFirstVerdict: true})
	}
	return e
}

func mustAppend(t *testing.T, s lm.State, e lm.Evidence, g lm.Graphs) lm.State {
	t.Helper()
	s, _, err := lm.Append(s, e, g)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// The thresholds of the sums of weights (docs/architecture.md, section 3.5): pages of 8 and 13
// items weigh 2.4 and 2.6, and the two make n = 5 and m = 5 within ε, whatever the float sum gives.
// The pages are of difficulty 5, so that answering them all right masters the skill.
func TestWeightThresholds(t *testing.T) {
	g := workbookGraphs(t)
	add := []string{skill("addition")}
	s := mustAppend(t, lm.State{}, hard(page(8, 8, add, 1)), g)
	if m := s.Modules[0]; len(m.Trial) != 8 || math.Abs(m.AnswerWeight-2.4) > 1e-12 {
		t.Fatalf("after the first page: %d trial answers, n %v", len(m.Trial), m.AnswerWeight)
	}
	s = mustAppend(t, s, hard(page(13, 13, add, 2)), g)
	m := s.Modules[0]
	sk := skillOf(t, s, add[0])
	if len(m.Trial) != 0 || math.Abs(m.AnswerWeight-5) > 1e-9 || math.Abs(sk.AnswerWeight-5) > 1e-9 {
		t.Fatalf("after the second page: %d trial answers, n %v, m %v; want the series over at 5", len(m.Trial), m.AnswerWeight, sk.AnswerWeight)
	}
	if !sk.Mastered || sk.MasteryBand != 1 || sk.Correction != 0 {
		t.Errorf("after the second page: mastered %v at %d, δ %v; want mastered at band 1, δ still 0", sk.Mastered, sk.MasteryBand, sk.Correction)
	}
	// The series is over, so the next page moves δ by the step of Elo.
	s = mustAppend(t, s, hard(page(4, 1, add, 3)), g)
	sk = skillOf(t, s, add[0])
	if sk.Correction >= 0 || sk.FailuresInARow != 1 || !sk.Mastered {
		t.Errorf("a failed page: δ %v, failures %d, mastered %v", sk.Correction, sk.FailuresInARow, sk.Mastered)
	}
}

// One step per piece of evidence (docs/architecture.md, section 3.3): its answers weigh u = 3 /
// (N + 2) each, a skill's by its own N_s, and a primary skill of an item with two gets half of the
// item's move; skills that are not skills of the graph with a verdict get nothing, and count in
// neither J nor δ̄.
func TestWeightsAndShares(t *testing.T) {
	g := workbookGraphs(t)
	add, sub := skill("addition"), skill("subtraction")
	before := lm.State{Kind: lm.StateKind, Modules: []lm.ModuleState{{
		Module: workbook, Level: 0.5, AnswerWeight: 20, StartBand: lm.BandRef{Place: 1, ID: "grade-2"},
		Skills: []lm.SkillState{{Skill: add, Correction: 0.4, AnswerWeight: 8, EvidenceCount: 3}},
	}}}
	e := page(3, 2, []string{add, sub, skill("essay"), skill("unknown")}, 7)
	// The third item names its own skill, one with no state to move.
	e.Task.Items[2].OwnSkills, e.Task.Items[2].Primary = true, []string{skill("essay")}
	s := mustAppend(t, before, e, g)

	// β = −2: band 1, difficulty 1; δ̄ = (0.4 + 0)/2 on the first two items, 0 on the third.
	sigma := func(x float64) float64 { return float64(1 / (1 + math.Exp(-x))) }
	p2, p1 := sigma(0.5+0.2+2), sigma(0.5+2)
	u := float64(3.0 / 5)
	theta := 0.5 + float64(0.1*u*(1-p2+1-p2+0-p1))
	us := float64(3.0 / 4)
	pull := float64((1-p2)/2) + float64((1-p2)/2)
	wantAdd := 0.4 + float64(float64(0.4/(1+float64(0.05*8)))*us*pull)
	wantSub := float64(float64(0.4*us) * pull)
	near(t, "θ", s.Modules[0].Level, theta)
	if math.Abs(s.Modules[0].AnswerWeight-21.8) > 1e-12 {
		t.Errorf("n = %v, want 21.8", s.Modules[0].AnswerWeight)
	}
	near(t, "δ of addition", skillOf(t, s, add).Correction, wantAdd)
	near(t, "δ of subtraction", skillOf(t, s, sub).Correction, wantSub)
	for _, id := range []string{add, sub} {
		sk := skillOf(t, s, id)
		if math.Abs(sk.AnswerWeight-(map[string]float64{add: 8, sub: 0}[id]+1.5)) > 1e-12 || sk.LastPractice != 7 {
			t.Errorf("%s: m %v, last practice %d", id, sk.AnswerWeight, sk.LastPractice)
		}
	}
	if len(s.Modules[0].Skills) != 2 {
		t.Errorf("skills with no state got one: %+v", s.Modules[0].Skills)
	}
}

func validEvidence() lm.Evidence {
	return choiceAnswer(topic("a"), 1, 3, true, ladderBand(1), 5)
}

func TestAppendRefuses(t *testing.T) {
	g := olympiadGraphs(t, map[string][]int32{topic("a"): {1}})
	evidence := []struct {
		name   string
		change func(*lm.Evidence)
		graphs lm.Graphs
	}{
		{"a kind below 1", func(e *lm.Evidence) { e.Kind = 0 }, g},
		{"a primary skill twice", func(e *lm.Evidence) { e.Task.Primary = []string{topic("a"), topic("a")} }, g},
		{"a task that names no module", func(e *lm.Evidence) { e.Task.ID = "example/olympiad/task" }, g},
		{"a task ID with an empty part", func(e *lm.Evidence) { e.Task.ID = "example//task/x" }, g},
		{"difficulty 6", func(e *lm.Evidence) { e.Task.Difficulty = 6 }, g},
		{"a task band at place 0", func(e *lm.Evidence) { e.Task.Band.Place = 0 }, g},
		{"a start band at place 0", func(e *lm.Evidence) { e.StartBand.Place = 0 }, g},
		{"a choice of one option", func(e *lm.Evidence) { e.Task.Items[0].Slots[0].Options = 1 }, g},
		{"a choice without its options", func(e *lm.Evidence) { e.Task.Items[0].Slots[0] = lm.Slot{Check: lm.CheckKindChoice} }, g},
		{"an answer with no slot", func(e *lm.Evidence) { e.Task.Items[0].Slots = nil }, g},
		{"a module with no entry", func(*lm.Evidence) {}, lm.Graphs{}},
		{"the graph of another module", func(*lm.Evidence) {}, lm.Graphs{olympiad: {Graph: lm.Graph{Module: workbook}, Current: true}}},
		{"skills out of their order", func(*lm.Evidence) {}, lm.Graphs{olympiad: {Current: true, Graph: lm.Graph{Module: olympiad,
			Skills: []lm.Skill{{ID: topic("b")}, {ID: topic("a")}}}}}},
		{"bands taught out of their order", func(*lm.Evidence) {}, lm.Graphs{olympiad: {Current: true, Graph: lm.Graph{Module: olympiad,
			Skills: []lm.Skill{{ID: topic("a"), BandsTaught: []int32{2, 1}}}}}}},
	}
	for _, c := range evidence {
		e := validEvidence()
		c.change(&e)
		if _, _, err := lm.Append(lm.State{}, e, c.graphs); !errors.Is(err, lm.ErrInvalid) {
			t.Errorf("%s: %v, want ErrInvalid", c.name, err)
		}
	}

	good := mustAppend(t, lm.State{}, validEvidence(), g)
	states := []struct {
		name   string
		change func(*lm.State)
		want   error
	}{
		{"a kind above StateKind", func(s *lm.State) { s.Kind = lm.StateKind + 1 }, lm.ErrInvalid},
		{"kind 0 with modules", func(s *lm.State) { s.Kind = 0 }, lm.ErrInvalid},
		{"modules out of their order", func(s *lm.State) {
			s.Modules = append(s.Modules, lm.ModuleState{Module: "example/a", StartBand: ladderBand(1)})
		}, lm.ErrInvalid},
		{"a start band at place 0", func(s *lm.State) { s.Modules[0].StartBand.Place = 0 }, lm.ErrInvalid},
		{"a trial answer of difficulty 0", func(s *lm.State) { s.Modules[0].Trial[0].Difficulty = 0 }, lm.ErrInvalid},
		{"a trial choice of no options", func(s *lm.State) { s.Modules[0].Trial[0].Slots[0].Options = 0 }, lm.ErrInvalid},
		{"a trial answer with no slot", func(s *lm.State) { s.Modules[0].Trial[0].Slots = nil }, lm.ErrInvalid},
		{"a trial answer of an unknown check", func(s *lm.State) { s.Modules[0].Trial[0].Slots[0] = lm.Slot{Check: 9} }, lm.ErrInvalid},
		{"a trial answer with no verdict", func(s *lm.State) { s.Modules[0].Trial[0].Verdict = 7 }, lm.ErrInvalid},
		{"a trial answer of weight 0", func(s *lm.State) { s.Modules[0].Trial[0].Weight = 0 }, lm.ErrInvalid},
		{"trial answers past the series", func(s *lm.State) { s.Modules[0].AnswerWeight = 5 }, lm.ErrInvalid},
		{"a level that is NaN", func(s *lm.State) { s.Modules[0].Level = math.Float64frombits(0x7ff8000000000001) }, lm.ErrNotFinite},
	}
	for _, c := range states {
		s, err := lm.DecodeState(encodeState(t, good))
		if err != nil {
			t.Fatal(err)
		}
		c.change(&s)
		if _, _, err := lm.Append(s, validEvidence(), g); !errors.Is(err, c.want) {
			t.Errorf("%s: %v, want %v", c.name, err, c.want)
		}
	}

	bad := validEvidence()
	bad.Task.Difficulty = 0
	_, err := lm.Recompute([]lm.Evidence{validEvidence(), bad}, g)
	if !errors.Is(err, lm.ErrInvalid) || !strings.Contains(err.Error(), "at index 1 of the log") {
		t.Errorf("Recompute: %v, want ErrInvalid at index 1", err)
	}
}

// What moves nothing, and the start band the step takes.
func TestAppendRules(t *testing.T) {
	a := topic("a")
	g := olympiadGraphs(t, map[string][]int32{a: {1, 2}})
	waiting := lm.Graphs{olympiad: {Current: false}}

	// The start band: the evidence's for a module with no state, and the state's after it.
	s, e, err := lm.Append(lm.State{}, choiceAnswer(a, 1, 3, true, ladderBand(2), 1), g)
	if err != nil || s.Kind != lm.StateKind || s.Modules[0].StartBand != ladderBand(2) || e.StartBand != ladderBand(2) {
		t.Fatalf("the first evidence: %+v, %+v, %v", s, e.StartBand, err)
	}
	if s.Modules[0].Level == 2.5 {
		t.Errorf("the trial series did not move θ from θ₀")
	}
	s2, e, err := lm.Append(s, choiceAnswer(a, 1, 3, true, ladderBand(1), 2), g)
	if err != nil || e.StartBand != ladderBand(2) || s2.Modules[0].StartBand != ladderBand(2) {
		t.Errorf("a later grade moved the start band: %+v, %v", e.StartBand, err)
	}

	// A module with no current version: its evidence moves nothing and waits, with the band of the
	// state when the module has one.
	same, e, err := lm.Append(s, choiceAnswer(a, 1, 3, true, ladderBand(1), 3), waiting)
	if err != nil || !sameState(t, same, s) || e.StartBand != ladderBand(2) {
		t.Errorf("no current version moved the state: %v, band %+v", err, e.StartBand)
	}
	empty, e, err := lm.Append(lm.State{}, choiceAnswer(a, 1, 3, true, ladderBand(3), 3), waiting)
	if err != nil || len(empty.Modules) != 0 || e.StartBand != ladderBand(3) {
		t.Errorf("no current version gave a module a state: %+v, %v", empty, err)
	}

	// Evidence with no answer moves nothing, and gives a module no state.
	for _, change := range []func(*lm.ItemAnswer){
		func(a *lm.ItemAnswer) { a.HasFirstVerdict, a.FirstVerdict = false, 0 },
		func(a *lm.ItemAnswer) { a.FirstVerdict = 7 },
	} {
		none := choiceAnswer(a, 1, 3, true, ladderBand(1), 4)
		change(&none.Items[0])
		if got := mustAppend(t, s, none, g); !sameState(t, got, s) {
			t.Errorf("evidence with no answer moved the state")
		}
		if got := mustAppend(t, lm.State{}, none, g); len(got.Modules) != 0 {
			t.Errorf("evidence with no answer gave a module a state")
		}
	}
	unknownKind := choiceAnswer(a, 1, 3, true, ladderBand(1), 4)
	unknownKind.Task.Items[0].Slots[0] = lm.Slot{Check: 9}
	if got := mustAppend(t, s, unknownKind, g); !sameState(t, got, s) {
		t.Errorf("an item of an unknown check moved the state")
	}

	// Evidence of a later kind is read as far as this version knows it.
	later := choiceAnswer(a, 1, 3, true, ladderBand(1), 5)
	later.Kind = lm.EvidenceKind + 1
	if _, _, err := lm.Append(s, later, g); err != nil {
		t.Errorf("evidence of a later kind: %v", err)
	}

	// Append changes nothing it was given, and gives a state that shares no memory with it.
	b := encodeState(t, s)
	out := mustAppend(t, s, choiceAnswer(a, 1, 3, false, ladderBand(1), 6), g)
	out.Modules[0].Skills[0].Correction = 99
	out.Modules[0].Trial[0].Slots[0].Options = 99
	if !bytes.Equal(encodeState(t, s), b) {
		t.Errorf("Append changed the state it was given")
	}
}

// The mastery is held at its band and raised only above it: a right answer at a lower band does
// not move it down, and one that clears a higher band raises it.
func TestMasteryHeldAndRaised(t *testing.T) {
	a := topic("a")
	g := olympiadGraphs(t, map[string][]int32{a: {1, 2, 3}})
	s := lm.State{Kind: lm.StateKind, Modules: []lm.ModuleState{{
		Module: olympiad, Level: 4, AnswerWeight: 200, StartBand: ladderBand(1),
		Skills: []lm.SkillState{{Skill: a, Correction: 1, AnswerWeight: 50, EvidenceCount: 50, Mastered: true, MasteryBand: 2}},
	}}}
	// A right answer at band 1: the band held, 2, stays.
	s = mustAppend(t, s, choiceAnswer(a, 1, 3, true, ladderBand(1), 1), g)
	if sk := skillOf(t, s, a); !sk.Mastered || sk.MasteryBand != 2 {
		t.Errorf("a right answer below the band held: %+v", sk)
	}
	// A right answer at band 3, which λ ≈ 4.6 clears for band 2 but not for band 3.
	s = mustAppend(t, s, choiceAnswer(a, 3, 3, true, ladderBand(1), 2), g)
	if sk := skillOf(t, s, a); sk.MasteryBand != 2 {
		t.Errorf("mastered at %d, want still 2", sk.MasteryBand)
	}
	s.Modules[0].Level = 7
	s = mustAppend(t, s, choiceAnswer(a, 3, 3, true, ladderBand(1), 3), g)
	if sk := skillOf(t, s, a); sk.MasteryBand != 3 {
		t.Errorf("mastered at %d, want raised to 3", sk.MasteryBand)
	}
}

// hard makes a page of the workbook one of difficulty 5.
func hard(e lm.Evidence) lm.Evidence {
	e.Task.Difficulty = 5
	return e
}
