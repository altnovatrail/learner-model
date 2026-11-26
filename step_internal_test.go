package learnermodel

import (
	"math"
	"slices"
	"strconv"
	"testing"

	"github.com/altnovatrail/learner-model/detmath"
)

// The first worked example of standalone's SPEC.md, at commit 6b29dd9 of its main: the step of Elo
// run from θ = 0 and no answers, on tasks of `1-2`, two topics — what evidence never meets, since
// its first five answers are the trial series, so the test runs the step past it. It holds to the
// four decimals printed there (standalone_test.go has the other examples).

const internalTolerance = 1e-4

// The other tolerances of the tests of the package: sumTolerance, of a float sum of a few terms
// against its exact value; productTolerance, of a product of a few exact fractions; threeDecimals
// and oneDecimal, of a value a document prints to three decimals and to one.
const (
	sumTolerance     = 1e-12
	productTolerance = 1e-15
	threeDecimals    = 1e-3
	oneDecimal       = 0.05
)

// guessOfFive is c, the chance of guessing a choice of five options.
const guessOfFive = 0.2

// summerBreak is the pause of the examples of forgetting: 60 days.
const summerBreak = 60 * day

// The learner of the examples of forgetting of docs/architecture.md: θ and δ of a skill mastered at
// the band `3-4`, its place, and the weights of answers n of the module and m of the skill.
const (
	exampleLevel      = 4
	exampleCorrection = 0.5
	examplePlace      = 2
	exampleN          = 100
	exampleM          = 20
)

