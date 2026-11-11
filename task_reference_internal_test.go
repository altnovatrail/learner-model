package learnermodel

import (
	"cmp"
	"slices"
	"strings"
	"testing"

	"github.com/altnovatrail/learner-model/detmath"
)

// The reference of Readout.Task (T40 in RUN.md): Task as v0.1.5 had it, which copied the tags and
// the items of each candidate, split its ID, and worked out R of its skills at t₀ anew. Its text
// and that of its helpers are v0.1.5's word for word, but for their names and for one line: which
// modules Read has checked comes from the field that holds them now. The numbers of the model —
// guess, beta, chance, retentionAt and detmath — are those of the code, so the reference holds the
// order of the operations and the checks, and TestTaskIsReference holds Task to it bit for bit.

// TaskReference is Readout.Task of v0.1.5, for the tests outside the package.
func (r Readout) TaskReference(t Task) (TaskReadout, error) {
	if _, err := checkTaskContentReference(t); err != nil {
		return TaskReadout{}, err
	}
	module, ok := moduleOfReference(t.ID)
	if !ok {
		return TaskReadout{}, invalid("task " + shownBytes(t.ID) + " names no module")
	}
	if err := checkTaskPointReference("task "+t.ID, t.Band.Place, t.Difficulty); err != nil {
		return TaskReadout{}, err
	}
	items := slices.Clone(t.Items)
	slices.SortFunc(items, func(a, b TaskItem) int { return cmp.Compare(a.Place, b.Place) })
	for _, it := range items {
		if err := checkChoices(it.Slots); err != nil {
			return TaskReadout{}, err
		}
	}
	mg, ok := r.graphs[module]
	if !ok {
		return TaskReadout{}, invalid("module " + module + " has no entry in the graphs")
	}
	if !mg.Current {
		return TaskReadout{}, invalid("module " + module + " of task " + t.ID + " has no current version")
	}
	// Read has checked the graphs of its modules; a candidate of another module is checked here.
	if _, checked := r.checkedModule(module); !checked {
		if err := checkGraph(module, mg.Graph); err != nil {
			return TaskReadout{}, err
		}
	}
	var ms ModuleState
	if i, ok := r.state.moduleIndex(module); ok {
		ms = r.state.Modules[i]
	} else {
		start, ok, err := startingPoint(r.grade, mg.Graph.Bands)
		if err != nil {
			return TaskReadout{}, err
		}
		if !ok {
			return TaskReadout{}, invalid("task " + t.ID + " of module " + module + ", whose current version has no bands")
		}
		ms = newModule(module, start)
	}
	trial := ms.inTrial()
	var rs []skillRetention
	out := TaskReadout{Items: make([]ItemProbability, 0, len(items))}
	for _, it := range items {
		ip := ItemProbability{Place: it.Place}
		if len(it.Slots) > 0 && knownChecks(it.Slots) {
			ip.HasVerdict = true
			ip.Probability = ms.chanceOfReference(t, it, &mg.Graph, &rs, r.at, trial).p
		}
		out.Items = append(out.Items, ip)
	}
	return out, nil
}

func checkTaskContentReference(t Task) ([]int32, error) {
	if err := utf8Strings(t.ID, t.Band.ID); err != nil {
		return nil, err
	}
	if err := idListReference("primary skill", t.Primary); err != nil {
		return nil, err
	}
	if err := idListReference("supporting skill", t.Supporting); err != nil {
		return nil, err
	}
	places := make([]int32, 0, len(t.Items))
	for _, item := range t.Items {
		if err := idListReference("primary skill of an item", item.Primary); err != nil {
			return nil, err
		}
		if err := idListReference("supporting skill of an item", item.Supporting); err != nil {
			return nil, err
		}
		if err := checkSlots(item.Slots); err != nil {
			return nil, err
		}
		places = append(places, item.Place)
	}
	slices.Sort(places)
	return places, uniqueSorted("item place", places)
}

