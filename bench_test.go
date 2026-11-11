package learnermodel_test

import (
	"fmt"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

// The benchmark of Readout.Task on a pool the size of olympiad maths 0.1.1, as the engine reads it
// (Д18 in RUN.md): the engine asks for the P of every task of the pool before each choice, 772
// calls. The module is synthetic, of the package's shape: three bands, 19 skills with a verdict,
// 772 tasks of one item and one primary skill each, 32 of them with a supporting skill, slots of a
// number, a choice of five or a fraction, and 19 templates. Read is api's call, not the engine's:
// it runs once, before the loop, which times Task alone.

const poolModule = "example/pool"

const poolSkills, poolTasks = 19, 772

func poolSkill(k int) string { return fmt.Sprintf("%s/skill/s%02d", poolModule, k) }

var poolBands = []lm.Band{
	{Place: 1, ID: "grades-1-2", Grades: []int32{1, 2}},
	{Place: 2, ID: "grades-3-4", Grades: []int32{3, 4}},
	{Place: 3, ID: "grades-5-6", Grades: []int32{5, 6}},
}

// poolTask is the k-th task of the pool, by ID: the tasks of a skill stand together, as the pool
// of a package is ordered by ID.
func poolTask(k int) lm.Task {
	b := poolBands[k%3]
	slot := lm.Slot{Check: lm.CheckKindNumber}
	switch {
	case k%110 == 7:
		slot = lm.Slot{Check: lm.CheckKindFraction}
	case k%4 == 0:
		slot = lm.Slot{Check: lm.CheckKindChoice, Options: 5, HasOptions: true}
	}
	s := k * poolSkills / poolTasks
	t := lm.Task{
		ID: fmt.Sprintf("%s/task/t%03d", poolModule, k), Template: k%41 == 0,
		Band: lm.BandRef{Place: b.Place, ID: b.ID}, Difficulty: int32(1 + k%5),
		Primary:       []string{poolSkill(s)},
		KeyProvenance: lm.ProvenanceAuthor, SkillTagsProvenance: lm.ProvenanceAuthor,
		Items: []lm.TaskItem{{Place: 1, Slots: []lm.Slot{slot}}},
	}
	if k%24 == 0 {
		t.Supporting = []string{poolSkill((s + 1) % poolSkills)}
	}
	return t
}

func poolGraphs(tb testing.TB) (lm.Graphs, []lm.Task) {
	tb.Helper()
	var skills []lm.Skill
	for k := range poolSkills {
		skills = append(skills, lm.Skill{ID: poolSkill(k), Assessability: lm.AssessabilityVerdict})
	}
	var tasks []lm.Task
	for k := range poolTasks {
		tasks = append(tasks, poolTask(k))
	}
	g, err := lm.NewGraph(poolModule, poolBands, skills, nil, tasks)
	if err != nil {
		tb.Fatal(err)
	}
	return lm.Graphs{poolModule: {Graph: g, Current: true}}, tasks
}

// poolAnswer is evidence of a task of the pool at a time, right or not yet.
func poolAnswer(t lm.Task, right bool, at int64) lm.Evidence {
	v := lm.VerdictNotYet
	if right {
		v = lm.VerdictCorrect
	}
	return lm.Evidence{
		Kind: lm.EvidenceKind, Task: t, StartBand: lm.BandRef{Place: 2, ID: "grades-3-4"},
		Outcome: lm.Outcome{Kind: lm.OutcomeFirstSubmission}, PracticeTime: at,
		Items: []lm.ItemAnswer{{Place: 1, Submitted: true, FirstVerdict: v, HasFirstVerdict: true}},
	}
}

const poolTime = int64(1762041600)

// poolStates are the states the benchmark reads: a learner who has just connected the module; one
// in the trial series; one half way, the first ten skills practised a day apart in the last ten
// days; and one far on, every skill practised over a year, most of them mastered, and some of
// those due for review at poolTime.
func poolStates(tb testing.TB, g lm.Graphs, tasks []lm.Task) []struct {
	name  string
	state lm.State
} {
	tb.Helper()
	appendAll := func(s lm.State, es ...lm.Evidence) lm.State {
		for _, e := range es {
			var err error
			if s, _, err = lm.Append(s, e, g); err != nil {
				tb.Fatal(err)
			}
		}
		return s
	}
	// The i-th task of a skill, by ID.
	of := func(s, i int) lm.Task {
		var own []lm.Task
		for _, t := range tasks {
			if t.Primary[0] == poolSkill(s) {
				own = append(own, t)
			}
		}
		return own[i%len(own)]
	}
	var trial, middle, far []lm.Evidence
	for k := range 3 {
		trial = append(trial, poolAnswer(of(k, 1), k != 1, poolTime-int64(3-k)*day))
	}
	for s := range 10 {
		for i := range 6 {
			at := poolTime - int64(10-s)*day + int64(i)*600
			middle = append(middle, poolAnswer(of(s, i), (s+i)%3 != 0, at))
		}
	}
	// Eight rounds 45 days apart, but a skill in four stops after two: its review has come.
	for round := range 8 {
		for s := range poolSkills {
			if s%4 == 3 && round >= 2 {
				continue
			}
			at := poolTime - int64((8-round)*45+s)*day
			far = append(far, poolAnswer(of(s, round), s%5 != 4 || round%2 == 0, at))
		}
	}
	return []struct {
		name  string
		state lm.State
	}{
		{"new", lm.State{}},
		{"trial", appendAll(lm.State{}, trial...)},
		{"middle", appendAll(lm.State{}, middle...)},
		{"far", appendAll(lm.State{}, far...)},
	}
}

// BenchmarkReadoutTask reads the P of every task of the pool, as the engine does before a choice.
// Its ns/op is a whole pool; ns/task is a call.
func BenchmarkReadoutTask(b *testing.B) {
	g, tasks := poolGraphs(b)
	for _, c := range poolStates(b, g, tasks) {
		r, err := lm.Read(c.state, g, poolTime, lm.Grade{Value: 3, Known: true}, []string{poolModule})
		if err != nil {
			b.Fatal(err)
		}
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for _, t := range tasks {
					if _, err := r.Task(t); err != nil {
						b.Fatal(err)
					}
				}
			}
			b.ReportMetric(float64(float64(b.Elapsed().Nanoseconds())/float64(b.N*len(tasks))), "ns/task")
		})
	}
}