// day is a day in seconds, for the times of the tests of the package.
const day = secondsPerDay

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
	// noMasteryYet is the step after which SPEC.md works out that the skill is not mastered yet.
	const noMasteryYet = 9
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
		if i+1 == noMasteryYet {
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

// The numbers of standalone's SPEC.md — the scale, the step and the floor of θ's step — and of the
// example of the scale of docs/architecture.md, that the tests through Append cannot see.
func TestModelNumbers(t *testing.T) {
	// A child of grade 3 at θ₀ = 2.5 on difficulty 3 of `3-4` with five options: P = 0.6.
	nearTo(t, "P at the start", probability(guessOfFive, shift(2)-beta(2, middleDifficulty)), 0.6)
	// The scale: a child at the start of their level expects these at the five difficulties.
	for d, want := range []float64{0.9046, 0.7848, 0.6000, 0.4152, 0.2954} {
		difficulty := int32(d + minDifficulty)
		nearTo(t, "P at difficulty "+itoa(difficulty), probability(guessOfFive, -beta(firstPlace, difficulty)), want)
	}
	// The floor of θ's step, the first answer: K_θ at its floor, K_δ, v and √v.
	nearTo(t, "K_θ at n = 79", levelStep(79), levelStepFloor)
	nearTo(t, "K_δ at m = 9", correctionStep(9), 0.2759)
	nearTo(t, "v at n = 80, m = 10", uncertainty(80, 10), 0.3045)
	nearTo(t, "√v at n = 80, m = 10", math.Sqrt(uncertainty(80, 10)), 0.5518)
	// The weight of an attempt: the table of docs/architecture.md.
	for n, want := range map[int]float64{1: 1, 2: 1.5, 4: 2, 10: 2.5} {
		if got := float64(float64(n) * seriesWeight(n)); math.Abs(got-want) > sumTolerance {
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
		if got := guess(c.slots); math.Abs(got-c.want) > productTolerance {
			t.Errorf("guess(%+v) = %v, want %v", c.slots, got, c.want)
		}
	}
	if gridCells(-1) != 0 || gridCells(1e300) != maxGridCells || gridCells(NaN) != 0 || gridCells(7.9) != 7 {
		t.Errorf("gridCells does not keep to its range")
	}
	const start = 2.5
	if trialLevel(start, nil) != start {
		t.Errorf("the trial level with no answers is not the start")
	}
}

// The example of forgetting of docs/architecture.md: a skill mastered at the band `3-4`,
// shift 2.5, with θ + δ = 4.5, n = 100 and m = 20. The document prints three decimals.
func TestForgettingExample(t *testing.T) {
	near3 := func(what string, got, want float64) {
		t.Helper()
		if math.Abs(got-want) > threeDecimals {
			t.Errorf("%s = %.6f, want %.3f", what, got, want)
		}
	}
	nearTo(t, "v", uncertainty(exampleN, exampleM), 0.2327)
	near3("λ", cautiousLevel(exampleLevel, exampleCorrection, exampleN, exampleM), 4.018)
	q := masteryChance(exampleLevel, exampleN,
		SkillState{Correction: exampleCorrection, AnswerWeight: exampleM, MasteryBand: examplePlace})
	near3("q", q, 0.820)
	days := func(r ReviewDate, from int64) float64 { return float64(float64(r.At-from) / day) }
	near3("the review after τ = 1", days(reviewDate(q, 1, 0), 0), 1.288)
	near1 := func(what string, got, want float64) {
		t.Helper()
		if math.Abs(got-want) > oneDecimal {
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
	r = retentionAt(SkillState{LastPractice: 0, Stability: 10}, summerBreak)
	near3("R(60; 10)", r, 0.645)
	sigma := detmath.Sigmoid(exampleLevel + exampleCorrection - beta(examplePlace, middleDifficulty))
	near3("P without forgetting", chance(guessOfFive, sigma, 1), 0.905)
	p := chance(guessOfFive, sigma, r)
	near3("P with forgetting", p, 0.654)
	w := retentionWeight(guessOfFive, sigma, r, chance(guessOfFive, sigma, 1), p)
	near3("w", w, 0.246)
	near3("the move of a right answer", float64(w*(1-p)), 0.085)
	near3("the move of a wrong one", float64(w*(0-p)), -0.161)
	near1("τ after a right answer", stabilityAfter(10, r, 1), 63.3)
	nearTo(t, "τ after a wrong one", stabilityAfter(10, r, 0), 5)

	// Within a lesson R stays near 1: two minutes after an answer at τ = 1.
	const twoMinutes = 120 // in seconds
	nearTo(t, "R after two minutes", retentionAt(SkillState{Stability: 1}, twoMinutes), 0.99984)
}

// The example of forgetting of docs/architecture.md, with the forgetting of the level: the
// learner of the example, mastered at `3-4`, comes back to the module after 60 days and answers
// difficulty 3 of `3-4` with five options right. The pause of the module is 60 days, so n falls
// from 100 to 20 before the step: K_θ is 0.1, not its floor of 0.05, and θ moves twice as far,
// while δ moves as it would without the rule. v widens, and the next review comes sooner. The
// document prints three decimals, and the review one.
func TestLevelForgettingExample(t *testing.T) {
	const module = "example/olympiad"
	a, b := module+"/skill/a", module+"/skill/b"
	band := BandRef{Place: 2, ID: "3-4"}
	g, err := NewGraph(module, []Band{{Place: 1, ID: "1-2"}, {Place: 2, ID: "3-4"}, {Place: 3, ID: "5-6"}},
		[]Skill{{ID: a, Assessability: AssessabilityVerdict}, {ID: b, Assessability: AssessabilityVerdict}}, nil,
		[]Task{{ID: module + "/task/a", Band: band, Difficulty: 3, Primary: []string{a}}})
	if err != nil {
		t.Fatal(err)
	}
	ms := ModuleState{Module: module, Level: exampleLevel, AnswerWeight: exampleN, StartBand: band,
		Skills: []SkillState{{Skill: a, Correction: exampleCorrection, AnswerWeight: exampleM, EvidenceCount: 20,
			LastPractice: 0, Stability: 10, Mastered: true, MasteryBand: examplePlace}}}
	e := Evidence{
		Kind: EvidenceKind,
		Task: Task{ID: module + "/task/x", Band: band, Difficulty: 3, Primary: []string{a},
			KeyProvenance: ProvenanceAuthor, SkillTagsProvenance: ProvenanceAuthor,
			Items: []TaskItem{{Place: 1, Slots: []Slot{{Check: CheckKindChoice, Options: 5, HasOptions: true}}}}},
		StartBand: band, PracticeTime: summerBreak,
		Items: []ItemAnswer{{Place: 1, Submitted: true, FirstVerdict: VerdictCorrect, HasFirstVerdict: true}},
	}
	near3 := func(what string, got, want float64) {
		t.Helper()
		if math.Abs(got-want) > threeDecimals {
			t.Errorf("%s = %.6f, want %.3f", what, got, want)
		}
	}
	nearTo(t, "K_θ at n = 100", levelStep(exampleN), levelStepFloor)
	nearTo(t, "K_θ at n = 20", levelStep(levelWeight), 0.1)
	// The same learner, but another skill of the module was practised a day before the answer: the
	// pause of the module is a day, the rule stays off, and the answer of skill a is the same.
	control := ms
	control.Skills = append(slices.Clone(ms.Skills), SkillState{Skill: b, AnswerWeight: 3, EvidenceCount: 1,
		LastPractice: summerBreak - day, Stability: 1})
	for _, m := range []*ModuleState{&ms, &control} {
		if !m.step(e, pairItems(e), &g) {
			t.Fatal("the answer moved nothing")
		}
	}
	sk, ck := ms.Skills[0], control.Skills[0]
	nearTo(t, "the move of θ", ms.Level-exampleLevel, 0.0085)
	nearTo(t, "the move of θ without the rule", control.Level-exampleLevel, 0.0043)
	nearTo(t, "the move of δ", sk.Correction-exampleCorrection, 0.0170)
	if ck.Correction != sk.Correction || ck.Stability != sk.Stability {
		t.Errorf("the rule moved δ or τ: %+v, want %+v", sk, ck)
	}
	// The answer weighs 1: n after the pause is levelWeight, and without it exampleN.
	if ms.AnswerWeight != levelWeight+1 || control.AnswerWeight != exampleN+1 || sk.AnswerWeight != exampleM+1 {
		t.Errorf("n %v, %v without the rule, and m %v after the answer, want 21, 101 and 21",
			ms.AnswerWeight, control.AnswerWeight, sk.AnswerWeight)
	}
	reviewIn := func(q float64) float64 {
		return float64(float64(reviewDate(q, sk.Stability, sk.LastPractice).At-summerBreak) / day)
	}
	for _, c := range []struct {
		what        string
		m           ModuleState
		v, q, after float64
	}{
		{"", ms, 0.465, 0.793, 58.7},
		{" without the rule", control, 0.228, 0.824, 84.9},
	} {
		near3("v"+c.what, uncertainty(c.m.AnswerWeight, c.m.Skills[0].AnswerWeight), c.v)
		q := masteryChance(c.m.Level, c.m.AnswerWeight, c.m.Skills[0])
		near3("q"+c.what, q, c.q)
		if got := reviewIn(q); math.Abs(got-c.after) > oneDecimal {
			t.Errorf("the next review%s in %.3f days, want %.1f", c.what, got, c.after)
		}
	}
	if !sk.Mastered || sk.MasteryBand != examplePlace {
		t.Errorf("the mastery after the answer: %v at %d, want held at 2", sk.Mastered, sk.MasteryBand)
	}
}

// The edges of the rules of forgetting.
func TestForgettingEdges(t *testing.T) {
	// The pause: none before the last practice, and below 2⁵³ seconds whatever the times.
	if pauseDays(5, 10) != 0 || pauseDays(10, 10) != 0 || pauseDays(2*day, 0) != 2 {
		t.Errorf("pauseDays of plain times")
	}
	if got := pauseDays(math.MaxInt64, math.MinInt64); got != float64(float64(maxPauseSeconds)/day) {
		t.Errorf("pauseDays of the widest times = %v", got)
	}
	if retentionAt(SkillState{LastPractice: 100, Stability: 1}, 50) != 1 {
		t.Errorf("a time before the last practice forgets")
	}

	// The weight: 1 without forgetting, 0 where σ is 1, and the formula of retentionWeight
	// elsewhere.
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
			if got := weight(c, sigma, rbar); math.Abs(got-want) > sumTolerance {
				t.Errorf("w(c = %v, σ = %v) = %v, want %v", c, sigma, got, want)
			}
		}
	}
	// No guessing and σ at its smallest: the weight is still a plain number.
	if w := weight(0, detmath.Sigmoid(math.Inf(-1)), 0.5); math.Abs(w-1) > sumTolerance {
		t.Errorf("w at the smallest σ = %v", w)
	}

	// The stability: S̄ = 1 is τ · (1 + g · (1 − R)) to the bit, S̄ = 0 halves it down to τ₀, and a
	// page right but for some grows it less.
	const tau = 3.0
	r := 0.9 // a variable, so that 1 − r is the float the code takes, not an exact constant
	if got, want := stabilityAfter(tau, r, 1), float64(tau*(1+float64(stabilityGrowth*(1-r)))); got != want {
		t.Errorf("τ after right answers = %v, want %v", got, want)
	}
	if stabilityAfter(tau, 0.9, 0) != tau*stabilityLoss || stabilityAfter(tau*stabilityLoss, 0.9, 0) != stabilityFloor {
		t.Errorf("τ after wrong answers")
	}
	nearTo(t, "τ after 8 of 10 at R = 0.9", stabilityAfter(1, 0.9, 0.8), 1.8119)
	nearTo(t, "τ after half at R = 0.9", stabilityAfter(1, 0.9, 0.5), 1.1180)
	if got := stabilityAfter(math.MaxFloat64, 0, 1); got != math.MaxFloat64 {
		t.Errorf("τ past the largest float = %v", got)
	}

	// The review date: "already" below the threshold, never sooner than τ₀, rounded up to the
	// second, and saturating past the largest time.
	const (
		belowThreshold = 0.7  // q just below masteryThreshold
		barelyHeld     = 0.72 // q just above it
		practice       = 1000 // the time of the last practice, in seconds
	)
	if !reviewDate(belowThreshold, 1, 0).Already {
		t.Errorf("q below the threshold is not already due")
	}
	if got := reviewDate(barelyHeld, 1, practice); got.Already || got.At != practice+stabilityFloor*day {
		t.Errorf("a mastery barely held is due at %+v, want a day after its practice", got)
	}
	q, tauReview := 0.9, 5.0
	rStar := float64(masteryThreshold / q)
	exact := float64(float64(tauReview*(float64(1/float64(rStar*rStar))-1)/detmath.RetentionFactor) * day)
	if got := reviewDate(q, tauReview, 0).At; float64(got) < exact || float64(got)-exact >= 1 {
		t.Errorf("the review at %d s, want %v rounded up", got, exact)
	}
	// Δt* lands on a whole second here, and its three roundings bring it to a second at which R · q
	// still holds the threshold: the date moves on to the next one.
	edge, tauEdge := 0.7509902336184128, 3.0
	at := reviewDate(edge, tauEdge, 0).At
	holds := func(s int64) bool {
		return float64(retentionAt(SkillState{Stability: tauEdge}, s)*edge) > masteryThreshold
	}
	const edgeSecond = 101357
	if at != edgeSecond || holds(at) || !holds(at-1) {
		t.Errorf("the review on the edge is at %d s, want %d, the first second past the threshold", at, edgeSecond)
	}
	const hugeTau = 1e300 // a stability whose review date is past the largest time
	if reviewDate(q, hugeTau, 0).At != math.MaxInt64 || addSeconds(math.MaxInt64-5, 10) != math.MaxInt64 ||
		addSeconds(5, math.Inf(1)) != math.MaxInt64 || addSeconds(5, -3) != 5 {
		t.Errorf("the review date does not saturate")
	}
}