func idListReference(what string, ids []string) error {
	if err := utf8Strings(ids...); err != nil {
		return err
	}
	return uniqueSorted(what, sortedCopy(ids))
}

func moduleOfReference(id string) (string, bool) {
	parts := strings.SplitN(id, "/", 4)
	if len(parts) < 4 || slices.Contains(parts, "") {
		return "", false
	}
	return parts[0] + "/" + parts[1], true
}

func checkTaskPointReference(what string, place, difficulty int32) error {
	if place < 1 {
		return invalid(what + " in band place " + itoa(place))
	}
	if difficulty < 1 || difficulty > 5 {
		return invalid(what + " of difficulty " + itoa(difficulty))
	}
	return nil
}

func transferReference(g *Graph, id string) []Skill {
	if s, ok := graphSkill(g, id); ok {
		if s.Assessability == AssessabilityVerdict {
			return []Skill{s}
		}
		return nil
	}
	i, ok := slices.BinarySearchFunc(g.Retired, id, func(r RetiredSkill, id string) int { return strings.Compare(r.ID, id) })
	if !ok {
		return nil
	}
	var out []Skill
	for _, r := range g.Retired[i].ReplacedBy {
		if s, ok := graphSkill(g, r); ok && s.Assessability == AssessabilityVerdict {
			out = append(out, s)
		}
	}
	return out
}

// chanceOfReference is chanceOf of v0.1.5, which gathered the shares into its result: here they
// stay in a variable of its own, which nothing reads.
func (ms *ModuleState) chanceOfReference(t Task, ti TaskItem, g *Graph, rs *[]skillRetention, at int64, trial bool) itemChance {
	var out itemChance
	var shares []share
	var corrections, retentions float64
	for _, x := range sortedCopy(itemSkills(t, ti)) {
		to := transferReference(g, x)
		if len(to) == 0 {
			continue
		}
		out.j++
		var correction, retention float64
		for _, s := range to {
			correction += ms.correction(s.ID)
			if r := ms.retention(rs, s.ID, at); !trial {
				retention += r
			}
			i, ok := slices.BinarySearchFunc(shares, s.ID, func(sh share, id string) int { return strings.Compare(sh.skill.ID, id) })
			if !ok {
				shares = slices.Insert(shares, i, share{skill: s})
			}
			shares[i].k++
		}
		parts := float64(len(to))
		corrections += float64(correction / parts)
		if !trial {
			retentions += float64(retention / parts)
		}
	}
	mean := 0.0
	out.rbar = 1
	if out.j > 0 {
		mean = float64(corrections / float64(out.j))
		if !trial {
			out.rbar = float64(retentions / float64(out.j))
		}
	}
	out.c = guess(ti.Slots)
	out.sigma = detmath.Sigmoid(ms.Level + mean - beta(t.Band.Place, t.Difficulty))
	out.p1, out.p = chance(out.c, out.sigma, 1), chance(out.c, out.sigma, out.rbar)
	return out
}

// moduleOf gives what the split of v0.1.5 gave, on IDs of every shape it tells apart.
func TestModuleOfIsReference(t *testing.T) {
	ids := []string{
		"", "/", "a", "a/b", "a/b/c", "a/b/c/", "a/b/c/d", "a/b/c/d/e", "a/b/c//", "a/b//d", "a//c/d", "/b/c/d",
		"a/b/c/d/", "ab/cd/task/x", "example/m0/task/t1", "example/m0/task/t/1", "example/m0/task", "///", "////",
		"a/b/c/\xff", "\xff/b/c/d",
	}
	for _, id := range ids {
		got, ok := moduleOf(id)
		want, wantOK := moduleOfReference(id)
		if got != want || ok != wantOK {
			t.Errorf("moduleOf(%q) = %q, %v; v0.1.5 gave %q, %v", id, got, ok, want, wantOK)
		}
	}
}
