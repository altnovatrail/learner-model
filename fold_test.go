package learnermodel_test

import (
	"bytes"
	"errors"
	"fmt"
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
	skills []string // the IDs evidence may name: skills with a verdict, one with a rubric, one unknown, retired ones
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
		// Skills of earlier versions: renamed, split, merged, retired with `[]` and replaced by a
		// skill with a rubric.
		s := func(name string) string { return id + "/skill/" + name }
		retired := []lm.RetiredSkill{
			{ID: s("r0"), ReplacedBy: []string{s("s0")}},
			{ID: s("r1"), ReplacedBy: []string{s("s1"), s("s2")}},
			{ID: s("r2"), ReplacedBy: []string{s("s3")}},
			{ID: s("r3"), ReplacedBy: []string{s("s3")}},
			{ID: s("r4")},
			{ID: s("r5"), ReplacedBy: []string{s("essay")}},
		}
		for _, r := range retired {
			mod.skills = append(mod.skills, r.ID)
		}
		g, err := lm.NewGraph(id, bands, skills, retired, tasks)
		if err != nil {
			t.Fatal(err)
		}
		w.graphs[id] = lm.ModuleGraph{Graph: g, Current: m != 2 || r.IntN(2) == 0}
		w.modules = append(w.modules, mod)
	}
	// A graph of a module no evidence names, as reading brings: it changes nothing.
	w.graphs["example/extra"] = lm.ModuleGraph{Graph: lm.Graph{Module: "example/extra"}, Current: true}
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
			a.WrongBeforeFirst = int32(r.IntN(3) / 2) // now and then a rung before it
		case v < 16:
			a.FirstVerdict, a.HasFirstVerdict = lm.VerdictNotYet, true
		case v < 18:
			a.Submitted = false
		case v < 19:
			reasons := []lm.UnverdictedReason{lm.UnverdictedUnreadableBlank, lm.UnverdictedUnreadableOther, lm.UnverdictedUnchecked}
			for range 1 + r.IntN(2) {
				a.Unverdicted = append(a.Unverdicted, reasons[r.IntN(len(reasons))])
			}
		default:
			a.FirstVerdict, a.HasFirstVerdict = 7, true // a verdict of a later version
		}
		answers = append(answers, a)
	}
	// Items out of their order: the step reads them by place.
	r.Shuffle(len(task.Items), func(i, j int) { task.Items[i], task.Items[j] = task.Items[j], task.Items[i] })
	e := lm.Evidence{Kind: lm.EvidenceKind, Task: task, StartBand: start, PracticeTime: at, Items: answers}
	outcomes := []lm.OutcomeKind{lm.OutcomeFirstSubmission, lm.OutcomeFirstSubmission, lm.OutcomeAfterHint,
		lm.OutcomeSolutionShown, lm.OutcomeSolutionShown, lm.OutcomeLeftUnfinished, lm.OutcomeUnchecked, lm.OutcomeWithSketch, 9}
	e.Outcome.Kind = outcomes[r.IntN(len(outcomes))]
	if e.Outcome.Kind == lm.OutcomeAfterHint {
		e.Outcome.Rung, e.Outcome.HasRung = int32(1+r.IntN(3)), true
	}
	e.TechniqueOnRequest = r.IntN(8) == 0
	e.TechniqueOnStep = r.IntN(4) == 0
	e.HasSolutionSeenFact = true
	e.SolutionSeenBefore = r.IntN(8) == 0
	if r.IntN(10) == 0 {
		provenances := []lm.Provenance{lm.ProvenanceUnspecified, lm.ProvenanceSolver, lm.ProvenanceAdult, lm.ProvenanceAI, 9}
		e.Task.KeyProvenance = provenances[r.IntN(len(provenances))]
		e.Task.SkillTagsProvenance = provenances[r.IntN(len(provenances))]
	}
	for i := range r.IntN(3) {
		trap := lm.FoundTrap{ID: m.id + "/trap/t" + itoa(i), ItemPlaces: []int32{int32(1 + r.IntN(items))}}
		trap.Skills = pick(r, m.skills, r.IntN(3))
		e.Traps = append(e.Traps, trap)
	}
	return e
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
// The logs name skills of earlier versions, which go to their replacements, and lose any of their
// evidence, as api cuts that of a version withdrawn for a wrong key.
func TestAppendIsRecompute(t *testing.T) {
	for seed := range uint64(50) {
		r := rand.New(rand.NewPCG(seed, 16))
		w := newWorld(t, r)
		var log []lm.Evidence
		at := int64(1762041600)
		for n := range 1 + r.IntN(50) {
			at += int64(r.IntN(200000)) - 20000 // now and then a time that comes back
			if e := journalEvidence(t, r, w, n, at); r.IntN(6) > 0 {
				log = append(log, e)
			}
		}
		var s lm.State
		var given, kept []lm.Evidence
		for n, e := range log {
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
		Skills: []lm.SkillState{{Skill: add, Correction: 0.4, AnswerWeight: 8, EvidenceCount: 3, LastPractice: 7, Stability: 1}},
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
		{"retired skills out of their order", func(*lm.Evidence) {}, lm.Graphs{olympiad: {Current: true, Graph: lm.Graph{Module: olympiad,
			Skills: []lm.Skill{{ID: topic("a")}}, Retired: []lm.RetiredSkill{{ID: topic("y")}, {ID: topic("x")}}}}}},
		{"a retired skill twice", func(*lm.Evidence) {}, lm.Graphs{olympiad: {Current: true, Graph: lm.Graph{Module: olympiad,
			Skills: []lm.Skill{{ID: topic("a")}}, Retired: []lm.RetiredSkill{{ID: topic("x")}, {ID: topic("x")}}}}}},
		{"replacements out of their order", func(*lm.Evidence) {}, lm.Graphs{olympiad: {Current: true, Graph: lm.Graph{Module: olympiad,
			Skills:  []lm.Skill{{ID: topic("a")}, {ID: topic("b")}},
			Retired: []lm.RetiredSkill{{ID: topic("x"), ReplacedBy: []string{topic("b"), topic("a")}}}}}}},
		{"a retired skill that is still a skill", func(*lm.Evidence) {}, lm.Graphs{olympiad: {Current: true, Graph: lm.Graph{Module: olympiad,
			Skills: []lm.Skill{{ID: topic("a")}}, Retired: []lm.RetiredSkill{{ID: topic("a")}}}}}},
		// A chain the graph has not made one step: x goes to y, which is retired too.
		{"a replacement that is no skill", func(*lm.Evidence) {}, lm.Graphs{olympiad: {Current: true, Graph: lm.Graph{Module: olympiad,
			Skills:  []lm.Skill{{ID: topic("a")}},
			Retired: []lm.RetiredSkill{{ID: topic("x"), ReplacedBy: []string{topic("y")}}, {ID: topic("y"), ReplacedBy: []string{topic("a")}}}}}}},
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
		{"a skill of a stability below τ₀", func(s *lm.State) { s.Modules[0].Skills[0].Stability = 0.5 }, lm.ErrInvalid},
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
		Skills: []lm.SkillState{{Skill: a, Correction: 1, AnswerWeight: 50, EvidenceCount: 50, LastPractice: 1, Stability: 10, Mastered: true, MasteryBand: 2}},
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

// The example page of docs/architecture.md, section 3.3: 10 items of a number slot with one
// primary skill, band 1, difficulty 3, so β = 0; before it θ = δ = 0.5, n = 20, m = 8 and τ = 1,
// and the skill was practised a day before, R = 0.9. 8 first verdicts are right and 2 not yet.
// The document prints three decimals, and τ two.
func TestForgettingPage(t *testing.T) {
	g := workbookGraphs(t)
	add := skill("addition")
	const t0 = 1762041600
	before := lm.State{Kind: lm.StateKind, Modules: []lm.ModuleState{{
		Module: workbook, Level: 0.5, AnswerWeight: 20, StartBand: lm.BandRef{Place: 1, ID: "grade-2"},
		Skills: []lm.SkillState{{Skill: add, Correction: 0.5, AnswerWeight: 8, EvidenceCount: 3, LastPractice: t0, Stability: 1}},
	}}}
	e := page(10, 8, []string{add}, t0+86400)
	e.Task.Difficulty = 3
	s := mustAppend(t, before, e, g)
	m, sk := s.Modules[0], skillOf(t, s, add)
	near3 := func(what string, got, want float64) {
		t.Helper()
		if math.Abs(got-want) > 1e-3 {
			t.Errorf("%s = %.6f, want %.3f", what, got, want)
		}
	}
	near3("θ", m.Level, 0.528)
	near3("δ", sk.Correction, 0.580)
	if math.Abs(sk.Stability-1.52) > 0.005 { // printed to two decimals
		t.Errorf("τ = %.6f, want 1.52", sk.Stability)
	}
	if math.Abs(m.AnswerWeight-22.5) > 1e-12 || math.Abs(sk.AnswerWeight-10.5) > 1e-12 {
		t.Errorf("n %v and m %v, want 22.5 and 10.5", m.AnswerWeight, sk.AnswerWeight)
	}
	if sk.Mastered || sk.FailuresInARow != 0 || sk.LastPractice != t0+86400 || sk.EvidenceCount != 4 {
		t.Errorf("the skill after the page: %+v", sk)
	}

	// A time that came back: no forgetting, R = 1, and the last practice stays where it was. A page
	// all right at R = 1 leaves τ as it is.
	back := page(10, 10, []string{add}, t0-1000)
	back.Task.Difficulty = 3
	s = mustAppend(t, before, back, g)
	if sk := skillOf(t, s, add); sk.LastPractice != t0 || sk.Stability != 1 {
		t.Errorf("a time that came back: last practice %d, τ %v", sk.LastPractice, sk.Stability)
	}
}

// In the trial series P is without retention: a pause between its pieces of evidence moves θ no
// differently, while τ and the last practice are written as usual.
func TestForgettingTrial(t *testing.T) {
	g := workbookGraphs(t)
	add := []string{skill("addition")}
	soon := mustAppend(t, mustAppend(t, lm.State{}, page(4, 3, add, 0), g), page(4, 4, add, 60), g)
	late := mustAppend(t, mustAppend(t, lm.State{}, page(4, 3, add, 0), g), page(4, 4, add, 30*86400), g)
	if len(late.Modules[0].Trial) == 0 || soon.Modules[0].Level != late.Modules[0].Level {
		t.Errorf("a pause in the trial series moved θ: %v and %v", soon.Modules[0].Level, late.Modules[0].Level)
	}
	a, b := skillOf(t, soon, add[0]), skillOf(t, late, add[0])
	if !(b.Stability > a.Stability && a.Stability > 1) || b.LastPractice != 30*86400 {
		t.Errorf("τ in the trial series: %v after a minute, %v after a month", a.Stability, b.Stability)
	}
}

// ready is a learner whose next right answer without help masters addition at band 1.
func ready(t *testing.T) lm.State {
	t.Helper()
	return lm.State{Kind: lm.StateKind, Modules: []lm.ModuleState{{
		Module: workbook, Level: 3, AnswerWeight: 50, StartBand: lm.BandRef{Place: 1, ID: "grade-2"},
		Skills: []lm.SkillState{{Skill: skill("addition"), Correction: 1, AnswerWeight: 20, EvidenceCount: 10,
			LastPractice: 100, Stability: 1}},
	}}}
}

// Help (docs/architecture.md, section 3.3): an answer with help moves the step, the counts and τ
// as it is, but masters nothing; the technique on the step and "explain again" are no help.
func TestHelp(t *testing.T) {
	g := workbookGraphs(t)
	add := skill("addition")
	cases := []struct {
		name     string
		change   func(*lm.Evidence)
		mastered bool
	}{
		{"no help", func(*lm.Evidence) {}, true},
		{"a rung after the wrong answer before one right item, the other without", func(e *lm.Evidence) {
			e.Outcome = lm.Outcome{Kind: lm.OutcomeAfterHint, Rung: 1, HasRung: true}
			e.Items[0].WrongBeforeFirst = 1
		}, true},
		{"a rung before both right items", func(e *lm.Evidence) {
			e.Outcome = lm.Outcome{Kind: lm.OutcomeAfterHint, Rung: 2, HasRung: true}
			e.Items[0].WrongBeforeFirst, e.Items[1].WrongBeforeFirst = 1, 1
		}, false},
		{"the technique on request", func(e *lm.Evidence) { e.TechniqueOnRequest = true }, false},
		{"the solution seen before", func(e *lm.Evidence) { e.SolutionSeenBefore, e.HasSolutionSeenFact = true, true }, false},
		{"a sketch", func(e *lm.Evidence) { e.Outcome = lm.Outcome{Kind: lm.OutcomeWithSketch} }, false},
		{"the technique on the step", func(e *lm.Evidence) { e.TechniqueOnStep = true }, true},
		{"explain again", func(e *lm.Evidence) { e.ExplainAgainCount = 3 }, true},
	}
	var plain lm.SkillState
	for i, c := range cases {
		e := page(3, 2, []string{add}, 100) // two right and, answered first, one not yet
		c.change(&e)
		s := mustAppend(t, ready(t), e, g)
		sk := skillOf(t, s, add)
		if sk.Mastered != c.mastered {
			t.Errorf("%s: mastered %v, want %v", c.name, sk.Mastered, c.mastered)
		}
		if i == 0 {
			plain = sk
		} else if sk.Correction != plain.Correction || sk.AnswerWeight != plain.AnswerWeight || sk.Stability != plain.Stability {
			t.Errorf("%s: help moved the step: %+v, want %+v", c.name, sk, plain)
		}
	}
}

// The outcomes (docs/architecture.md, section 3.3): a shown solution makes an item not submitted
// or only blank "not solved", S = 0; another "unreadable" or "unchecked" gives nothing; left
// unfinished, unchecked and an outcome of a later version read only the verdicts there were.
func TestOutcomes(t *testing.T) {
	g := workbookGraphs(t)
	add := skill("addition")
	e := page(6, 1, []string{add}, 100)
	e.Items[1] = lm.ItemAnswer{Place: 2}
	e.Items[2] = lm.ItemAnswer{Place: 3, Submitted: true, Unverdicted: []lm.UnverdictedReason{lm.UnverdictedUnreadableBlank, lm.UnverdictedUnreadableBlank}}
	e.Items[3] = lm.ItemAnswer{Place: 4, Submitted: true, Unverdicted: []lm.UnverdictedReason{lm.UnverdictedUnreadableBlank, lm.UnverdictedUnreadableOther}}
	e.Items[4] = lm.ItemAnswer{Place: 5, Submitted: true, Unverdicted: []lm.UnverdictedReason{lm.UnverdictedUnchecked}}
	e.Items[5] = lm.ItemAnswer{Place: 6, Submitted: true} // submitted, with neither a verdict nor a reason: nothing
	for _, c := range []struct {
		outcome  lm.OutcomeKind
		answers  float64 // N · u
		failures int32
	}{
		{lm.OutcomeSolutionShown, float64(3 * float64(3.0/5)), 1}, // one right, two not solved
		{lm.OutcomeLeftUnfinished, 1, 0},
		{lm.OutcomeUnchecked, 1, 0},
		{9, 1, 0},
	} {
		e.Outcome = lm.Outcome{Kind: c.outcome}
		before := ready(t)
		s := mustAppend(t, before, e, g)
		sk := skillOf(t, s, add)
		if math.Abs(s.Modules[0].AnswerWeight-50-c.answers) > 1e-12 || sk.FailuresInARow != c.failures {
			t.Errorf("outcome %d: n grew by %v and failures are %d, want %v and %d",
				c.outcome, s.Modules[0].AnswerWeight-50, sk.FailuresInARow, c.answers, c.failures)
		}
	}
	// A "not solved" item is an answer, and an answer needs a slot.
	e.Outcome = lm.Outcome{Kind: lm.OutcomeSolutionShown}
	e.Task.Items[1].Slots = nil
	if _, _, err := lm.Append(ready(t), e, g); !errors.Is(err, lm.ErrInvalid) {
		t.Errorf("a not solved item with no slot: %v, want ErrInvalid", err)
	}
}

// The provenance (docs/architecture.md, section 3.3; concepts 9.4 of spec): a key and skill tags from
// the author, the solver or an adult weigh 1 — an adult's key counts when the skills are not from AI;
// a key or tags from AI, a provenance not given and one of a later version weigh 0, and such
// evidence moves nothing, nor gives a module a state.
func TestProvenance(t *testing.T) {
	g := workbookGraphs(t)
	add := []string{skill("addition")}
	for _, c := range []struct {
		key, tags lm.Provenance
		weighs    bool
	}{
		{lm.ProvenanceAuthor, lm.ProvenanceAuthor, true},
		{lm.ProvenanceSolver, lm.ProvenanceAuthor, true},
		{lm.ProvenanceAdult, lm.ProvenanceAuthor, true},
		{lm.ProvenanceAI, lm.ProvenanceAuthor, false},
		{lm.ProvenanceAdult, lm.ProvenanceAI, false},
		{lm.ProvenanceUnspecified, lm.ProvenanceAuthor, false},
		{lm.ProvenanceAuthor, lm.ProvenanceUnspecified, false},
		{9, lm.ProvenanceAuthor, false},
		{lm.ProvenanceAuthor, lm.ProvenanceSolver, true},
		{lm.ProvenanceAdult, lm.ProvenanceAdult, true},
		{lm.ProvenanceAI, lm.ProvenanceAI, false},
		{lm.ProvenanceAuthor, 9, false},
	} {
		e := page(3, 2, add, 200)
		e.Task.KeyProvenance, e.Task.SkillTagsProvenance = c.key, c.tags
		e.Traps = []lm.FoundTrap{{ID: workbook + "/trap/t", Skills: add, ItemPlaces: []int32{3}}}
		moved := !sameState(t, mustAppend(t, ready(t), e, g), ready(t))
		created := len(mustAppend(t, lm.State{}, e, g).Modules) > 0
		if moved != c.weighs || created != c.weighs {
			t.Errorf("key %d, tags %d: moved %v, gave a module a state %v; want %v", c.key, c.tags, moved, created, c.weighs)
		}
	}
}

// The traps (docs/architecture.md, section 3.3): a trap found goes to its skills that have a state,
// once per piece of evidence, with the time of the latest; with none such, to the module.
func TestTraps(t *testing.T) {
	g := workbookGraphs(t)
	add, sub := skill("addition"), skill("subtraction")
	trap := func(name string, skills ...string) lm.FoundTrap {
		return lm.FoundTrap{ID: workbook + "/trap/" + name, Skills: skills, ItemPlaces: []int32{1, 2}}
	}
	e := page(2, 0, []string{add}, 500)
	e.Traps = []lm.FoundTrap{
		trap("borrow", add, skill("essay")),
		trap("loose"),
		trap("other", sub), // subtraction has no state: the trap is the module's
	}
	s := mustAppend(t, ready(t), e, g)
	// The same traps again, at a time that came back.
	e.PracticeTime = 300
	s = mustAppend(t, s, e, g)
	sk := skillOf(t, s, add)
	want := []lm.TrapCount{{ID: workbook + "/trap/borrow", Count: 2, Last: 500}}
	if fmt.Sprint(sk.Traps) != fmt.Sprint(want) {
		t.Errorf("the traps of addition: %+v, want %+v", sk.Traps, want)
	}
	wantModule := []lm.TrapCount{{ID: workbook + "/trap/loose", Count: 2, Last: 500}, {ID: workbook + "/trap/other", Count: 2, Last: 500}}
	if fmt.Sprint(s.Modules[0].Traps) != fmt.Sprint(wantModule) {
		t.Errorf("the traps of the module: %+v, want %+v", s.Modules[0].Traps, wantModule)
	}
	if len(s.Modules[0].Skills) != 1 {
		t.Errorf("a trap gave a skill a state: %+v", s.Modules[0].Skills)
	}
}

// Reasons of submissions on an item never submitted are no evidence api writes; and evidence that
// weighs nothing needs no slot at its answers, since it is never read.
func TestAppendRefusesItems(t *testing.T) {
	g := workbookGraphs(t)
	e := page(2, 1, []string{skill("addition")}, 100)
	e.Items[1] = lm.ItemAnswer{Place: 2, Unverdicted: []lm.UnverdictedReason{lm.UnverdictedUnreadableBlank}}
	if _, _, err := lm.Append(ready(t), e, g); !errors.Is(err, lm.ErrInvalid) {
		t.Errorf("reasons on an item never submitted: %v, want ErrInvalid", err)
	}
	ai := page(2, 1, []string{skill("addition")}, 100)
	ai.Outcome = lm.Outcome{Kind: lm.OutcomeSolutionShown}
	ai.Items[1] = lm.ItemAnswer{Place: 2}
	ai.Task.Items[1].Slots = nil
	ai.Task.SkillTagsProvenance = lm.ProvenanceAI
	if s, _, err := lm.Append(ready(t), ai, g); err != nil || !sameState(t, s, ready(t)) {
		t.Errorf("evidence that weighs nothing, with a not solved item and no slot: %v", err)
	}
}
