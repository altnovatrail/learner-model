package learnermodel

import (
	"math"
	"strconv"
	"testing"

	"github.com/altnovatrail/learner-model/detmath"
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
				KeyProvenance: ProvenanceAuthor, SkillTagsProvenance: ProvenanceAuthor,
				Items: []TaskItem{{Place: 1, Slots: []Slot{{Check: CheckKindChoice, Options: 5, HasOptions: true}}}}},
			StartBand: BandRef{Place: 1, ID: "1-2"},
			Items:     []ItemAnswer{{Place: 1, Submitted: true, FirstVerdict: verdict, HasFirstVerdict: true}},
		}
		items := pairItems(e)
		answers, _ := ms.answers(e, items, &g, false)
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
	// The guessing of an item: a choice of three is 1/3, two of them 1/9, a number or a time beside
	// one 0.
	choice3 := Slot{Check: CheckKindChoice, Options: 3, HasOptions: true}
	for _, c := range []struct {
		slots []Slot
		want  float64
	}{
		{[]Slot{choice3}, float64(1.0 / 3)},
		{[]Slot{choice3, choice3}, float64(1.0 / 9)},
		{[]Slot{choice3, {Check: CheckKindNumber}}, 0},
		{[]Slot{choice3, {Check: CheckKindTime}}, 0},
		{[]Slot{{Check: CheckKindTime}}, 0},
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

// The example of docs/architecture.md, section 3.2: a skill mastered at the band `3-4`, shift
// 2.5, with θ + δ = 4.5, n = 100 and m = 20. The document prints three decimals.
func TestForgettingExample(t *testing.T) {
	const day = 86400
	near3 := func(what string, got, want float64) {
		t.Helper()
		if math.Abs(got-want) > 1e-3 {
			t.Errorf("%s = %.6f, want %.3f", what, got, want)
		}
	}
	nearTo(t, "v", uncertainty(100, 20), 0.2327)
	near3("λ", cautiousLevel(4, 0.5, 100, 20), 4.018)
	q := masteryChance(4, 100, SkillState{Correction: 0.5, AnswerWeight: 20, MasteryBand: 2})
	near3("q", q, 0.820)
	days := func(r ReviewDate, from int64) float64 { return float64(float64(r.At-from) / day) }
	near3("the review after τ = 1", days(reviewDate(q, 1, 0), 0), 1.288)
	near1 := func(what string, got, want float64) {
		t.Helper()
		if math.Abs(got-want) > 0.05 {
			t.Errorf("%s = %.6f, want %.1f", what, got, want)
		}
	}
	near1("the review after τ = 10", days(reviewDate(q, 10, 0), 0), 12.9)

	// Back after a day and right: R(1; 1) = 0.9, τ = 2.5, and the next review in 3.2 days.
	r := retentionAt(SkillState{LastPractice: 0, Stability: 1}, day)
	nearTo(t, "R(1; 1)", r, 0.9)
	tau := stabilityAfter(1, r, 1)
	nearTo(t, "τ after a right answer at R = 0.9", tau, 2.5)
	near1("the review after it", days(reviewDate(q, tau, day), day), 3.2)

	// Back after a summer break of 60 days at τ = 10, on difficulty 3 of `3-4` with five options.
	r = retentionAt(SkillState{LastPractice: 0, Stability: 10}, 60*day)
	near3("R(60; 10)", r, 0.645)
	sigma := detmath.Sigmoid(4.5 - beta(2, 3))
	near3("P without forgetting", chance(0.2, sigma, 1), 0.905)
	p := chance(0.2, sigma, r)
	near3("P with forgetting", p, 0.654)
	w := retentionWeight(0.2, sigma, r, chance(0.2, sigma, 1), p)
	near3("w", w, 0.246)
	near3("the move of a right answer", float64(w*(1-p)), 0.085)
	near3("the move of a wrong one", float64(w*(0-p)), -0.161)
	near1("τ after a right answer", stabilityAfter(10, r, 1), 63.3)
	nearTo(t, "τ after a wrong one", stabilityAfter(10, r, 0), 5)

	// Within a lesson R stays near 1: two minutes after an answer at τ = 1.
	nearTo(t, "R after two minutes", retentionAt(SkillState{Stability: 1}, 120), 0.99984)
}

// The edges of the rules of forgetting.
func TestForgettingEdges(t *testing.T) {
	const day = 86400
	// The pause: none before the last practice, and below 2⁵³ seconds whatever the times.
	if pauseDays(5, 10) != 0 || pauseDays(10, 10) != 0 || pauseDays(2*day, 0) != 2 {
		t.Errorf("pauseDays of plain times")
	}
	if got := pauseDays(math.MaxInt64, math.MinInt64); got != float64(float64(1<<53)/day) {
		t.Errorf("pauseDays of the widest times = %v", got)
	}
	if retentionAt(SkillState{LastPractice: 100, Stability: 1}, 50) != 1 {
		t.Errorf("a time before the last practice forgets")
	}

	// The weight: 1 without forgetting, 0 where σ is 1, and the formula of section 3.2 elsewhere.
	weight := func(c, sigma, rbar float64) float64 {
		return retentionWeight(c, sigma, rbar, chance(c, sigma, 1), chance(c, sigma, rbar))
	}
	if weight(0.2, 0.7, 1) != 1 || weight(0, 1, 0.5) != 0 {
		t.Errorf("the weight at its ends")
	}
	for _, c := range []float64{0, 0.2, float64(1.0 / 9)} {
		for _, sigma := range []float64{0.01, 0.5, 0.99} {
			rbar := 0.7
			p1, p := chance(c, sigma, 1), chance(c, sigma, rbar)
			want := float64(rbar * p1 * (1 - p1) / (p * (1 - p)))
			if got := weight(c, sigma, rbar); math.Abs(got-want) > 1e-12 {
				t.Errorf("w(c = %v, σ = %v) = %v, want %v", c, sigma, got, want)
			}
		}
	}
	// No guessing and σ at its smallest: the weight is still a plain number.
	if w := weight(0, detmath.Sigmoid(-708), 0.5); math.Abs(w-1) > 1e-12 {
		t.Errorf("w at the smallest σ = %v", w)
	}

	// The stability: S̄ = 1 is τ · (1 + g · (1 − R)) to the bit, S̄ = 0 halves it down to τ₀, and a
	// page right but for some grows it less.
	r := 0.9 // a variable, so that 1 − r is the float the code takes, not an exact constant
	if got, want := stabilityAfter(3, r, 1), float64(3*(1+float64(15*(1-r)))); got != want {
		t.Errorf("τ after right answers = %v, want %v", got, want)
	}
	if stabilityAfter(3, 0.9, 0) != 1.5 || stabilityAfter(1.5, 0.9, 0) != 1 {
		t.Errorf("τ after wrong answers")
	}
	nearTo(t, "τ after 8 of 10 at R = 0.9", stabilityAfter(1, 0.9, 0.8), 1.8119)
	nearTo(t, "τ after half at R = 0.9", stabilityAfter(1, 0.9, 0.5), 1.1180)
	if got := stabilityAfter(math.MaxFloat64, 0, 1); got != math.MaxFloat64 {
		t.Errorf("τ past the largest float = %v", got)
	}

	// The review date: "already" below the threshold, never sooner than τ₀, rounded up to the
	// second, and saturating past the largest time.
	if !reviewDate(0.7, 1, 0).Already {
		t.Errorf("q below the threshold is not already due")
	}
	if got := reviewDate(0.72, 1, 1000); got.Already || got.At != 1000+day {
		t.Errorf("a mastery barely held is due at %+v, want a day after its practice", got)
	}
	q := 0.9
	rStar := float64(masteryThreshold / q)
	exact := float64(float64(5*(float64(1/float64(rStar*rStar))-1)/detmath.RetentionFactor) * day)
	if got := reviewDate(q, 5, 0).At; float64(got) < exact || float64(got)-exact >= 1 {
		t.Errorf("the review at %d s, want %v rounded up", got, exact)
	}
	// Δt* lands on a whole second here, and its three roundings bring it to a second at which R · q
	// still holds the threshold: the date moves on to the next one.
	edge, tauEdge := 0.7509902336184128, 3.0
	at := reviewDate(edge, tauEdge, 0).At
	holds := func(s int64) bool {
		return float64(retentionAt(SkillState{Stability: tauEdge}, s)*edge) > masteryThreshold
	}
	if at != 101357 || holds(at) || !holds(at-1) {
		t.Errorf("the review on the edge is at %d s, want 101357, the first second past the threshold", at)
	}
	if reviewDate(q, 1e300, 0).At != math.MaxInt64 || addSeconds(math.MaxInt64-5, 10) != math.MaxInt64 ||
		addSeconds(5, math.Inf(1)) != math.MaxInt64 || addSeconds(5, -3) != 5 {
		t.Errorf("the review date does not saturate")
	}
}
