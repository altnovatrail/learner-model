package learnermodel_test

import (
	"cmp"
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	lm "github.com/altnovatrail/learner-model"
)

// Readout.Task against its reference, Task of v0.1.5 (task_reference_internal_test.go): the same
// P to the bit and the same refusal, word for word, on every candidate here.

// sameTask holds Task to its reference on a candidate.
func sameTask(t *testing.T, what string, r lm.Readout, task lm.Task) {
	t.Helper()
	got, err := r.Task(task)
	want, wantErr := r.TaskReference(task)
	switch {
	case (err == nil) != (wantErr == nil):
		t.Fatalf("%s, task %q: Task gives %v, v0.1.5 gave %v", what, task.ID, err, wantErr)
	case err != nil:
		if err.Error() != wantErr.Error() {
			t.Fatalf("%s, task %q: Task refuses with %q, v0.1.5 with %q", what, task.ID, err, wantErr)
		}
		return
	case len(got.Items) != len(want.Items):
		t.Fatalf("%s, task %q: %d items, v0.1.5 gave %d", what, task.ID, len(got.Items), len(want.Items))
	}
	for i, p := range got.Items {
		q := want.Items[i]
		if p.Place != q.Place || p.HasVerdict != q.HasVerdict || math.Float64bits(p.Probability) != math.Float64bits(q.Probability) {
			t.Fatalf("%s, task %q: item %+v, v0.1.5 gave %+v", what, task.ID, p, q)
		}
	}
}

// variants gives a candidate as it is, and as the engine or a mistake might give it: its tags and
// items in another order, and each check of Task broken once.
func variants(task lm.Task) []lm.Task {
	clone := func() lm.Task {
		c := task
		c.Primary, c.Supporting = slices.Clone(task.Primary), slices.Clone(task.Supporting)
		c.Items = slices.Clone(task.Items)
		for i := range c.Items {
			c.Items[i].Primary, c.Items[i].Slots = slices.Clone(c.Items[i].Primary), slices.Clone(c.Items[i].Slots)
		}
		return c
	}
	var out []lm.Task
	add := func(change func(*lm.Task)) {
		c := clone()
		change(&c)
		out = append(out, c)
	}
	add(func(*lm.Task) {})
	add(func(c *lm.Task) { slices.Reverse(c.Primary); slices.Reverse(c.Items) })
	add(func(c *lm.Task) {
		slices.SortFunc(c.Items, func(a, b lm.TaskItem) int { return cmp.Compare(a.Place, b.Place) })
	})
	add(func(c *lm.Task) {
		for i := range c.Items {
			slices.Reverse(c.Items[i].Primary)
		}
	})
	if len(task.Primary) > 0 {
		add(func(c *lm.Task) { c.Primary = append(c.Primary, c.Primary[0]) })
	}
	if len(task.Items) > 0 {
		add(func(c *lm.Task) { c.Items = append(c.Items, c.Items[0]) })
		add(func(c *lm.Task) {
			c.Items[0].Slots = append(c.Items[0].Slots, lm.Slot{Check: lm.CheckKindChoice, Options: 1, HasOptions: true})
		})
		add(func(c *lm.Task) {
			c.Items[0].Slots = append(c.Items[0].Slots, lm.Slot{Check: lm.CheckKindNumber, Options: 2})
		})
	}
	add(func(c *lm.Task) { c.Band.Place = 0 })
	add(func(c *lm.Task) { c.Difficulty = 6 })
	add(func(c *lm.Task) { c.ID = "example/m0/task" })
	add(func(c *lm.Task) { c.ID = "example//task/x" })
	add(func(c *lm.Task) { c.ID += "\xff" })
	add(func(c *lm.Task) { c.ID = "example/none/task/x" })
	add(func(c *lm.Task) { c.ID = "example/extra/task/x" })
	return out
}

