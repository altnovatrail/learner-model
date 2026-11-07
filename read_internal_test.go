package learnermodel

import "testing"

// Readout.Task gives the P the step takes of evidence at the time of reading
// (docs/architecture.md, section 3.2): in the trial series and past it, with a pause and with the
// transfer, P read at the evidence's practice time is the step's P to the bit. The learner starts
// at the second band, θ₀ = 2.5, so that the first evidence of the module, read before it, holds
// the step to the same start.
func TestTaskIsTheStepsChance(t *testing.T) {
	const module = "example/workbook"
	s := func(name string) string { return module + "/skill/" + name }
	bands := []Band{{Place: 1, ID: "grade-2", Grades: []int32{2}}, {Place: 2, ID: "grade-3", Grades: []int32{3}}}
	g, err := NewGraph(module, bands,
		[]Skill{{ID: s("a"), Assessability: AssessabilityVerdict}, {ID: s("b"), Assessability: AssessabilityVerdict}},
		[]RetiredSkill{{ID: s("old"), ReplacedBy: []string{s("a"), s("b")}}},
		[]Task{{ID: module + "/task/g", Band: BandRef{Place: 2, ID: "grade-3"}, Difficulty: 3, Primary: []string{s("a"), s("b")}}})
	if err != nil {
		t.Fatal(err)
	}
	graphs := Graphs{module: {Graph: g, Current: true}}
	page := func(at int64, primary ...string) Evidence {
		e := Evidence{
			Kind: EvidenceKind,
			Task: Task{ID: module + "/task/p", Band: BandRef{Place: 2, ID: "grade-3"}, Difficulty: 2, Primary: primary,
				KeyProvenance: ProvenanceAuthor, SkillTagsProvenance: ProvenanceAuthor},
			StartBand: BandRef{Place: 2, ID: "grade-3"}, Outcome: Outcome{Kind: OutcomeFirstSubmission}, PracticeTime: at,
		}
		for i := range 4 {
			item := TaskItem{Place: int32(i + 1), Slots: []Slot{{Check: CheckKindChoice, Options: int32(2 + i), HasOptions: true}}}
			if i == 3 {
				item.OwnSkills, item.Primary = true, []string{s("old")}
			}
			e.Task.Items = append(e.Task.Items, item)
			verdict := VerdictCorrect
			if i == 1 {
				verdict = VerdictNotYet
			}
			e.Items = append(e.Items, ItemAnswer{Place: int32(i + 1), Submitted: true, FirstVerdict: verdict, HasFirstVerdict: true})
		}
		return e
	}

	var state State
	for n := range 6 {
		at := int64(n) * 9 * 86400 // pauses of nine days
		e := page(at, s("a"), s("b"))
		r, err := Read(state, graphs, at, Grade{Value: 3, Known: true}, []string{module})
		if err != nil {
			t.Fatal(err)
		}
		read, err := r.Task(e.Task)
		if err != nil {
			t.Fatal(err)
		}
		in, err := checkInput(e, graphs, nil)
		if err != nil {
			t.Fatal(err)
		}
		ms := newModule(module, e.StartBand)
		if i, ok := state.moduleIndex(module); ok {
			ms = state.Modules[i]
		}
		answers, _ := ms.answers(e, in.items, &g, ms.inTrial())
		for i, a := range answers {
			if p := read.Items[i]; !p.HasVerdict || p.Probability != a.p {
				t.Errorf("evidence %d, item %d: read %v, the step %v", n, i+1, p, a.p)
			}
		}
		if state, _, err = Append(state, e, graphs); err != nil {
			t.Fatal(err)
		}
	}
	if state.Modules[0].AnswerWeight < trialWeight {
		t.Errorf("the pages never left the trial series: n = %v", state.Modules[0].AnswerWeight)
	}
}
