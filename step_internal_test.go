package learnermodel

import (
	"math"
	"strconv"
	"testing"
)

// SPEC 2.6 of standalone, at commit 6b29dd9 of its main, read 2026-11-07: the step of Elo run from
// θ = 0 and no answers, on tasks of `1-2`, two topics — what evidence never meets, since its first
// five answers are the trial series, so the test runs the step past it. It holds to the four
// decimals printed there (standalone_test.go has the other examples).

const internalTolerance = 1e-4

func nearTo(t *testing.T, what string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > internalTolerance {
		t.Errorf("%s = %.6f, want %.4f", what, got, want)
	}
}

func TestStandaloneTenAnswers(t *testing.T) {
	const module = "example/olympiad"
	a, b := module+"/skill/combinatorics.enumeration", module+"/skill/parity.alternation"
	bands := []Band{{Place: 1, ID: "1-2", Grades: []int32{1, 2}}}
	g, err := NewGraph(module, bands,
		[]Skill{{ID: a, Assessability: AssessabilityVerdict}, {ID: b, Assessability: AssessabilityVerdict}}, nil,
		[]Task{{ID: module + "/task/a", Band: BandRef{Place: 1, ID: "1-2"}, Difficulty: 3, Primary: []string{a}},
			{ID: module + "/task/b", Band: BandRef{Place: 1, ID: "1-2"}, Difficulty: 3, Primary: []string{b}}})
	if err != nil {
		t.Fatal(err)
	}
	rows := []struct {
		skill        string
		d            int32
		p            float64
		correct      bool
		theta, delta float64
	}{
		{a, 3, 0.6000, true, 0.0800, 0.1600},
		{a, 3, 0.6478, true, 0.1471, 0.2942},
		{a, 4, 0.4911, false, 0.0578, 0.1156},
		{a, 3, 0.6346, true, 0.1214, 0.2427},
		{b, 3, 0.6242, false, 0.0173, -0.2497},
		{b, 2, 0.7464, true, 0.0579, -0.1531},
		{a, 4, 0.4656, true, 0.1401, 0.4209},
		{a, 4, 0.5136, true, 0.2122, 0.5765},
		{a, 4, 0.5579, true, 0.2753, 0.7125},
		{b, 3, 0.6244, true, 0.3271, -0.0165},
	}
	ms := ModuleState{Module: module, StartBand: BandRef{Place: 1, ID: "1-2"}}
	for i, r := range rows {
		step := "answer " + strconv.Itoa(i+1)
		verdict := VerdictNotYet
		if r.correct {
			verdict = VerdictCorrect
		}
		e := Evidence{
			Kind: EvidenceKind,
			Task: Task{ID: module + "/task/x", Band: BandRef{Place: 1, ID: "1-2"}, Difficulty: r.d, Primary: []string{r.skill},
				Items: []TaskItem{{Place: 1, Slots: []Slot{{Check: CheckKindChoice, Options: 5, HasOptions: true}}}}},
			StartBand: BandRef{Place: 1, ID: "1-2"},
			Items:     []ItemAnswer{{Place: 1, Submitted: true, FirstVerdict: verdict, HasFirstVerdict: true}},
		}
		items := pairItems(e)
		answers := ms.answers(e.Task, items, &g)
		if len(answers) != 1 {
			t.Fatalf("%s: %d answers", step, len(answers))
		}
		nearTo(t, "P of "+step, answers[0].p, r.p)
		if !ms.move(e, items, &g, false) {
			t.Fatalf("%s moved nothing", step)
		}
		nearTo(t, "θ after "+step, ms.Level, r.theta)
		nearTo(t, "δ after "+step, ms.correction(r.skill), r.delta)
		if ms.AnswerWeight != float64(i+1) || len(ms.Trial) != 0 {
			t.Errorf("%s: n %v, %d trial answers", step, ms.AnswerWeight, len(ms.Trial))
		}
		if i == 8 {
			// "No mastery yet": after step 9, n = 9 and m = 7, v = 0.9093 and λ = 0.0342, too low
			// for the middle task of `1-2`.
			i, _ := ms.skillIndex(a)
			sk := ms.Skills[i]
			nearTo(t, "v after step 9", uncertainty(ms.AnswerWeight, sk.AnswerWeight), 0.9093)
			// λ is printed from the rounded θ and δ of the table, so it is worked out from them.
			nearTo(t, "λ after step 9", cautiousLevel(0.2753, 0.7125, ms.AnswerWeight, sk.AnswerWeight), 0.0342)
			if sk.AnswerWeight != 7 || sk.Mastered {
				t.Errorf("after step 9: m %v, mastered %v", sk.AnswerWeight, sk.Mastered)
			}
		}
	}
}

// The numbers of SPEC 2.1, 2.2 and 2.12 of standalone and the example of docs/architecture.md,
// section 3.1, that the tests through Append cannot see.
func TestModelNumbers(t *testing.T) {
	// A child of grade 3 at θ₀ = 2.5 on difficulty 3 of `3-4` with five options: P = 0.6.
	nearTo(t, "P at the start", probability(0.2, shift(2)-beta(2, 3)), 0.6)
	// SPEC 2.1: a child at the start of their level expects these at the five difficulties.
	for d, want := range []float64{0.9046, 0.7848, 0.6000, 0.4152, 0.2954} {
		nearTo(t, "P at difficulty "+strconv.Itoa(d+1), probability(0.2, -beta(1, int32(d+1))), want)
	}
	// SPEC 2.12, the first answer: K_θ at its floor, K_δ, v and √v.
	nearTo(t, "K_θ at n = 79", levelStep(79), 0.05)
	nearTo(t, "K_δ at m = 9", correctionStep(9), 0.2759)
	nearTo(t, "v at n = 80, m = 10", uncertainty(80, 10), 0.3045)
	nearTo(t, "√v at n = 80, m = 10", math.Sqrt(uncertainty(80, 10)), 0.5518)
	// The weight of an attempt: the table of docs/architecture.md, section 3.3.
	for n, want := range map[int]float64{1: 1, 2: 1.5, 4: 2, 10: 2.5} {
		if got := float64(float64(n) * seriesWeight(n)); math.Abs(got-want) > 1e-12 {
			t.Errorf("an attempt of %d answers weighs %v, want %v", n, got, want)
		}
	}
	// The guessing of an item: a choice of three is 1/3, two of them 1/9, a number beside one 0.
	choice3 := Slot{Check: CheckKindChoice, Options: 3, HasOptions: true}
	for _, c := range []struct {
		slots []Slot
		want  float64
	}{
		{[]Slot{choice3}, float64(1.0 / 3)},
		{[]Slot{choice3, choice3}, float64(1.0 / 9)},
		{[]Slot{choice3, {Check: CheckKindNumber}}, 0},
	} {
		if got := guess(c.slots); math.Abs(got-c.want) > 1e-15 {
			t.Errorf("guess(%+v) = %v, want %v", c.slots, got, c.want)
		}
	}
	if gridCells(-1) != 0 || gridCells(1e300) != 1<<20 || gridCells(math.Float64frombits(0x7ff8000000000001)) != 0 || gridCells(7.9) != 7 {
		t.Errorf("gridCells does not keep to its range")
	}
	if trialLevel(2.5, nil) != 2.5 {
		t.Errorf("the trial level with no answers is not the start")
	}
}