// TestTaskIsReference reads every prefix of the golden logs at the three times of their golden
// file, and the random logs of the property, with every task of the log and its variants as
// candidates, and the pool of the benchmark in its four states.
func TestTaskIsReference(t *testing.T) {
	for _, j := range goldenJournals {
		log, g := goldenLog(t, j.name), j.graphs(t)
		var engine []string
		for _, id := range sortedKeys(g) {
			if g[id].Current {
				engine = append(engine, id)
			}
		}
		for n := range len(log) + 1 {
			s, err := lm.Recompute(log[:n], g)
			if err != nil {
				t.Fatal(err)
			}
			last := log[max(n-1, 0)].PracticeTime
			for _, rd := range []struct {
				at     int64
				grade  lm.Grade
				engine []string
			}{
				{last, lm.Grade{Value: 3, Known: true}, engine},
				{last + day, lm.Grade{Value: 3, Known: true}, engine},
				{last + 60*day, lm.Grade{}, nil},
			} {
				r := mustRead(t, s, g, rd.at, rd.grade, rd.engine...)
				for _, e := range log {
					for _, task := range variants(e.Task) {
						sameTask(t, j.name, r, task)
					}
				}
			}
		}
	}

	for seed := range uint64(30) {
		r := rand.New(rand.NewPCG(seed, 40))
		w := newWorld(t, r)
		var engine []string
		for _, m := range w.modules {
			if r.IntN(2) == 0 && w.graphs[m.id].Current {
				engine = append(engine, m.id)
			}
		}
		var s lm.State
		var log []lm.Evidence
		at := int64(1762041600)
		for n := range 1 + r.IntN(40) {
			at += int64(r.IntN(200000))
			e := journalEvidence(t, r, w, n, at)
			s = mustAppend(t, s, e, w.graphs)
			log = append(log, e)
			if r.IntN(4) == 0 || n == 0 {
				rd := mustRead(t, s, w.graphs, at+int64(r.IntN(100*day)), lm.Grade{Value: int32(r.IntN(9)), Known: r.IntN(3) > 0}, engine...)
				for _, e := range log {
					for _, task := range variants(e.Task) {
						sameTask(t, "seed "+itoa(int(seed)), rd, task)
					}
				}
			}
		}
	}

	g, tasks := poolGraphs(t)
	for _, c := range poolStates(t, g, tasks) {
		for _, at := range []int64{poolTime, poolTime + 30*day} {
			r := mustRead(t, c.state, g, at, lm.Grade{Value: 3, Known: true}, poolModule)
			for _, task := range tasks {
				sameTask(t, "pool "+c.name, r, task)
			}
		}
	}
}

// TestTaskAllocatesOnce holds the fast path of Task (R29): on the pool, whose tags and items are in
// their order, a call allocates only the slice it gives. A skill transfer gives whose R at t₀ Read
// had not worked out would be inserted into a copy, and this test would see it.
func TestTaskAllocatesOnce(t *testing.T) {
	g, tasks := poolGraphs(t)
	for _, c := range poolStates(t, g, tasks) {
		r := mustRead(t, c.state, g, poolTime, lm.Grade{Value: 3, Known: true}, poolModule)
		var failed error
		allocs := testing.AllocsPerRun(3, func() {
			for _, task := range tasks {
				if _, err := r.Task(task); err != nil {
					failed = err
				}
			}
		})
		if failed != nil {
			t.Fatal(failed)
		}
		if allocs != float64(len(tasks)) {
			t.Errorf("%s: %v allocations for %d candidates, not one each", c.name, allocs, len(tasks))
		}
	}
}

// TestTaskConcurrently calls Task on one Readout from several goroutines at once, as a caller
// choosing for many requests may: Task only reads the Readout (docs/architecture.md, section 3.6,
// rule 7), so -race sees no race and every call gives the bits of a call alone.
func TestTaskConcurrently(t *testing.T) {
	g, tasks := poolGraphs(t)
	states := poolStates(t, g, tasks)
	r := mustRead(t, states[len(states)-1].state, g, poolTime, lm.Grade{Value: 3, Known: true}, poolModule)
	read := func() ([]lm.TaskReadout, error) {
		out := make([]lm.TaskReadout, 0, len(tasks))
		for _, task := range tasks {
			tr, err := r.Task(task)
			if err != nil {
				return nil, err
			}
			out = append(out, tr)
		}
		return out, nil
	}
	want, err := read()
	if err != nil {
		t.Fatal(err)
	}
	type result struct {
		out []lm.TaskReadout
		err error
	}
	const callers = 4
	results := make(chan result, callers)
	for range callers {
		go func() {
			out, err := read()
			results <- result{out, err}
		}()
	}
	for range callers {
		got := <-results
		if got.err != nil {
			t.Fatal(got.err)
		}
		for k, tr := range got.out {
			for i, p := range tr.Items {
				q := want[k].Items[i]
				if p.Place != q.Place || p.HasVerdict != q.HasVerdict || math.Float64bits(p.Probability) != math.Float64bits(q.Probability) {
					t.Fatalf("task %s: %+v at once, %+v alone", tasks[k].ID, p, q)
				}
			}
		}
	}
}
